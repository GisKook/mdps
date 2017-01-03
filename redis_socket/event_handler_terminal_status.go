package redis_socket

import (
	"github.com/giskook/mdps/conf"
	"strconv"
	"time"
)

const (
	PREFIX_STATUS string = "TSTATE:"
	SEP_STATUS    string = ","

	STATUS_KEY_UUID   string = "uuid"
	STATUS_KEY_STATUS string = "status"
	STATUS_KEY_TIME   string = "time"
)

func (socket *RedisSocket) ProccessTerminalStatus() {
	if len(socket.TerminalStatus) > 0 {
		conn := socket.GetConn()
		defer conn.Close()
		//conn.Do("SELECT", 14)

		for terminal_id := range socket.TerminalStatus {
			conn.Send("EXPIRE",
				PREFIX_STATUS+
					strconv.FormatUint(terminal_id, 10),
				conf.GetConf().Redis.StatusExpire)
			conn.Send("HMSET",
				PREFIX_STATUS+
					strconv.FormatUint(terminal_id, 10),
				STATUS_KEY_UUID,
				socket.TerminalStatus[terminal_id].Uuid,
				STATUS_KEY_STATUS,
				socket.TerminalStatus[terminal_id].Status,
				STATUS_KEY_TIME,
				strconv.FormatInt(time.Now().Unix(), 10),
			)
			delete(socket.TerminalStatus, terminal_id)
		}

		conn.Do("")
	}
}
