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
)

func (socket *RedisSocket) ProccessTerminalStatus() {
	defer socket.Mutex_Terminal_Status.Unlock()

	socket.Mutex_Terminal_Status.Lock()
	if len(socket.Terminal_Status) > 0 {
		conn := socket.GetConn()
		defer conn.Close()

		for i, terminal_status := range socket.Terminal_Status {
			conn.Send("HMSET",
				PREFIX_STATUS+
					strconv.FormatUint(terminal_status.Tid, 10),
				STATUS_KEY_UUID,
				terminal_status.Uuid,
				STATUS_KEY_STATUS,
				terminal_status.Status,
				STATUS_KEY_ID,
				terminal_status.Tid,
			)

			socket.Terminal_Status[i] = nil
		}

		socket.Terminal_Status = socket.Terminal_Status[:0]

		conn.Do("")
	}
}
