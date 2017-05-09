package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"github.com/giskook/mdps/redis_socket"
)

func (s *ZmqServer) ProccessControlRepDataQuery(command *Report.ControlCommand) {
	tid := command.Tid
	serial := command.SerialNumber

	chan_key := _GenerateKey(tid, serial)
	GetHttpServer().DoResponse(&HttpResponsePair{
		Key:     chan_key,
		Command: command,
	})
	redis_socket.GetRedisSocket().UpdateStatus(tid)
}
