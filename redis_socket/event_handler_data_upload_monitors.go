package redis_socket

import (
	//	"bytes"
	//	"encoding/binary"
	"github.com/giskook/mdps/conf"
	"log"
	"strconv"
)

const PREFIX_MONITORS string = "TMDATA:"
const SEP_MONITORS string = "+"

func (socket *RedisSocket) ProccessDataUploadMonitors() {
	log.Println("prcccess data upload monitors")
	conn := socket.GetConn()
	defer conn.Close()

	log.Println(len(socket.DataUploadMonitors))
	for _, data_command := range socket.DataUploadMonitors {
		conn.Send("EXPIRE", PREFIX_MONITORS+
			strconv.FormatUint(data_command.Tid, 10),
			conf.GetConf().Redis.MonitorsKeyExpire)

		log.Println(len(data_command.Monitors))
		for i, monitor := range data_command.Monitors {
			//reader := bytes.NewReader(monitor.Data)
			if monitor.DataType == 0 {
				//	var byte_value byte
				//			for i := 0; i < int(monitor.DataLen); i++ {
				//	binary.Read(reader, binary.LittleEndian, &byte_value)
				conn.Send("HMSET",
					PREFIX_MONITORS+
						strconv.FormatUint(data_command.Tid, 10),
					strconv.Itoa(int(monitor.Id)+i)+
						SEP_MONITORS+
						strconv.Itoa(int(monitor.DataType))+
						SEP_MONITORS+
						strconv.Itoa(int(monitor.DataLen)),
					monitor.Data)
				//		}
			} else if monitor.DataType == 1 {
				//		var word_value uint16
				//	for i := 0; i < int(monitor.DataLen); i++ {
				//			binary.Read(reader, binary.LittleEndian, &word_value)
				conn.Send("HMSET",
					PREFIX_MONITORS+
						strconv.FormatUint(data_command.Tid, 10),
					strconv.Itoa(int(monitor.Id)+i*2)+
						SEP_MONITORS+
						strconv.Itoa(int(monitor.DataType))+
						SEP_MONITORS+
						strconv.Itoa(int(monitor.DataLen)),
					monitor.Data)
				//}
			}
		}
		data_command = nil
	}
	socket.DataUploadMonitors = socket.DataUploadMonitors[:0]
	conn.Do("")
	conn.Close()
}
