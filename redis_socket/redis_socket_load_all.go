package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func PreTreatmentKeys(keys []string) []string {
	for i, key := range keys {
		if !strings.Contains(key, PREFIX_STATUS) {
			keys = append(keys[0:i], keys[i+1:]...)
		}
	}

	return keys
}

func (socket *RedisSocket) LoadAll() {
	conn := socket.GetConn()
	defer func() {
		conn.Close()
		log.Println("end proccess router")
	}()
	var value interface{}
	var cursor_keys []interface{}
	var cursor string = "0"
	var keys []string
	var e error
	for {
		value, e = conn.Do("SCAN", cursor)
		CheckError(e)
		cursor_keys, e = redis.Values(value, e)
		CheckError(e)
		cursor, e = redis.String(cursor_keys[0], nil)
		CheckError(e)
		keys, e = redis.Strings(cursor_keys[1], nil)
		keys = PreTreatmentKeys(keys)
		CheckError(e)
		socket.PipelineGetValue(keys)
		log.Println(cursor)
		log.Println(keys)
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
			log.Println(v)

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
