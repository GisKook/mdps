package zmq_server

import (
	"github.com/giskook/mdps/db_socket"
	"github.com/giskook/mdps/pb"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq3"
	"log"
	"strconv"
)

func (s *ZmqServer) ProccessManageUpRegister(command *Report.ManageCommand) {
	log.Println("manage up register")
	uuid := command.Uuid
	tid := command.Tid
	w_c_id := command.Paras[0].Npara*100000 + command.Paras[1].Npara
	s.Socket_Terminal_Manage_Down_Socket.Send(uuid, zmq.SNDMORE)
	s_tid := strconv.FormatUint(tid, 10)
	s.Socket_Terminal_Manage_Down_Socket.Send(s_tid, zmq.SNDMORE)
	s_w_c_id := strconv.FormatUint(w_c_id, 10)
	s.Socket_Terminal_Manage_Down_Socket.Send(s_w_c_id, zmq.SNDMORE)

	log.Println(command.Cpuid)
	plc_id := db_socket.GetDBSocket().GetPlcID(string(command.Cpuid))

	para := []*Report.Param{
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: 0,
		},
		&Report.Param{
			Type:  Report.Param_UINT64,
			Npara: plc_id,
		},
	}
	command_rep := &Report.ManageCommand{
		Type:  Report.ManageCommand_CMT_REP_REGISTER,
		Paras: para,
	}

	data, _ := proto.Marshal(command_rep)
	s.Socket_Terminal_Manage_Down_Socket.Send(string(data), 0)

}
