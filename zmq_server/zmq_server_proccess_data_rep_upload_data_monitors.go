package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"github.com/giskook/mdps/redis_socket"
)

func (s *ZmqServer) ProccessDataRepDataUploadMonitors(command *Report.DataCommand) {
	redis_socket.GetRedisSocket().RecvZmqDataUploadMonitors(command)
	redis_socket.GetRedisSocket().UpdateStatus(command.Tid)
}
