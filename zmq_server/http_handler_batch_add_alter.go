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

type Alter struct {
	ID          uint8
	Modbus_Addr uint16
	Data_type   uint8
	Data_len    uint8
	CondOp      uint8
	Threshold   uint32
}

type BatchAddAlter struct {
	Plc_id      uint64
	Serial      uint32
	Serial_Port uint8
	Alters      []*Alter
}

type BatchAddAlterResponse struct {
	SerialPort uint8
	Result     uint8
}

func BatchAddAlterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var batch_add_alter BatchAddAlter
	err := decoder.Decode(&batch_add_alter)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Println(batch_add_alter)
	log.Println(batch_add_alter.Alters[0])
	log.Println(batch_add_alter.Alters[1])

	var jsonResult []byte
	defer func() {
		if x := recover(); x != nil {
			jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_SERVER_FAILED})
			fmt.Fprint(w, string(jsonResult))
			log.Println("3")
		}
	}()

	paras := []*Report.Param{
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(batch_add_alter.Serial_Port),
		},
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(len(batch_add_alter.Alters)),
		},
	}
	for _, alter := range batch_add_alter.Alters {
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(alter.ID),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT16,
			Npara: uint64(alter.Modbus_Addr),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(alter.Data_type),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(alter.Data_len),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(alter.CondOp),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(alter.Threshold),
		})
	}

	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          batch_add_alter.Plc_id,
		SerialNumber: batch_add_alter.Serial,
		Type:         Report.ControlCommand_CMT_REQ_BATCH_ADD_MONITOR,
		Paras:        paras,
	}

	chan_key := GenerateKey(batch_add_alter.Plc_id, batch_add_alter.Serial)

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
		result := uint8((*Report.ControlCommand)(res).Paras[1].Npara)

		batch_add_alter_response := &BatchAddAlterResponse{
			SerialPort: serial_port,
			Result:     result,
		}
		jsonResult, _ = json.Marshal(batch_add_alter_response)
		fmt.Fprint(w, string(jsonResult))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		close(chan_response)
		delete(GetHttpServer().HttpRespones, chan_key)
		jsonResult, _ = json.Marshal(map[string]string{HTTP_RESPONSE_RESULT: HTTP_RESPONSE_RESULT_TIMEOUT})
		fmt.Fprint(w, string(jsonResult))
	}
}
