package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"log"
)

func (s *ZmqServer) ProccessControlRestart(command *Report.ControlCommand) {
	log.Println("ProccessControlRestart")
	tid := command.Tid
	serial := command.SerialNumber

	chan_key := GenerateKey(tid, serial)
	chan_response, ok := GetHttpServer().HttpRespones[chan_key]
	if ok {
		log.Println("hava key")
		chan_response <- command
	}
}
