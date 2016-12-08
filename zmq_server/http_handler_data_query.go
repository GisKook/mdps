package zmq_server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"net/http"
	"sync"
	"time"
)

type DataQuery struct {
	Plc_id      uint64
	Serial      uint32
	Serial_Port uint8
	Modbus_Addr uint16
	Data_type   uint8
	Data_len    uint8
}

type DataQueryResponse struct {
	SerialPort uint8
	Datatype   uint8
	DataLen    uint8
	Data       string
}

func DataQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var data_query DataQuery
	err := decoder.Decode(&data_query)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Println(data_query)

	var jsonResult []byte
	defer func() {
		if x := recover(); x != nil {
			jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_SERVER_FAILED})
			fmt.Fprint(w, string(jsonResult))
			log.Println("3")
		}
	}()

	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          data_query.Plc_id,
		SerialNumber: data_query.Serial,
		Type:         Report.ControlCommand_CMT_REQ_DATA_QUERY,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(data_query.Serial_Port),
			},
			&Report.Param{
				Type:  Report.Param_UINT16,
				Npara: uint64(data_query.Modbus_Addr),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(data_query.Data_type),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(data_query.Data_len),
			},
		},
	}

	chan_key := GenerateKey(data_query.Plc_id, data_query.Serial)

	chan_response, ok := GetHttpServer().HttpRespones[chan_key]

	if !ok {
		chan_response = make(chan *Report.ControlCommand)
		var once sync.Once
		once.Do(func() { GetHttpServer().HttpRespones[chan_key] = chan_response })
	}

	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		serial_port := uint8((*Report.ControlCommand)(res).Paras[0].Npara)
		data_type := uint8((*Report.ControlCommand)(res).Paras[1].Npara)
		data_len := uint8((*Report.ControlCommand)(res).Paras[2].Npara)
		data := (*Report.ControlCommand)(res).Paras[3].Bpara
		data_base64 := base64.StdEncoding.EncodeToString(data)

		data_query_response := &DataQueryResponse{
			SerialPort: serial_port,
			Datatype:   data_type,
			DataLen:    data_len,
			Data:       data_base64,
		}
		jsonResult, _ = json.Marshal(data_query_response)
		fmt.Fprint(w, string(jsonResult))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		close(chan_response)
		delete(GetHttpServer().HttpRespones, chan_key)
		jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_TIMEOUT})
		fmt.Fprint(w, string(jsonResult))
	}
}
