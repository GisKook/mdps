package zmq_server

import (
	"github.com/giskook/mdps/pb"
	"github.com/giskook/mdps/redis_socket"
	//"github.com/golang/protobuf/proto"
	//	zmq "github.com/pebbe/zmq3"
	//"strconv"
)

func (s *ZmqServer) ProccessDataRepDataUploadAlters(command *Report.DataCommand) {
	//log.Println(command)
	redis_socket.GetRedisSocket().RecvZmqDataUploadAlters(command)
	redis_socket.GetRedisSocket().UpdateStatus(command.Tid)
}
