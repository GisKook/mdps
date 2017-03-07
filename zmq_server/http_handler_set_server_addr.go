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

type Server_Addr struct {
	Connection_Type *uint8
	Server_Class    *uint8
	Server_Addr     *string
	Port            *uint16
}

type SetServer_Addr struct {
	Plc_id  *uint64
	Serial  *uint32
	Servers *[]*Server_Addr
}

func CheckParamtersSetServerAddrErr(server_addrs *SetServer_Addr) bool {
	if server_addrs.Plc_id == nil ||
		server_addrs.Serial == nil ||
		server_addrs.Servers == nil {
		return true
	}

	for _, server_addr := range *server_addrs.Servers {
		if server_addr == nil {
			return true
		}
		if server_addr.Connection_Type == nil ||
			server_addr.Server_Class == nil ||
			server_addr.Server_Addr == nil ||
			server_addr.Port == nil {
			return true
		}
	}

	return false

}

func SetServerAddrHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SetServerAddrHandler")
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	var set_server_addr SetServer_Addr
	err := decoder.Decode(&set_server_addr)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if CheckParamtersSetServerAddrErr(&set_server_addr) {
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
			Npara: uint64(len(*set_server_addr.Servers)),
		},
	}

	for _, server := range *set_server_addr.Servers {
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*server.Connection_Type),
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*server.Server_Class),
		})
		paras = append(paras, &Report.Param{
			Type:    Report.Param_STRING,
			Strpara: *server.Server_Addr,
		})
		paras = append(paras, &Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(*server.Port),
		})
	}
	log.Println(paras)
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *set_server_addr.Plc_id,
		SerialNumber: *set_server_addr.Serial,
		Type:         Report.ControlCommand_CMT_REQ_SET_SERVER_ADDR,
		Paras:        paras,
	}

	chan_key := GenerateKey(*set_server_addr.Plc_id, *set_server_addr.Serial)

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
		result := (*Report.ControlCommand)(res).Paras[0].Npara
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
