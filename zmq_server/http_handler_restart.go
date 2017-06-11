package zmq_server

import (
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"net/http"
	"time"
)

func RestartHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()
	PrintRequest(r)

	r.ParseForm()
	if !CheckParamters(r, HTTP_PLC_ID, HTTP_PLC_SERIAL, HTTP_RESTART_DELAY) {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	plc_id, serial := GetPLCIDAndSerial(r)
	delay := GetUint8Value(r, HTTP_RESTART_DELAY)

	_serial := uint32(GetHttpServer().SetSerialID(serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_RESTART,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(delay),
			},
		},
	}

	//data, _ := proto.Marshal(req)

	chan_key := GenerateKey(plc_id, _serial)
	chan_response := GetHttpServer().SendRequest(chan_key)

	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		value := uint8((*Report.ControlCommand)(res).Paras[0].Npara)
		fmt.Fprint(w, EncodingGeneralResponse(value))
		GetHttpServer().DelRequest(chan_key)

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		if try_time < conf.GetConf().Http.TryTime {
			try_time++
			goto cmd
		} else {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_TIMEOUT))
			GetHttpServer().DelRequest(chan_key)
		}
	}
}
