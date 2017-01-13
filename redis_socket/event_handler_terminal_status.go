package redis_socket

import (
	"strconv"
)

const (
	PREFIX_STATUS string = "TSTATE:"
	SEP_STATUS    string = ","

	STATUS_KEY_UUID   string = "uuid"
	STATUS_KEY_STATUS string = "status"
	STATUS_KEY_ID     string = "id"
	//STATUS_KEY_TIME   string = "time"
)

func (socket *RedisSocket) ProccessTerminalStatus() {
	if len(socket.Terminal_Status) > 0 {
		conn := socket.GetConn()
		defer conn.Close()

		for i, terminal_status := range socket.Terminal_Status {
			//			conn.Send("EXPIRE",
			//				PREFIX_STATUS+
			//					strconv.FormatUint(terminal_id, 10),
			//				conf.GetConf().Redis.StatusExpire)
			conn.Send("HMSET",
				PREFIX_STATUS+
					strconv.FormatUint(terminal_status.Tid, 10),
				STATUS_KEY_UUID,
				terminal_status.Uuid,
				STATUS_KEY_STATUS,
				terminal_status.Status,
				STATUS_KEY_ID,
				terminal_status.Tid,
				//			STATUS_KEY_TIME,
				//			strconv.FormatInt(time.Now().Unix(), 10),
			)

			socket.Terminal_Status[i] = nil
		}

		socket.Mutex_Terminal_Status.Lock()
		socket.Terminal_Status = socket.Terminal_Status[:0]
		socket.Mutex_Terminal_Status.Unlock()

		conn.Do("")
	}
}
