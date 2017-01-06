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

type Rs485GetConfig struct {
	Plc_id      *uint64
	Serial      *uint32
	Serial_Port *uint8
}

type Rs485GetConfigResponse struct {
	Result      uint8  `json:"result"`
	Desc        string `json:"desc"`
	Serial_Port uint8  `json:"serial_port"`
	StartBit    uint8  `json:"start_bit"`
	EndBit      uint8  `json:"end_bit"`
	DataBit     uint8  `json:"data_bit"`
	CheckBit    uint8  `json:"check_bit"`
	BaudRate    uint32 `json:"baud_rate"`
}

func EncodeRs485GetConfigResponse(response *Report.ControlCommand) string {
	serial_port := uint8(response.Paras[0].Npara)
	start_bit := uint8(response.Paras[1].Npara)
	end_bit := uint8(response.Paras[2].Npara)
	data_bit := uint8(response.Paras[3].Npara)
	check_bit := uint8(response.Paras[4].Npara)
	baud_rate := uint32(response.Paras[5].Npara)

	response_json, _ := json.Marshal(Rs485GetConfigResponse{
		Result:      HTTP_RESPONSE_RESULT_SUCCESS,
		Desc:        HTTP_RESULT[HTTP_RESPONSE_RESULT_SUCCESS],
		Serial_Port: serial_port,
		StartBit:    start_bit,
		EndBit:      end_bit,
		DataBit:     data_bit,
		CheckBit:    check_bit,
		BaudRate:    baud_rate,
	})

	return string(response_json)
}

func Rs485GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	log.Println(r.PostForm)
	decoder := json.NewDecoder(r.Body)
	var rs485_get_config Rs485GetConfig
	err := decoder.Decode(&rs485_get_config)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if rs485_get_config.Plc_id == nil ||
		rs485_get_config.Serial == nil ||
		rs485_get_config.Serial_Port == nil {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}
	log.Println(rs485_get_config)

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *rs485_get_config.Plc_id,
		SerialNumber: *rs485_get_config.Serial,
		Type:         Report.ControlCommand_CMT_REQ_RS485_GET_CONFIG,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs485_get_config.Serial_Port),
			},
		},
	}

	chan_key := GenerateKey(*rs485_get_config.Plc_id, *rs485_get_config.Serial)

	chan_response, ok := GetHttpServer().HttpRespones[chan_key]

	if !ok {
		chan_response = make(chan *Report.ControlCommand)
		var once sync.Once
		once.Do(func() { GetHttpServer().HttpRespones[chan_key] = chan_response })
	}

	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		fmt.Fprint(w, EncodeRs485GetConfigResponse(res))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		close(chan_response)
		delete(GetHttpServer().HttpRespones, chan_key)

		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_TIMEOUT))
	}
}
