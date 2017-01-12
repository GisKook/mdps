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

type Monitor struct {
	ID          *uint8
	Modbus_Addr *uint16
	Data_type   *uint8
	Data_len    *uint8
}

type BatchAddMonitor struct {
	Plc_id      *uint64
	Serial      *uint32
	Serial_Port *uint8
	Monitors    *[]*Monitor
}

//type BatchAddMonitorResponse struct {
//	SerialPort uint8
//	Result     uint8
//}
func CheckParamtersBatchAddMonitorErr(batch_add_monitor *BatchAddMonitor) bool {
	if batch_add_monitor.Plc_id == nil ||
		batch_add_monitor.Serial == nil ||
		batch_add_monitor.Serial_Port == nil ||
		batch_add_monitor.Monitors == nil {
		return true
	}

	for _, monitor := range *batch_add_monitor.Monitors {
		if monitor == nil {
			return true
		}
		if monitor.ID == nil ||
			monitor.Modbus_Addr == nil ||
			monitor.Data_type == nil ||
			monitor.Data_len == nil {

			return true
		}
	}

	return false
}

func BatchAddMonitorHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var batch_add_monitor BatchAddMonitor
	err := decoder.Decode(&batch_add_monitor)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	if CheckParamtersBatchAddMonitorErr(&batch_add_monitor) {
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
			Npara: uint64(*batch_add_monitor.Serial_Port),
		},
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(len(*batch_add_monitor.Monitors)),
		},
	}
	for _, monitor := range *batch_add_monitor.Monitors {
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*monitor.ID),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT16,
			Npara: uint64(*monitor.Modbus_Addr),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*monitor.Data_type),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*monitor.Data_len),
		})
	}

	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *batch_add_monitor.Plc_id,
		SerialNumber: *batch_add_monitor.Serial,
		Type:         Report.ControlCommand_CMT_REQ_BATCH_ADD_MONITOR,
		Paras:        paras,
	}

	chan_key := GenerateKey(*batch_add_monitor.Plc_id, *batch_add_monitor.Serial)

	chan_response, ok := GetHttpServer().HttpRespones[chan_key]

	if !ok {
		chan_response = make(chan *Report.ControlCommand)
		var once sync.Once
		once.Do(func() { GetHttpServer().HttpRespones[chan_key] = chan_response })
	}

	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		//	serial_port := uint8((*Report.ControlCommand)(res).Paras[0].Npara)
		result := uint8((*Report.ControlCommand)(res).Paras[1].Npara)
		log.Printf("normal result %d\n", result)

		fmt.Fprint(w, EncodingGeneralResponse(result))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		if try_time < conf.GetConf().Http.TryTime {
			try_time++
			goto cmd
		} else {
			close(chan_response)
			delete(GetHttpServer().HttpRespones, chan_key)
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_TIMEOUT))
		}
	}
}
