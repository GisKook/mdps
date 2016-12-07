package zmq_server

import (
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq3"
	"log"
	"strconv"
	"sync"
)

type ZmqServer struct {
	Socket_Terminal_Manage_Up_Socket   *zmq.Socket
	Socket_Terminal_Manage_Down_Socket *zmq.Socket
	Socket_Terminal_Manage_Up_Chan     chan string
	Socket_Terminal_Manage_Down_Chan   chan string

	Socket_Terminal_Control_Up_Socket   *zmq.Socket
	Socket_Terminal_Control_Down_Socket *zmq.Socket
	Socket_Terminal_Control_Up_Chan     chan string
	Socket_Terminal_Control_Down_Chan   chan string

	Socket_Terminal_Data_Up_Socket *zmq.Socket
	Socket_Terminal_Data_Up_Chan   chan string

	ExitChan  chan struct{}
	waitGroup *sync.WaitGroup
}

var G_ZmqServer *ZmqServer

func (s *ZmqServer) Init(config *conf.ZmqConf) bool {
	s.Socket_Terminal_Manage_Up_Socket, _ = zmq.NewSocket(zmq.PULL)
	s.Socket_Terminal_Manage_Up_Socket.Bind(config.TerminalManageUp)
	s.Socket_Terminal_Manage_Down_Socket, _ = zmq.NewSocket(zmq.PUB)
	s.Socket_Terminal_Manage_Down_Socket.Bind(config.TerminalManageDown)

	s.Socket_Terminal_Control_Up_Socket, _ = zmq.NewSocket(zmq.PULL)
	s.Socket_Terminal_Control_Up_Socket.Bind(config.TerminalControlUp)
	s.Socket_Terminal_Control_Down_Socket, _ = zmq.NewSocket(zmq.PUB)
	s.Socket_Terminal_Control_Down_Socket.Bind(config.TerminalControlDown)

	s.Socket_Terminal_Data_Up_Socket, _ = zmq.NewSocket(zmq.PULL)
	s.Socket_Terminal_Data_Up_Socket.Bind(config.TerminalDataUp)

	return true
}

func NewZmqServer() *ZmqServer {
	G_ZmqServer =
		&ZmqServer{
			Socket_Terminal_Manage_Up_Chan:    make(chan string),
			Socket_Terminal_Manage_Down_Chan:  make(chan string),
			Socket_Terminal_Control_Up_Chan:   make(chan string, 20),
			Socket_Terminal_Control_Down_Chan: make(chan string),
			Socket_Terminal_Data_Up_Chan:      make(chan string),

			ExitChan:  make(chan struct{}),
			waitGroup: &sync.WaitGroup{},
		}

	return G_ZmqServer
}

func GetZmqServer() *ZmqServer {
	return G_ZmqServer
}

func (s *ZmqServer) RecvManageUp() {
	for {
		p, _ := s.Socket_Terminal_Manage_Up_Socket.Recv(0)
		s.Socket_Terminal_Manage_Up_Chan <- p
	}
}

func (s *ZmqServer) RecvControlUp() {
	for {
		p, _ := s.Socket_Terminal_Control_Up_Socket.Recv(0)
		log.Println("recv control up from zmq")
		s.Socket_Terminal_Control_Up_Chan <- p
	}
}

func (s *ZmqServer) SendControlDown(command *Report.ControlCommand) {
	log.Println("SendControlDown")
	uuid := command.Uuid
	s.Socket_Terminal_Control_Down_Socket.Send(uuid, zmq.SNDMORE)

	tid := command.Tid
	s_tid := strconv.FormatUint(tid, 10)
	s.Socket_Terminal_Control_Down_Socket.Send(s_tid, zmq.SNDMORE)

	data, _ := proto.Marshal(command)
	s.Socket_Terminal_Control_Down_Socket.Send(string(data), 0)
}

