package zmq_server

import (
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"net/http"
	"sync"
	"time"
)

func RestartHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r)

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	r.ParseForm()
	if !CheckParamters(r, HTTP_PLC_ID, HTTP_PLC_SERIAL, HTTP_RESTART_DELAY) {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	plc_id, serial := GetPLCIDAndSerial(r)
	delay := GetUint8Value(r, HTTP_RESTART_DELAY)

	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          plc_id,
		SerialNumber: serial,
		Type:         Report.ControlCommand_CMT_REQ_RESTART,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(delay),
			},
		},
	}

	//data, _ := proto.Marshal(req)

	chan_key := GenerateKey(plc_id, serial)

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
		value := uint8((*Report.ControlCommand)(res).Paras[0].Npara)
		fmt.Fprint(w, EncodingGeneralResponse(value))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		if try_time < conf.GetConf().Http.TryTime {
			try_time++
			goto cmd
		} else {
			close(chan_response)
			delete(GetHttpServer().HttpRespones, chan_key)
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_TIMEOUT))
		}
	}
}
