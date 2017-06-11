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

type DataDownload struct {
	Plc_id      *uint64
	Serial      *uint32
	Serial_Port *uint8
	Modbus_Addr *uint16
	Data_type   *uint8
	Data_len    *uint8
	Data        *string
}

func DataDownloadHandler(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var data_download DataDownload
	err := decoder.Decode(&data_download)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if data_download.Plc_id == nil ||
		data_download.Serial == nil ||
		data_download.Serial_Port == nil ||
		data_download.Modbus_Addr == nil ||
		data_download.Data_type == nil ||
		data_download.Data_len == nil ||
		data_download.Data == nil {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	data, _ := base64.StdEncoding.DecodeString(*data_download.Data)
	_serial := uint32(GetHttpServer().SetSerialID(*data_download.Serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *data_download.Plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_DATA_DOWNLOAD,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*data_download.Serial_Port),
			},
			&Report.Param{
				Type:  Report.Param_UINT16,
				Npara: uint64(*data_download.Modbus_Addr),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*data_download.Data_type),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*data_download.Data_len),
			},
			&Report.Param{
				Type:  Report.Param_BYTES,
				Bpara: data,
			},
		},
	}

	chan_key := GenerateKey(*data_download.Plc_id, _serial)

	chan_response := GetHttpServer().SendRequest(chan_key)
	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		result := (*Report.ControlCommand)(res).Paras[1].Npara
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