func (s *ZmqServer) Run() {
	s.waitGroup.Add(1)
	defer func() {
		s.Socket_Terminal_Manage_Up_Socket.Close()
		s.Socket_Terminal_Manage_Down_Socket.Close()
		s.waitGroup.Done()
	}()

	go s.RecvManageUp()
	go s.RecvControlUp()
	for {
		select {
		case <-s.ExitChan:
			return
		case t_m_u := <-s.Socket_Terminal_Manage_Up_Chan:
			s.ProccessManageUp(t_m_u)
		case t_c_u := <-s.Socket_Terminal_Control_Up_Chan:
			log.Println("control recv chan")
			s.ProccessControlUp(t_c_u)
		}
	}
}

func (s *ZmqServer) Stop() {
	close(s.ExitChan)
	s.waitGroup.Wait()
}

func (s *ZmqServer) ProccessManageUp(p string) {
	command := &Report.ManageCommand{}
	err := proto.Unmarshal([]byte(p), command)
	if err != nil {
		log.Println("unmarshal error")
	} else {
		switch command.Type {
		case Report.ManageCommand_CMT_REQ_REGISTER:
			s.ProccessManageUpRegister(command)
		case Report.ManageCommand_CMT_REQ_LOGIN:
			s.ProccessManageUpLogin(command)
		}
	}
}

func (s *ZmqServer) ProccessControlUp(p string) {
	command := &Report.ControlCommand{}
	err := proto.Unmarshal([]byte(p), command)
	if err != nil {
		log.Println("unmarshal error")
	} else {
		switch command.Type {
		case Report.ControlCommand_CMT_REP_RESTART:
			s.ProccessControlRestart(command)
		}
	}
}

//func (s *ZmqServer) ProccessManageUpLogin(command *Report.ManageCommand) {
//	log.Println(command)
//	uuid := command.Uuid
//	s.Socket_Terminal_Manage_Down_Socket.Send(uuid, zmq.SNDMORE)
//
//	tid := command.Tid
//	s_tid := strconv.FormatUint(tid, 10)
//	s.Socket_Terminal_Manage_Down_Socket.Send(s_tid, zmq.SNDMORE)
//
//	para := []*Report.Param{
//		&Report.Param{
//			Type:  Report.Param_UINT8,
//			Npara: 0,
//		},
//	}
//	command_rep := &Report.ManageCommand{
//		Type:  Report.ManageCommand_CMT_REP_LOGIN,
//		Paras: para,
//	}
//
//	data, _ := proto.Marshal(command_rep)
//	s.Socket_Terminal_Manage_Down_Socket.Send(string(data), 0)
//}
//
//func (s *ZmqServer) ProccessManageUpRegister(command *Report.ManageCommand) {
//	log.Println(command)
//	uuid := command.Uuid
//	tid := command.Tid
//	w_c_id := command.Paras[0].Npara*100000 + command.Paras[1].Npara
//	s.Socket_Terminal_Manage_Down_Socket.Send(uuid, zmq.SNDMORE)
//	s_tid := strconv.FormatUint(tid, 10)
//	s.Socket_Terminal_Manage_Down_Socket.Send(s_tid, zmq.SNDMORE)
//	s_w_c_id := strconv.FormatUint(w_c_id, 10)
//	s.Socket_Terminal_Manage_Down_Socket.Send(s_w_c_id, zmq.SNDMORE)
//
//	para := []*Report.Param{
//		&Report.Param{
//			Type:  Report.Param_UINT8,
//			Npara: 0,
//		},
//		&Report.Param{
//			Type:  Report.Param_UINT64,
//			Npara: 1001,
//		},
//	}
//	command_rep := &Report.ManageCommand{
//		Type:  Report.ManageCommand_CMT_REP_REGISTER,
//		Paras: para,
//	}
//
//	data, _ := proto.Marshal(command_rep)
//	s.Socket_Terminal_Manage_Down_Socket.Send(string(data), 0)
//
//}
