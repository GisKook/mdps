package zmq_server

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"net/http"
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
	PrintRequest(r)
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
	_serial := uint32(GetHttpServer().SetSerialID(*set_server_addr.Serial))
	req := &Report.ControlCommand{
		Uuid:         "das",
		Tid:          *set_server_addr.Plc_id,
		SerialNumber: _serial,
		Type:         Report.ControlCommand_CMT_REQ_SET_SERVER_ADDR,
		Paras:        paras,
	}

	chan_key := GenerateKey(*set_server_addr.Plc_id, _serial)

	chan_response := GetHttpServer().SendRequest(chan_key)
	try_time := uint8(0)
cmd:
	GetZmqServer().SendControlDown(req)

	select {
	case res := <-chan_response:
		result := (*Report.ControlCommand)(res).Paras[0].Npara
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
