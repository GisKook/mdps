package zmq_server

import (
	"github.com/giskook/mdps/pb"
)

func (s *ZmqServer) ProccessControlRepRs232GetConfig(command *Report.ControlCommand) {
	tid := command.Tid
	serial := command.SerialNumber

	chan_key := GenerateKey(tid, serial)
	chan_response, ok := GetHttpServer().HttpRespones[chan_key]
	if ok {
		chan_response <- command
	}
}
