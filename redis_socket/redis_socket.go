package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

type RedisSocket struct {
	conf               *conf.RedisConfigure
	Pool               *redis.Pool
	DataUploadMonitors []*Report.DataCommand
	DataUploadAlters   []*Report.DataCommand

	ticker *time.Ticker
}

func NewRedisSocket(config *conf.RedisConfigure) (*RedisSocket, error) {
	return &RedisSocket{
		conf: config,
		Pool: &redis.Pool{
			MaxIdle:     config.MaxIdle,
			IdleTimeout: time.Duration(config.IdleTimeOut) * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", config.Addr)
				if err != nil {
					log.Println(err.Error())
					return nil, err
				}

				if len(config.Passwd) > 0 {
					if _, err := c.Do("AUTH", config.Passwd); err != nil {
						log.Println(err.Error())
						c.Close()
						return nil, err
					}
				}

				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}

				_, err := c.Do("PING")

				return err
			},
		},
		DataUploadMonitors: make([]*Report.DataCommand, 0),
		DataUploadAlters:   make([]*Report.DataCommand, 0),
		ticker:             time.NewTicker(time.Duration(config.TranInterval) * time.Second),
	}, nil
}

func (socket *RedisSocket) DoWork() {
	defer func() {
		socket.Close()
	}()

	for {
		select {
		case <-socket.ticker.C:
			go socket.ProcessDataUploadMonitors()
		}
	}
}

func (socket *RedisSocket) GetConn() redis.Conn {
	return socket.Pool.Get()
}

func (socket *RedisSocket) Close() {
	socket.ticker.Stop()
}

func (socket *RedisSocket) RecvNsqDataUpload(message []byte) {
	data_upload := &Report.DataCommand{}
	err := proto.Unmarshal(message, data_upload)
	if err != nil {
		log.Println("unmarshal error")
	} else {
		log.Printf("<IN NSQ> %s %d \n", data_upload.DasUuid, data_upload.Cpid)
		socket.DataUpload = append(socket.DataUpload, data_upload)
	}
}
