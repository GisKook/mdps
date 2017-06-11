package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/mdps/base"
	"log"
	"strconv"
	"time"
)

func (socket *RedisSocket) LoadAll() {
	conn := socket.GetConn()
	defer func() {
		conn.Close()
	}()
	var value interface{}
	var cursor_keys []interface{}
	var cursor string = "0"
	var keys []string
	var e error
	for {
		value, e = conn.Do("SCAN", cursor)
		base.CheckError(e)
		cursor_keys, e = redis.Values(value, e)
		base.CheckError(e)
		cursor, e = redis.String(cursor_keys[0], nil)
		base.CheckError(e)
		keys, e = redis.Strings(cursor_keys[1], nil)
		keys = base.FilterStringArray(keys, PREFIX_STATUS)
		base.CheckError(e)
		socket.PipelineGetValue(keys)
		if cursor == "0" {
			return
		}
	}
}

func (socket *RedisSocket) PipelineGetValue(keys []string) {
	if len(keys) != 0 {
		conn := socket.GetConn()
		defer func() {
			conn.Close()
		}()

		var index int = 0
		var key string = ""
		for index, key = range keys {
			conn.Send("HMGET", key, STATUS_KEY_STATUS, STATUS_KEY_ID)
		}

		conn.Flush()

		for i := 0; i < index+1; i++ {
			v_redis, err := conn.Receive()

			if err != nil {
				log.Println(err)
				continue
			}

			v, _ := redis.Strings(v_redis, nil)

			if err != nil {
				log.Println("unmarshal error PipelineGetValue")
			} else {
				if v[0] == "0" {
					id, _ := strconv.ParseUint(v[1], 10, 64)
					GetStatusChecker().Insert(id, time.Now().Unix())
				}
			}
		}
		conn.Do("")
	}
}
