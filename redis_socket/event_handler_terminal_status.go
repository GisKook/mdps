package redis_socket

import (
	"github.com/giskook/mdps/conf"
	"strconv"
)

const (
	PREFIX_STATUS string = "TSTATE:"
	SEP_STATUS    string = ","

	STATUS_KEY_UUID   string = "uuid"
	STATUS_KEY_STATUS string = "status"
)

func (socket *RedisSocket) ProccessTerminalStatus() {
	conn := socket.GetConn()
	defer conn.Close()
	conn.Do("SELECT", 14)

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
			socket.TerminalStatus[terminal_id].Status)
	}

	conn.Do("")
}
