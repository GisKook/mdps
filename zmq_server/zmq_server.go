package zmq_server

import (
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
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

	ExitChan chan struct{}
}

var mutex_server sync.Mutex
var G_ZmqServer *ZmqServer

func (s *ZmqServer) Init(config *conf.ZmqConf) bool {
	s.Socket_Terminal_Manage_Up_Socket, _ = zmq.NewSocket(zmq.PULL)
	log.Printf("Socket_Terminal_Manage_Up_Socket %x\n", s.Socket_Terminal_Manage_Up_Socket)
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

			ExitChan: make(chan struct{}),
		}

	return G_ZmqServer
}

func GetZmqServer() *ZmqServer {
	mutex_server.Lock()
	defer mutex_server.Unlock()
	return G_ZmqServer
}

func (s *ZmqServer) RecvManageUp() {
	for {
		p, _ := s.Socket_Terminal_Manage_Up_Socket.Recv(0)
		log.Println("recv manage up from zmq")
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

func (s *ZmqServer) RecvDataUp() {
	for {
		p, _ := s.Socket_Terminal_Data_Up_Socket.Recv(0)
		log.Println("recv data up from zmq")
		s.Socket_Terminal_Data_Up_Chan <- p
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
	defer func() {
		s.Socket_Terminal_Manage_Up_Socket.Close()
		s.Socket_Terminal_Manage_Down_Socket.Close()
	}()
	log.Println("zmq run")

	go s.RecvManageUp()
	go s.RecvControlUp()
	go s.RecvDataUp()

	for {
		select {
		case <-s.ExitChan:
			return
		case t_m_u := <-s.Socket_Terminal_Manage_Up_Chan:
			log.Println("manage recv chan")
			s.ProccessManageUp(t_m_u)
		case t_c_u := <-s.Socket_Terminal_Control_Up_Chan:
			log.Println("control recv chan")
			s.ProccessControlUp(t_c_u)
		case t_d_u := <-s.Socket_Terminal_Data_Up_Chan:
			log.Println("data recv chan")
			s.ProccessDataUp(t_d_u)
		}
	}
}

func (s *ZmqServer) Stop() {
	close(s.ExitChan)
}

func (s *ZmqServer) ProccessManageUp(p string) {
	log.Println("proccess manage up")
	command := &Report.ManageCommand{}
	err := proto.Unmarshal([]byte(p), command)
	log.Println(command)
	if err != nil {
		log.Println("unmarshal error")
	} else {
		switch command.Type {
		case Report.ManageCommand_CMT_REQ_REGISTER:
			s.ProccessManageUpRegister(command)
		case Report.ManageCommand_CMT_REQ_LOGIN:
			s.ProccessManageUpLogin(command)
		case Report.ManageCommand_CMT_REP_HEART:
			s.ProccessManageUpHeart(command)
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
		case Report.ControlCommand_CMT_REP_DATA_DOWNLOAD:
			s.ProccessControlRepDataDownload(command)
		case Report.ControlCommand_CMT_REP_DATA_QUERY:
			s.ProccessControlRepDataQuery(command)
		case Report.ControlCommand_CMT_REP_BATCH_ADD_MONITOR:
			s.ProccessControlRepBatchAddMonitor(command)
		case Report.ControlCommand_CMT_REP_BATCH_ADD_ALTER:
			s.ProccessControlRepBatchAddAlter(command)
		case Report.ControlCommand_CMT_REP_RS232_GET_CONFIG:
			s.ProccessControlRepRs232GetConfig(command)
		case Report.ControlCommand_CMT_REP_RS232_SET_CONFIG:
			s.ProccessControlRepRs232SetConfig(command)
		case Report.ControlCommand_CMT_REP_GET_SERVER_ADDR:
			s.ProccessControlRepGetServerAddr(command)
		case Report.ControlCommand_CMT_REP_SET_SERVER_ADDR:
			s.ProccessControlRepSetServerAddr(command)
		}
	}
}

func (s *ZmqServer) ProccessDataUp(p string) {
	command := &Report.DataCommand{}
	err := proto.Unmarshal([]byte(p), command)
	if err != nil {
		log.Println(err)
	} else {
		switch command.Type {
		case Report.DataCommand_CMT_REP_DATA_UPLOAD_MONITORS:
			s.ProccessDataRepDataUploadMonitors(command)
		case Report.DataCommand_CMT_REP_DATA_UPLOAD_ALTERS:
			s.ProccessDataRepDataUploadAlters(command)
		}

	}

}
