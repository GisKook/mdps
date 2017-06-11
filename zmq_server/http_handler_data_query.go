package zmq_server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"net/http"
	"time"
)

type DataQuery struct {
	Plc_id      *uint64
	Serial      *uint32
	Serial_Port *uint8
	Modbus_Addr *uint16
	Data_type   *uint8
	Data_len    *uint8
}

type DataQueryResponse struct {
	Result     uint8  `json:"result"`
	Desc       string `json:"desc"`
	SerialPort uint8  `json:"serial_port"`
	Datatype   uint8  `json:"data_type"`
	DataLen    uint8  `json:"data_len"`
	Data       string `json:"data"`
}

func CheckParamtersDataQueryErr(data_query *DataQuery) bool {
	if data_query.Plc_id == nil ||
		data_query.Serial == nil ||
		data_query.Serial_Port == nil ||
		data_query.Modbus_Addr == nil ||
		data_query.Data_type == nil ||
		data_query.Data_len == nil {
		return true
	}

	return false
}

func EncodingDataQueryResponse(data_query_response *DataQueryResponse) string {
	response, _ := json.Marshal(data_query_response)

	return string(response)
}

func DataQueryHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingDataQueryResponse(
				&DataQueryResponse{
					Result: HTTP_RESPONSE_RESULT_SERVER_FAILED,
					Desc:   HTTP_RESULT[HTTP_RESPONSE_RESULT_SERVER_FAILED],
				}))
		}
	}()
	PrintRequest(r)

	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var data_query DataQuery
	err := decoder.Decode(&data_query)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	if CheckParamtersDataQueryErr(&data_query) {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))
		return
	}

	_serial := uint32(GetHttpServer().SetSerialID(*data_query.Serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *data_query.Plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_DATA_QUERY,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*data_query.Serial_Port),
			},
			&Report.Param{
				Type:  Report.Param_UINT16,
				Npara: uint64(*data_query.Modbus_Addr),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*data_query.Data_type),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*data_query.Data_len),
			},
		},
	}

	chan_key := GenerateKey(*data_query.Plc_id, _serial)
	chan_response := GetHttpServer().SendRequest(chan_key)

	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		serial_port := uint8((*Report.ControlCommand)(res).Paras[0].Npara)
		data_type := uint8((*Report.ControlCommand)(res).Paras[1].Npara)
		data_len := uint8((*Report.ControlCommand)(res).Paras[2].Npara)
		data := (*Report.ControlCommand)(res).Paras[3].Bpara
		data_base64 := base64.StdEncoding.EncodeToString(data)

		data_query_response := &DataQueryResponse{
			Result:     HTTP_RESPONSE_RESULT_SUCCESS,
			Desc:       HTTP_RESULT[HTTP_RESPONSE_RESULT_SUCCESS],
			SerialPort: serial_port,
			Datatype:   data_type,
			DataLen:    data_len,
			Data:       data_base64,
		}
		fmt.Fprint(w, EncodingDataQueryResponse(data_query_response))
		GetHttpServer().DelRequest(chan_key)

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		if try_time < conf.GetConf().Http.TryTime {
			try_time++
			goto cmd
		} else {
			fmt.Fprint(w, EncodingDataQueryResponse(&DataQueryResponse{
				Result: HTTP_RESPONSE_RESULT_TIMEOUT,
				Desc:   HTTP_RESULT[HTTP_RESPONSE_RESULT_TIMEOUT],
			}))
			GetHttpServer().DelRequest(chan_key)
		}
	}
}
