package redis_socket

import (
	"bytes"
	"encoding/binary"
	"github.com/giskook/mdps/conf"
	"strconv"
)

const PREFIX_MONITORS string = "TMDATA:"
const SEP_MONITORS string = "+"

func (socket *RedisSocket) ProccessDataUploadMonitors() {
	conn := socket.GetConn()
	defer conn.Close()

	for _, data_command := range socket.DataUploadMonitors {
		for _, monitor := range data_command.Monitors {
			conn.Send("EXPIRE", monitor.Id, conf.GetConf().Redis.MonitorsKeyExpire)
			reader := bytes.NewReader(monitor.Data)
			if monitor.DataType == 0 {
				var byte_value byte
				for i := uint32(0); i < monitor.DataLen; i++ {
					binary.Read(reader, binary.LittleEndian, &byte_value)
					conn.Send("HMSET",
						PREFIX_MONITORS+
							strconv.FormatUint(data_command.Tid, 10),
						strconv.Itoa(int(monitor.Id+i))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataType))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataLen)),
						byte_value)
				}
			} else if monitor.DataType == 1 {
				var word_value uint16
				for i := uint32(0); i < monitor.DataLen; i++ {
					binary.Read(reader, binary.LittleEndian, &word_value)
					conn.Send("HMSET",
						PREFIX_MONITORS+
							strconv.FormatUint(data_command.Tid, 10),
						strconv.Itoa(int(monitor.Id+i*2))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataType))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataLen)),
						word_value)
				}
			}
		}
	}
	conn.Do("")
}
