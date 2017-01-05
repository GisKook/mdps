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

type Rs485SetConfig struct {
	Plc_id     *uint64
	Serial     *uint32
	SerialPort *uint8
	StartBit   *uint8
	EndBit     *uint8
	DataBit    *uint8
	Check      *uint8
	BaudRate   *uint32
}

func Rs485SetConfigHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	log.Println(r.PostForm)
	decoder := json.NewDecoder(r.Body)
	var rs485_set_cnfig Rs485SetConfig
	err := decoder.Decode(&rs485_set_cnfig)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if rs485_set_cnfig.Plc_id == nil ||
		rs485_set_cnfig.Serial == nil ||
		rs485_set_cnfig.SerialPort == nil ||
		rs485_set_cnfig.StartBit == nil ||
		rs485_set_cnfig.EndBit == nil ||
		rs485_set_cnfig.DataBit == nil ||
		rs485_set_cnfig.Check == nil ||
		rs485_set_cnfig.BaudRate == nil {
		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_PARAMTER_ERR))

		return
	}
	log.Println(rs485_set_cnfig)

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_SERVER_FAILED))
		}
	}()

	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *rs485_set_cnfig.Plc_id,
		SerialNumber: *rs485_set_cnfig.Serial,
		Type:         Report.ControlCommand_CMT_REQ_RS485_SET_CONFIG,
		Paras: []*Report.Param{
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs485_set_cnfig.SerialPort),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs485_set_cnfig.StartBit),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs485_set_cnfig.EndBit),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs485_set_cnfig.DataBit),
			},
			&Report.Param{
				Type:  Report.Param_UINT8,
				Npara: uint64(*rs485_set_cnfig.Check),
			},
			&Report.Param{
				Type:  Report.Param_UINT32,
				Npara: uint64(*rs485_set_cnfig.BaudRate),
			},
		},
	}

	chan_key := GenerateKey(*rs485_set_cnfig.Plc_id, *rs485_set_cnfig.Serial)

	chan_response, ok := GetHttpServer().HttpRespones[chan_key]

	if !ok {
		chan_response = make(chan *Report.ControlCommand)
		var once sync.Once
		once.Do(func() { GetHttpServer().HttpRespones[chan_key] = chan_response })
	}

	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		result := (*Report.ControlCommand)(res).Paras[0].Npara
		fmt.Fprint(w, EncodingGeneralResponse(uint8(result)))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		close(chan_response)
		delete(GetHttpServer().HttpRespones, chan_key)

		fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_TIMEOUT))
	}
}
