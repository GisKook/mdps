package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"net/http"
	"time"
)

type Alter struct {
	ID          *uint8
	Modbus_Addr *uint16
	Data_type   *uint8
	Data_len    *uint8
	CondOp      *uint8
	Threshold   *uint32
}

type BatchAddAlter struct {
	Plc_id      *uint64
	Serial      *uint32
	Serial_Port *uint8
	Alters      *[]*Alter
}

//type BatchAddAlterResponse struct {
//	Result uint8  `json:result`
//	Desc   string `json:desc`
//}

//func EncodingResponse(result uint8) string {
//	batch_add_alter_response := &BatchAddAlterResponse{
//		Result: result,
//		Desc:   HTTP_RESULT[result],
//	}
//
//	response, _ := json.Marshal(batch_add_alter_response)
//
//	return string(response)
//}

func CheckParamtersErr(batch_add_alter *BatchAddAlter) bool {
	if batch_add_alter.Plc_id == nil ||
		batch_add_alter.Serial == nil ||
		batch_add_alter.Serial_Port == nil ||
		batch_add_alter.Alters == nil {
		return true
	}

	for _, alter := range *batch_add_alter.Alters {
		if alter == nil {
			return true
		}
		if alter.ID == nil ||
			alter.Modbus_Addr == nil ||
			alter.Data_type == nil ||
			alter.Data_len == nil ||
			alter.CondOp == nil ||
			alter.Threshold == nil {
			return true
		}
	}

	return false
}

func BatchAddAlterHandler(w http.ResponseWriter, r *http.Request) {
	PrintRequest(r)
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var batch_add_alter BatchAddAlter
	err := decoder.Decode(&batch_add_alter)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	if CheckParamtersErr(&batch_add_alter) {
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
			Npara: uint64(*batch_add_alter.Serial_Port),
		},
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(len(*batch_add_alter.Alters)),
		},
	}
	for _, alter := range *batch_add_alter.Alters {
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*alter.ID),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT16,
			Npara: uint64(*alter.Modbus_Addr),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*alter.Data_type),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*alter.Data_len),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*alter.CondOp),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*alter.Threshold),
		})
	}

	_serial := uint32(GetHttpServer().SetSerialID(*batch_add_alter.Serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *batch_add_alter.Plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_BATCH_ADD_ALTER,
		Paras:        paras,
	}

	chan_key := GenerateKey(*batch_add_alter.Plc_id, _serial)

	chan_response := GetHttpServer().SendRequest(chan_key)
	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		result := uint8((*Report.ControlCommand)(res).Paras[1].Npara)
		fmt.Fprint(w, EncodingGeneralResponse(result))
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
