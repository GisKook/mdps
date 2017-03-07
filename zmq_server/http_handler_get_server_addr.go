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

type GetServerAddr struct {
	Plc_id     *uint64
	Serial     *uint32
	Addr_Count *uint8
}

type ServerAddr struct {
	ConnectionType uint8  `json:"connection_type"`
	ServerClass    uint8  `json:"server_class"`
	Addr           string `json:"addr"`
	Port           uint16 `json:"port"`
}

type GetServerAddrResponse struct {
	Result uint8         `json:"result"`
	Desc   string        `json:"desc"`
	Addrs  []*ServerAddr `json:"servers"`
}

func EncodeGetServerAddrResponse(response *Report.ControlCommand) string {
	server_count := response.Paras[0].Npara
	var get_server_addr_reponse GetServerAddrResponse
	get_server_addr_reponse.Result = HTTP_RESPONSE_RESULT_SUCCESS
	get_server_addr_reponse.Desc = HTTP_RESULT[HTTP_RESPONSE_RESULT_SUCCESS]
	for i := 0; i < int(server_count); i++ {
		get_server_addr_reponse.Addrs = append(get_server_addr_reponse.Addrs, &ServerAddr{
			ConnectionType: uint8(response.Paras[i*4+1].Npara),
			ServerClass:    uint8(response.Paras[i*4+2].Npara),
			Addr:           response.Paras[i*4+3].Strpara,
			Port:           uint16(response.Paras[i*4+4].Npara),
		})
	}
	response_json, _ := json.Marshal(get_server_addr_reponse)

	return string(response_json)
}

func GetServerAddrHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetServerAddrHandler")
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var get_server_addr GetServerAddr
	err := decoder.Decode(&get_server_addr)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if get_server_addr.Plc_id == nil ||
		get_server_addr.Serial == nil ||
		get_server_addr.Addr_Count == nil {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *get_server_addr.Plc_id,
		SerialNumber: *get_server_addr.Serial,
		Type:         Report.ControlCommand_CMT_REQ_GET_SERVER_ADDR,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*get_server_addr.Addr_Count),
			},
		},
	}

	chan_key := GenerateKey(*get_server_addr.Plc_id, *get_server_addr.Serial)

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
		fmt.Fprint(w, EncodeGetServerAddrResponse(res))

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
