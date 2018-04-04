package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"net/http"
	"time"
)

type GetParamters struct {
	Plc_id *uint64
	Serial *uint32
}

type ParamtersResponse struct {
	Result     uint8  `json:"result"`
	Desc       string `json:"desc"`
	LinkMode   uint8  `json:"link_mode"`
	Wifi       string `json:"wifi"`
	WifiPasswd string `json:"wifi_passwd"`
}

func EncodeGetParamtersResponse(response *Report.ControlCommand) string {
	link_mode := response.Paras[0].Npara
	wifi := response.Paras[1].Strpara
	wifi_passwd := response.Paras[2].Strpara
	var paramters_response ParamtersResponse
	paramters_response.Result = HTTP_RESPONSE_RESULT_SUCCESS
	paramters_response.Desc = HTTP_RESULT[HTTP_RESPONSE_RESULT_SUCCESS]
	paramters_response.LinkMode = uint8(link_mode)
	paramters_response.Wifi = wifi
	paramters_response.WifiPasswd = wifi_passwd

	response_json, _ := json.Marshal(paramters_response)

	return string(response_json)
}

func GetParamtersHandler(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var get_paramters GetParamters
	err := decoder.Decode(&get_paramters)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if get_paramters.Plc_id == nil ||
		get_paramters.Serial == nil {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	_serial := uint32(GetHttpServer().SetSerialID(*get_paramters.Serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *get_paramters.Plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_GET_PARAMTERS,
	}

	chan_key := GenerateKey(*get_paramters.Plc_id, _serial)

	chan_response := GetHttpServer().SendRequest(chan_key)
	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		fmt.Fprint(w, EncodeGetParamtersResponse(res))
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
