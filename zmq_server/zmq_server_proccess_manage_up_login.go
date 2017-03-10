package zmq_server

import (
	"github.com/giskook/mdps/db_socket"
	"github.com/giskook/mdps/pb"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"strconv"
)

func (s *ZmqServer) SendFeedbackLogin(login *ZmqSendValueLogin) {
	tid := strconv.FormatUint(login.Tid, 10)

	s.Socket_Terminal_Manage_Down_Socket.Send(login.Uuid, zmq.SNDMORE)
	s.Socket_Terminal_Manage_Down_Socket.Send(tid, zmq.SNDMORE)
	para := []*Report.Param{
		&Report.Param{
			Type:  Report.Param_UINT8,
			Npara: uint64(login.Check),
		},
	}
	command_rep := &Report.ManageCommand{
		Tid:   login.Tid,
		Type:  Report.ManageCommand_CMT_REP_LOGIN,
		Paras: para,
	}

	data, _ := proto.Marshal(command_rep)
	s.Socket_Terminal_Manage_Down_Socket.Send(string(data), 0)
}

func (s *ZmqServer) ProccessManageUpLogin(command *Report.ManageCommand) {
	uuid := command.Uuid

	tid := command.Tid
	check := db_socket.GetDBSocket().CheckPlcID(tid)
	if check == 1 {
		check = 0
	} else {
		check = 1
	}

	s.CollectSend(&ZmqSendValue{
		SocketType: SOCKET_TERMINAL_MANAGE_DOWN_LOGIN,
		SocketValueLogin: &ZmqSendValueLogin{
			Uuid:  uuid,
			Tid:   tid,
			Check: check,
		},
	})
}
