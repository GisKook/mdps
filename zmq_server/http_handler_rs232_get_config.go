package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"net/http"
	"time"
)

type Rs232GetConfig struct {
	Plc_id      *uint64
	Serial      *uint32
	Serial_Port *uint8
}

type Rs232GetConfigResponse struct {
	Result      uint8  `json:"result"`
	Desc        string `json:"desc"`
	Serial_Port uint8  `json:"serial_port"`
	NodeType    uint8  `json:"node_type"`
	StationID   uint8  `json:"station_id"`
	StartBit    uint8  `json:"start_bit"`
	EndBit      uint8  `json:"end_bit"`
	DataBit     uint8  `json:"data_bit"`
	CheckBit    uint8  `json:"check_bit"`
	BaudRate    uint32 `json:"baud_rate"`
}

func EncodeRs232GetConfigResponse(response *Report.ControlCommand) string {
	serial_port := uint8(response.Paras[0].Npara)
	node_type := uint8(response.Paras[1].Npara)
	station_id := uint8(response.Paras[2].Npara)
	start_bit := uint8(response.Paras[3].Npara)
	end_bit := uint8(response.Paras[4].Npara)
	data_bit := uint8(response.Paras[5].Npara)
	check_bit := uint8(response.Paras[6].Npara)
	baud_rate := uint32(response.Paras[7].Npara)

	response_json, _ := json.Marshal(Rs232GetConfigResponse{
		Result:      HTTP_RESPONSE_RESULT_SUCCESS,
		Desc:        HTTP_RESULT[HTTP_RESPONSE_RESULT_SUCCESS],
		Serial_Port: serial_port,
		NodeType:    node_type,
		StationID:   station_id,
		StartBit:    start_bit,
		EndBit:      end_bit,
		DataBit:     data_bit,
		CheckBit:    check_bit,
		BaudRate:    baud_rate,
	})

	return string(response_json)
}

func Rs232GetConfigHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var rs232_get_config Rs232GetConfig
	err := decoder.Decode(&rs232_get_config)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if rs232_get_config.Plc_id == nil ||
		rs232_get_config.Serial == nil ||
		rs232_get_config.Serial_Port == nil {
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
		Tid:          *rs232_get_config.Plc_id,
		SerialNumber: *rs232_get_config.Serial,
		Type:         Report.ControlCommand_CMT_REQ_RS232_GET_CONFIG,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs232_get_config.Serial_Port),
			},
		},
	}

	chan_key := GenerateKey(*rs232_get_config.Plc_id, *rs232_get_config.Serial)

	chan_response := GetHttpServer().SendRequest(chan_key)

	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		fmt.Fprint(w, EncodeRs232GetConfigResponse(res))
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
