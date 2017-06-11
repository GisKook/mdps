package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"net/http"
	"time"
)

type TransparentTransmission struct {
	Plc_id                            *uint64
	Serial                            *uint32
	Serial_Port                       *uint8
	Connection_Type                   *uint8
	Server_Class                      *uint8
	Server_Addr                       *string
	Port                              *uint16
	Transparent_Transmission_ClientID *uint64
	Transparent_Transmission_Key      *uint32
}

func CheckParamtersTransparentTransmissionAddrErr(transparent_transmission *TransparentTransmission) bool {
	if transparent_transmission.Plc_id == nil ||
		transparent_transmission.Serial == nil ||
		transparent_transmission.Serial_Port == nil ||
		transparent_transmission.Connection_Type == nil ||
		transparent_transmission.Server_Class == nil ||
		transparent_transmission.Server_Addr == nil ||
		transparent_transmission.Port == nil ||
		transparent_transmission.Transparent_Transmission_ClientID == nil ||
		transparent_transmission.Transparent_Transmission_Key == nil {
		return true
	}

	return false
}

func TransparentTransmissionHandler(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var transparent_transmission TransparentTransmission
	err := decoder.Decode(&transparent_transmission)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if CheckParamtersTransparentTransmissionAddrErr(&transparent_transmission) {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	paras := []*Report.Param{
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*transparent_transmission.Serial_Port),
		},
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*transparent_transmission.Connection_Type),
		},
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*transparent_transmission.Server_Class),
		},
		&Report.Param{
			Type:    Report.Param_STRING,
			Strpara: *transparent_transmission.Server_Addr,
		},
		&Report.Param{
			Type:  Report.Param_UINT16,
			Npara: uint64(*transparent_transmission.Port),
		},
		&Report.Param{
			Type:  Report.Param_UINT64,
			Npara: uint64(*transparent_transmission.Transparent_Transmission_ClientID),
		},
		&Report.Param{
			Type:  Report.Param_UINT32,
			Npara: uint64(*transparent_transmission.Transparent_Transmission_Key),
		},
	}

	_serial := uint32(GetHttpServer().SetSerialID(*transparent_transmission.Serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *transparent_transmission.Plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_TRANSPARENT_TRANSMISSION,
		Paras:        paras,
	}

	chan_key := GenerateKey(*transparent_transmission.Plc_id, _serial)

	chan_response := GetHttpServer().SendRequest(chan_key)
	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		result := (*Report.ControlCommand)(res).Paras[1].Npara
		fmt.Fprint(w, EncodingGeneralResponse(uint8(result)))

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
