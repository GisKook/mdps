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

type SetServerAddr struct {
	Plc_id          *uint64
	Serial          *uint32
	Connection_Type *uint8
	Server_Class    *uint8
	ServerAddr      *string
	Port            *uint16
}

func SetServerAddrHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	log.Println(r.PostForm)
	decoder := json.NewDecoder(r.Body)
	var set_server_addr SetServerAddr
	err := decoder.Decode(&set_server_addr)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if set_server_addr.Plc_id == nil ||
		set_server_addr.Serial == nil ||
		set_server_addr.Connection_Type == nil ||
		set_server_addr.Server_Class == nil ||
		set_server_addr.ServerAddr == nil ||
		set_server_addr.Port == nil {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}
	log.Println(set_server_addr)

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *set_server_addr.Plc_id,
		SerialNumber: *set_server_addr.Serial,
		Type:         Report.ControlCommand_CMT_REQ_SET_SERVER_ADDR,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*set_server_addr.Connection_Type),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*set_server_addr.Server_Class),
			},
			&Report.Param{
				Type:    Report.Param_STRING,
				Strpara: *set_server_addr.ServerAddr,
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*set_server_addr.Port),
			},
		},
	}

	chan_key := GenerateKey(*set_server_addr.Plc_id, *set_server_addr.Serial)

	chan_response, ok := GetHttpServer().HttpRespones[chan_key]

	if !ok {
		chan_response = make(chan *Report.ControlCommand)
		var once sync.Once
		once.Do(func() { GetHttpServer().HttpRespones[chan_key] = chan_response })
	}

	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		result := (*Report.ControlCommand)(res).Paras[0].Npara
		fmt.Fprint(w, EncodingGeneralResponse(uint8(result)))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		close(chan_response)
		delete(GetHttpServer().HttpRespones, chan_key)

		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_TIMEOUT))
	}
}
