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

type DataDownload struct {
	Plc_id      uint64
	Serial      uint32
	Serial_Port uint8
	Modbus_Addr uint16
	Data_type   uint8
	Data_len    uint8
	Data        string
}

func DataDownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var data_download DataDownload
	err := decoder.Decode(&data_download)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Println(data_download)

	var jsonResult []byte
	defer func() {
		if x := recover(); x != nil {
			jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_SERVER_FAILED})
			fmt.Fprint(w, string(jsonResult))
			log.Println("3")
		}
	}()

	data, _ := base64.StdEncoding.DecodeString(data_download.Data)
	log.Println(data)
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          data_download.Plc_id,
		SerialNumber: data_download.Serial,
		Type:         Report.ControlCommand_CMT_REQ_DATA_DOWNLOAD,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(data_download.Serial_Port),
			},
			&Report.Param{
				Type:  Report.Param_UINT16,
				Npara: uint64(data_download.Modbus_Addr),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(data_download.Data_type),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(data_download.Data_len),
			},
			&Report.Param{
				Type:  Report.Param_BYTES,
				Bpara: data,
			},
		},
	}

	chan_key := GenerateKey(data_download.Plc_id, data_download.Serial)

	chan_response, ok := GetHttpServer().HttpRespones[chan_key]

	if !ok {
		chan_response = make(chan *Report.ControlCommand)
		var once sync.Once
		once.Do(func() { GetHttpServer().HttpRespones[chan_key] = chan_response })
	}

	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		serial_port := (*Report.ControlCommand)(res).Paras[0].Npara
		result := (*Report.ControlCommand)(res).Paras[1].Npara
		jsonResult, _ = json.Marshal(map[string]uint64{HTTP_RESPONSE_SERIAL_PORT: serial_port, HTTP_RESPONSE_RESULT: result})
		fmt.Fprint(w, string(jsonResult))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		close(chan_response)
		delete(GetHttpServer().HttpRespones, chan_key)
		jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_TIMEOUT})
		fmt.Fprint(w, string(jsonResult))
	}
}
