package redis_socket

import (
	//	"bytes"
	//	"encoding/binary"
	//	"github.com/giskook/mdps/conf"
	//"github.com/giskook/mdps/base"
	//"log"
	"strconv"
)

const PREFIX_MONITORS string = "TMDATA:"
const SEP_MONITORS string = "+"
const SEP_KEY_MONITORS string = ":"

func (socket *RedisSocket) ProccessDataUploadMonitors() {
	defer socket.MutexMonitors.Unlock()
	socket.MutexMonitors.Lock()
	if len(socket.DataUploadMonitors) > 0 {
		conn := socket.GetConn()
		defer conn.Close()

		//log.Println(len(socket.DataUploadMonitors))
		for _, data_command := range socket.DataUploadMonitors {
			monitor_key :=
				PREFIX_MONITORS +
					strconv.FormatUint(data_command.Tid, 10) +
					SEP_KEY_MONITORS +
					strconv.FormatUint(uint64(data_command.SerialPort), 10)
			conn.Send("DEL", monitor_key)
			conn.Send("EXPIRE", monitor_key)

			//log.Println(len(data_command.Monitors))
			for _, monitor := range data_command.Monitors {
				//log.Println(monitor)
				//reader := bytes.NewReader(monitor.Data)
				if monitor.DataType == 0 {
					//	var byte_value byte
					//			for i := 0; i < int(monitor.DataLen); i++ {
					//	binary.Read(reader, binary.LittleEndian, &byte_value)
					conn.Send("HMSET",
						monitor_key,
						strconv.Itoa(int(monitor.ModusAddr))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataType))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataLen)),
						monitor.Data)
					//base.GetString(monitor.Data))
					//		}
				} else if monitor.DataType == 1 {
					//		var word_value uint16
					//	for i := 0; i < int(monitor.DataLen); i++ {
					//			binary.Read(reader, binary.LittleEndian, &word_value)
					conn.Send("HMSET",
						monitor_key,
						strconv.Itoa(int(monitor.ModusAddr))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataType))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataLen)),
						monitor.Data)
					//base.GetString(monitor.Data))
					//}
				}
			}
			data_command = nil
		}
		socket.DataUploadMonitors = socket.DataUploadMonitors[:0]
		conn.Do("")
		conn.Close()
	}
}
