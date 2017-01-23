package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"github.com/giskook/mdps/redis_socket"
	"log"
)

func (s *ZmqServer) ProccessDataRepDataUploadMonitors(command *Report.DataCommand) {
	log.Println("data up upload monitors")
	//log.Println(command)
	redis_socket.GetRedisSocket().RecvZmqDataUploadMonitors(command)
	redis_socket.GetRedisSocket().UpdateStatus(command.Tid)
}
