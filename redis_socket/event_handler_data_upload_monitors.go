package redis_socket

import (
	"strconv"
	"time"
)

const (
	PREFIX_MONITORS  string = "TMDATA:"
	SEP_MONITORS     string = "+"
	SEP_KEY_MONITORS string = ":"
	TIMESTAMP        string = "TIMESTAMP"
)

func (socket *RedisSocket) ProccessDataUploadMonitors() {
	defer socket.MutexMonitors.Unlock()
	socket.MutexMonitors.Lock()
	if len(socket.DataUploadMonitors) > 0 {
		conn := socket.GetConn()
		defer conn.Close()

		for _, data_command := range socket.DataUploadMonitors {
			monitor_key :=
				PREFIX_MONITORS +
					strconv.FormatUint(data_command.Tid, 10) +
					SEP_KEY_MONITORS +
					strconv.FormatUint(uint64(data_command.SerialPort), 10)
			conn.Send("DEL", monitor_key)
			conn.Send("EXPIRE", monitor_key)
			conn.Send("HMSET", TIMESTAMP, time.Now().Unix())

			for _, monitor := range data_command.Monitors {
				if monitor.DataType == 0 {
					conn.Send("HMSET",
						monitor_key,
						strconv.Itoa(int(monitor.ModusAddr))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataType))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataLen)),
						monitor.Data)
				} else if monitor.DataType == 1 {
					conn.Send("HMSET",
						monitor_key,
						strconv.Itoa(int(monitor.ModusAddr))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataType))+
							SEP_MONITORS+
							strconv.Itoa(int(monitor.DataLen)),
						monitor.Data)
				}
			}
			data_command = nil
		}
		socket.DataUploadMonitors = socket.DataUploadMonitors[:0]
		conn.Do("")
		conn.Close()
	}
}
