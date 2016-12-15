package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"time"
)

type TStatus struct {
	Uuid   string
	Status uint8
}

type RedisSocket struct {
	conf               *conf.RedisConf
	Pool               *redis.Pool
	DataUploadMonitors []*Report.DataCommand
	DataUploadAlters   []*Report.DataCommand
	TerminalStatus     map[uint64]*TStatus

	ticker *time.Ticker
}

var G_RedisSocket *RedisSocket = nil

func NewRedisSocket(config *conf.RedisConf) (*RedisSocket, error) {
	if G_RedisSocket == nil {
		G_RedisSocket =
			&RedisSocket{
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
			}
	}
	return G_RedisSocket, nil
}

func GetRedisSocket() *RedisSocket {
	return G_RedisSocket
}
func (socket *RedisSocket) DoWork() {
	defer func() {
		socket.Close()
	}()

	for {
		select {
		case <-socket.ticker.C:
			go socket.ProccessDataUploadMonitors()
			go socket.ProccessDataUploadAlters()
			go socket.ProccessTerminalStatus()
		}
	}
}

func (socket *RedisSocket) GetConn() redis.Conn {
	return socket.Pool.Get()
}

func (socket *RedisSocket) Close() {
	socket.ticker.Stop()
}

func (socket *RedisSocket) RecvZmqDataUploadMonitors(monitors *Report.DataCommand) {
	//	monitors := &Report.DataCommand{}
	//	err := proto.Unmarshal(message, monitors)
	//	if err != nil {
	//		log.Println("unmarshal error")
	//	} else {
	log.Printf("<IN ZMQ>  monitors %s %d \n", monitors.Uuid, monitors.Tid)
	socket.DataUploadMonitors = append(socket.DataUploadMonitors, monitors)
	//	}
}

func (socket *RedisSocket) RecvZmqDataUploadAlters(alters *Report.DataCommand) {
	log.Printf("<IN ZMQ>  alters %s %d \n", alters.Uuid, alters.Tid)
	socket.DataUploadAlters = append(socket.DataUploadAlters, alters)
}

func (socket *RedisSocket) RecvZmqStatus(tid uint64, status *TStatus) {
	socket.TerminalStatus[tid].Uuid = status.Uuid
	socket.TerminalStatus[tid].Status = status.Status
}
