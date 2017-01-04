package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"github.com/giskook/mdps/redis_socket"
	"log"
)

func (s *ZmqServer) ProccessManageUpHeart(command *Report.ManageCommand) {
	log.Println("manage up heart")
	uuid := command.Uuid
	tid := command.Tid
	status := uint8(command.Paras[0].Npara)
	log.Println(redis_socket.GetRedisSocket())
	redis_socket.GetRedisSocket().RecvZmqStatus(&redis_socket.TStatus{
		Tid:    tid,
		Uuid:   uuid,
		Status: status,
	})
}
