package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/mdps/pb"
	"github.com/golang/protobuf/proto"
	"log"
)

func (socket *RedisSocket) ProcessDataUploadMonitors() {
	conn := socket.GetConn()
	defer conn.Close()
	conn.Do("SELECT", 1)

	var index int = 0
	var pkg *Report.DataCommand
	for index, pkg = range socket.DataUploadMonitors {
		conn.Send("GET", pkg.Cpid)
	}

	conn.Flush()

	tobe_commit_data_upload_monitors := make([]*Report.DataCommand, 0)
	for i := 0; i < index+1; i++ {
		v_redis, err := conn.Receive()

		if err != nil {
			log.Println(err.Error())
			continue
		}

		v, _ := redis.Bytes(v_redis, nil)

		redis_data_upload_monitors := &Report.DataCommand{}
		err = proto.Unmarshal(v, redis_data_upload_monitors)
		if err != nil {
			log.Println("unmarshal error")
		} else {
			if !proto.Equal(redis_data_upload_monitors, socket.ChargingPiles[i]) {
				if redis_data_upload_monitors.Timestamp < socket.ChargingPiles[i].Timestamp {
					tobe_commit_data_upload_monitors = append(tobe_commit_data_upload_monitors, socket.ChargingPiles[i])
				}
			}
			socket.ChargingPiles[i] = nil
		}
	}

	socket.ChargingPiles = socket.ChargingPiles[:0]

	for _, new_pkg := range tobe_commit_data_upload_monitors {
		data, _ := proto.Marshal(new_pkg)
		conn.Send("SET", new_pkg.Cpid, data)
		new_pkg = nil
	}

	tobe_commit_data_upload_monitors = tobe_commit_data_upload_monitors[:0]

	conn.Flush()
	conn.Do("EXEC")
}
