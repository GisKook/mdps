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

type ReleaseTransparentTransmission struct {
	Plc_id      *uint64
	Serial      *uint32
	Serial_Port *uint8
}

func CheckParamtersReleaseTransparentTransmissionAddrErr(release_release_transparent_transmission *ReleaseTransparentTransmission) bool {
	if release_release_transparent_transmission.Plc_id == nil ||
		release_release_transparent_transmission.Serial == nil ||
		release_release_transparent_transmission.Serial_Port == nil {
		return true
	}

	return false
}

func ReleaseTransparentTransmissionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ReleaseTransparentTransmission")
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var release_transparent_transmission ReleaseTransparentTransmission
	err := decoder.Decode(&release_transparent_transmission)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if CheckParamtersReleaseTransparentTransmissionAddrErr(&release_transparent_transmission) {
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
			Npara: uint64(*release_transparent_transmission.Serial_Port),
		},
	}

	log.Println(paras)
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *release_transparent_transmission.Plc_id,
		SerialNumber: *release_transparent_transmission.Serial,
		Type:         Report.ControlCommand_CMT_REQ_RELEASE_TRANSPARENT_TRANSMISSION,
		Paras:        paras,
	}

	chan_key := GenerateKey(*release_transparent_transmission.Plc_id, *release_transparent_transmission.Serial)

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
		result := (*Report.ControlCommand)(res).Paras[1].Npara
		fmt.Fprint(w, EncodingGeneralResponse(uint8(result)))

		break
	case <-time.After(time.Duration(conf.GetConf().Http.Timeout) * time.Second):
		if try_time < conf.GetConf().Http.TryTime {
			try_time++
			goto cmd
		} else {
			close(chan_response)
			var once sync.Once
			once.Do(func() { delete(GetHttpServer().HttpRespones, chan_key) })

			fmt.Fprint(w, EncodingGeneralResponse(HTTP_RESPONSE_RESULT_TIMEOUT))
		}
	}
}