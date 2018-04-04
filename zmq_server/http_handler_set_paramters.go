package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"net/http"
	"time"
)

type SetParamters struct {
	Plc_id     *uint64 `json:"plc_id"`
	Serial     *uint32 `json:"serial"`
	LinkMode   *uint8  `json:"link_mode"`
	Wifi       *string `json:"wifi"`
	WifiPasswd *string `json:"wifi_passwd"`
}

func SetParamtersHandler(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var set_paramters SetParamters
	err := decoder.Decode(&set_paramters)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if set_paramters.Plc_id == nil ||
		set_paramters.Serial == nil ||
		set_paramters.LinkMode == nil {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	} else if *set_paramters.LinkMode == 3 && // 3 means wifi
		(set_paramters.Wifi == nil || set_paramters.WifiPasswd == nil) {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	_serial := uint32(GetHttpServer().SetSerialID(*set_paramters.Serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *set_paramters.Plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_SET_PARAMTERS,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*set_paramters.LinkMode),
			},
			&Report.Param{
				Type:    Report.Param_STRING,
				Strpara: *set_paramters.Wifi,
			},
			&Report.Param{
				Type:    Report.Param_STRING,
				Strpara: *set_paramters.WifiPasswd,
			},
		},
	}

	chan_key := GenerateKey(*set_paramters.Plc_id, _serial)

	chan_response := GetHttpServer().SendRequest(chan_key)
	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		result := (*Report.ControlCommand)(res).Paras[0].Npara
		fmt.Fprint(w, EncodingGeneralResponse(uint8(result)))
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
