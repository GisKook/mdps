package zmq_server

import (
	"github.com/giskook/mdps/pb"
	//"github.com/golang/protobuf/proto"
	//zmq "github.com/pebbe/zmq3"
	"log"
	//"strconv"
)

func (s *ZmqServer) ProccessDataRepDataUploadMonitors(command *Report.DataCommand) {
	log.Println("data up upload monitors")
	//	uuid := command.Uuid
	//	tid := command.Tid
}
