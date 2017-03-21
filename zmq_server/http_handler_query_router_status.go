package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"net/http"
	"sync"
	"time"
)

type QueryRouterStatus struct {
	Plc_id *uint64
	Serial *uint32
}

func CheckParamtersQueryRouterStatusAddrErr(release_query_router_status *QueryRouterStatus) bool {
	if release_query_router_status.Plc_id == nil ||
		release_query_router_status.Serial == nil {
		return true
	}

	return false
}

type QueryRouterStatusResponse struct {
	Result uint8  `json:"result"`
	Desc   string `json:"desc"`
	Status uint8  `json:"status"`
}

func EncodeQueryRouterStatusResponse(response *Report.ControlCommand) string {
	status := uint8(response.Paras[0].Npara)

	response_json, _ := json.Marshal(Rs232GetConfigResponse{
		Result: HTTP_RESPONSE_RESULT_SUCCESS,
		Desc:   HTTP_RESULT[HTTP_RESPONSE_RESULT_SUCCESS],
		Status: status,
	})
}

func QueryRouterStatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("QueryRouterStatus")
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var query_router_status QueryRouterStatus
	err := decoder.Decode(&query_router_status)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if CheckParamtersQueryRouterStatusAddrErr(&query_router_status) {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	paras := []*Report.Param{
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*query_router_status.Serial_Port),
		},
	}

	log.Println(paras)
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *query_router_status.Plc_id,
		SerialNumber: *query_router_status.Serial,
		Type:         Report.ControlCommand_CMT_REQ_QUERY_ROUTER_STATUS,
		Paras:        paras,
	}

	chan_key := GenerateKey(*query_router_status.Plc_id, *query_router_status.Serial)

	chan_response, ok := GetHttpServer().HttpRespones[chan_key]

	if !ok {
		chan_response = make(chan *Report.ControlCommand)
		var once sync.Once
		once.Do(func() { GetHttpServer().HttpRespones[chan_key] = chan_response })
	}

	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		fmt.Fprint(w, EncodeQueryRouterStatusResponse(res))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		if try_time < conf.GetConf().Http.TryTime {
			try_time++
			goto cmd
		} else {
			close(chan_response)
			var once sync.Once
			once.Do(func() { delete(GetHttpServer().HttpRespones, chan_key) })

			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_TIMEOUT))
		}
	}
}
