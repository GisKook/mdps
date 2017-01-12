package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"github.com/giskook/mdps/redis_socket"
)

func (s *ZmqServer) ProccessControlRepBatchAddMonitor(command *Report.ControlCommand) {
	tid := command.Tid
	serial := command.SerialNumber

	chan_key := _GenerateKey(tid, serial)
	chan_response, ok := GetHttpServer().HttpRespones[chan_key]
	if ok {
		chan_response <- command
	}
	redis_socket.GetRedisSocket().UpdateStatus(tid)
}
