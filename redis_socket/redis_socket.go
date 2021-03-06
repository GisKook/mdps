package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"sync"
	"time"
)

const (
	TERMINAL_STATUS_ONLINE  uint8 = 0
	TERMINAL_STATUS_OFFLINE uint8 = 1
	TERMINAL_STATUS_TT      uint8 = 2
	TERMINAL_STATUS_KEEP    uint8 = 3
)

type TStatus struct {
	Uuid      string
	Tid       uint64
	Status    uint8
	Timestamp int64
}

type RedisSocket struct {
	conf *conf.RedisConf
	Pool *redis.Pool

	DataUploadMonitors []*Report.DataCommand
	MutexMonitors      sync.Mutex
	DataUploadAlters   []*Report.DataCommand
	MutexAlters        sync.Mutex

	Mutex_Terminal_Status sync.Mutex
	Terminal_Status       []*TStatus
	Terminal_Status_Chan  chan *TStatus

	ticker *time.Ticker
}

var G_Mutex_RedisSocket sync.Mutex
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
				DataUploadMonitors:   make([]*Report.DataCommand, 0),
				DataUploadAlters:     make([]*Report.DataCommand, 0),
				ticker:               time.NewTicker(time.Duration(config.TranInterval) * time.Second),
				Terminal_Status:      make([]*TStatus, 0),
				Terminal_Status_Chan: make(chan *TStatus),
			}
	}
	return G_RedisSocket, nil
}

func GetRedisSocket() *RedisSocket {
	defer G_Mutex_RedisSocket.Unlock()
	G_Mutex_RedisSocket.Lock()
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
			go GetStatusChecker().Check()
		case p := <-socket.Terminal_Status_Chan:
			socket.Mutex_Terminal_Status.Lock()
			if p.Status != TERMINAL_STATUS_KEEP {
				socket.Terminal_Status = append(socket.Terminal_Status, p)
			}
			if p.Status != TERMINAL_STATUS_OFFLINE {
				GetStatusChecker().Insert(p.Tid, time.Now().Unix())
			}
			socket.Mutex_Terminal_Status.Unlock()
		}
	}
}

func (socket *RedisSocket) GetConn() redis.Conn {
	return socket.Pool.Get()
}

func (socket *RedisSocket) Close() {
	close(socket.Terminal_Status_Chan)
	socket.ticker.Stop()
}

func (socket *RedisSocket) RecvZmqDataUploadMonitors(monitors *Report.DataCommand) {
	socket.MutexMonitors.Lock()
	socket.DataUploadMonitors = append(socket.DataUploadMonitors, monitors)
	socket.MutexMonitors.Unlock()
}

func (socket *RedisSocket) RecvZmqDataUploadAlters(alters *Report.DataCommand) {
	socket.MutexAlters.Lock()
	socket.DataUploadAlters = append(socket.DataUploadAlters, alters)
	socket.MutexAlters.Unlock()
}

func (socket *RedisSocket) RecvZmqStatus(status *TStatus) {
	socket.Terminal_Status_Chan <- status
}

func (socket *RedisSocket) UpdateStatus(tid uint64) {
	socket.RecvZmqStatus(&TStatus{
		Tid: tid,
		//	Status: TERMINAL_STATUS_KEEP,
		Status: TERMINAL_STATUS_ONLINE,
	})
}
