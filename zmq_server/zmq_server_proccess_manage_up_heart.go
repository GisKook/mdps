package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"github.com/giskook/mdps/redis_socket"
)

func (s *ZmqServer) ProccessManageUpHeart(command *Report.ManageCommand) {
	uuid := command.Uuid
	tid := command.Tid
	status := uint8(command.Paras[0].Npara)
	redis_socket.GetRedisSocket().RecvZmqStatus(&redis_socket.TStatus{
		Tid:    tid,
		Uuid:   uuid,
		Status: status,
	})
}
