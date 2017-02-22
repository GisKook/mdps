package conn

import (
	"time"
)

type Conn_status struct {
	conn   *Conn
	status uint8 // 0 normal 1 error
}

type Conns struct {
	CS     map[uint64]*Conn_status
	ticker *time.Ticker //time.NewTicker(time.Duration(5) * time.Second)
}

var connsInstance *Conns

func GetConns() *Conns {
	if connsInstance == nil {
		connsInstance = &Conns{
			CS:     make(map[uint64]*Conn_status),
			ticker: time.NewTicker(time.Duration(5) * time.Second),
		}
	}

	return connsInstance
}
