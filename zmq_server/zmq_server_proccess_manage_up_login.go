package zmq_server

import (
	"github.com/giskook/mdps/db_socket"
	"github.com/giskook/mdps/pb"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq3"
	"log"
	"strconv"
)

func (s *ZmqServer) ProccessManageUpLogin(command *Report.ManageCommand) {
	log.Println(command)
	uuid := command.Uuid
	s.Socket_Terminal_Manage_Down_Socket.Send(uuid, zmq.SNDMORE)

	tid := command.Tid
	s_tid := strconv.FormatUint(tid, 10)
	s.Socket_Terminal_Manage_Down_Socket.Send(s_tid, zmq.SNDMORE)
	log.Println("proccess login")
	check := db_socket.GetDBSocket().CheckPlcID(122)
	log.Println(check)
	if check == 1 {
		check = 0
	} else {
		check = 1
	}

	para := []*Report.Param{
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(check),
		},
	}
	command_rep := &Report.ManageCommand{
		Type:  Report.ManageCommand_CMT_REP_LOGIN,
		Paras: para,
	}

	data, _ := proto.Marshal(command_rep)
	s.Socket_Terminal_Manage_Down_Socket.Send(string(data), 0)
}
