package redis_socket

import (
	"strconv"
)

func (socket *RedisSocket) ProccessDataUploadAlters() {
	defer socket.MutexAlters.Unlock()
	socket.MutexAlters.Lock()
	if len(socket.DataUploadAlters) > 0 {
		conn := socket.GetConn()
		defer conn.Close()

		router_alters := socket.ProccessDataUploadAltersFetch(socket.DataUploadAlters)

		for _, data_command := range socket.DataUploadAlters {
			alter_key := socket.GenAlterKey(data_command)
			conn.Send("DEL", alter_key)
			conn.Send("EXPIRE", alter_key)
			for _, alter := range data_command.Alters {
				if alter.DataType == 0 {
					conn.Send("HMSET",
						alter_key,
						strconv.Itoa(int(alter.ModusAddr))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataType))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataLen)),
						append([]byte{byte(alter.Status)}, alter.Data...))
				} else if alter.DataType == 1 {
					conn.Send("HMSET",
						alter_key,
						strconv.Itoa(int(alter.ModusAddr))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataType))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataLen)),
						append([]byte{byte(alter.Status)}, alter.Data...))
				}
			}
			data_command = nil
		}
		conn.Do("")
		conn.Close()

		socket.DataUploadAlters = socket.DataUploadAlters[:0]
	}
}
