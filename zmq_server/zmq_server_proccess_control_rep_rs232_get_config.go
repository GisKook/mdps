package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"github.com/giskook/mdps/redis_socket"
	"log"
)

func (s *ZmqServer) ProccessControlRepRs232GetConfig(command *Report.ControlCommand) {
	tid := command.Tid
	serial := command.SerialNumber

	chan_key := _GenerateKey(tid, serial)
	chan_response, ok := GetHttpServer().HttpRespones[chan_key]
	if ok {
		chan_response <- command
		log.Println("ProccessControlRepRs232GetConfig have chan")
	} else {
		log.Println("ProccessControlRepRs232GetConfig do not have chan")
	}
	redis_socket.GetRedisSocket().UpdateStatus(tid)
}
