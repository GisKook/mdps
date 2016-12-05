package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"net/http"
	"time"
)

func RestartHandler(w http.ResponseWriter, r *http.Request) {

	var jsonResult []byte
	defer func() {
		if x := recover(); x != nil {
			jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_SERVER_FAILED})
			log.Println("3")
		}
		log.Println("2")
		fmt.Fprint(w, string(jsonResult))
	}()

	r.ParseForm()
	if !CheckParamters(r, HTTP_PLC_ID, HTTP_PLC_SERIAL, HTTP_RESTART_DELAY) {
		jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_PARAMTER_ERR})
		log.Println("1")
		fmt.Fprint(w, string(jsonResult))

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
		GetHttpServer().HttpRespones[chan_key] = chan_response
	}

	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		value := (*Report.ControlCommand)(res).Paras[0].Npara
		jsonResult, _ = json.Marshal(map[string]uint64{HTTP_RESPONSE_RESULT: value})
		fmt.Fprint(w, string(jsonResult))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout)):
		close(chan_response)
		delete(GetHttpServer().HttpRespones, chan_key)
		jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_TIMEOUT})
		fmt.Fprint(w, string(jsonResult))
	}
}
