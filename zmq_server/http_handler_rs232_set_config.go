package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type Rs232SetConfig struct {
	Plc_id      *uint64
	Serial      *uint32
	Serial_Port *uint8
	Node_Type   *uint8
	Station_id  *uint8
	Start_Bit   *uint8
	End_Bit     *uint8
	Data_Bit    *uint8
	Check_Bit   *uint8
	Baud_Rate   *uint32
}

func Rs232SetConfigHandler(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var rs232_set_cnfig Rs232SetConfig
	err := decoder.Decode(&rs232_set_cnfig)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if rs232_set_cnfig.Plc_id == nil ||
		rs232_set_cnfig.Serial == nil ||
		rs232_set_cnfig.Serial_Port == nil ||
		rs232_set_cnfig.Node_Type == nil ||
		rs232_set_cnfig.Station_id == nil ||
		rs232_set_cnfig.Start_Bit == nil ||
		rs232_set_cnfig.End_Bit == nil ||
		rs232_set_cnfig.Data_Bit == nil ||
		rs232_set_cnfig.Check_Bit == nil ||
		rs232_set_cnfig.Baud_Rate == nil {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	defer func() {
		if x := recover(); x != nil {
			log.Printf("%s %s\n", x, debug.Stack())
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	_serial := uint32(GetHttpServer().SetSerialID(*rs232_set_cnfig.Serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *rs232_set_cnfig.Plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_RS232_SET_CONFIG,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs232_set_cnfig.Serial_Port),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs232_set_cnfig.Node_Type),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs232_set_cnfig.Station_id),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs232_set_cnfig.Start_Bit),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs232_set_cnfig.End_Bit),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs232_set_cnfig.Data_Bit),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs232_set_cnfig.Check_Bit),
			},
			&Report.Param{
				Type:  Report.Param_UINT32,
				Npara: uint64(*rs232_set_cnfig.Baud_Rate),
			},
		},
	}

	chan_key := GenerateKey(*rs232_set_cnfig.Plc_id, _serial)

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
