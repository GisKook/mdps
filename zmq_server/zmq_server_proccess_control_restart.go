package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"log"
)

func (s *ZmqServer) ProccessControlRestart(command *Report.ControlCommand) {
	log.Println("proccess control restart")

	tid := command.Tid
	serial := command.SerialNumber

	chan_key := GenerateKey(tid, serial)
	_, ok := GetHttpServer().HttpRespones[chan_key]
	if ok {
		GetHttpServer().HttpRespones[tid] <- command
	}
}
