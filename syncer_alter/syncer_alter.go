package syncer_alter

import (
	"github.com/giskook/mdps/base"
	"github.com/giskook/mdps/conf"
	"log"
	"sync"
	"time"
)

type SyncerAlter struct {
	ticker           *time.Ticker
	RouterAlters     []*base.RouterAlter
	ChanRouterAlters chan *base.RouterAlter
}

func NewSyncerAlter() *SyncerAlter {
	return &SyncerAlter{
		ticker:           time.NewTicker(time.Duration(conf.GetConf().Redis.SyncAlterInterval) * time.Second),
		RouterAlters:     make([]*base.RouterAlter, 0),
		ChanRouterAlters: make(chan *base.RouterAlter),
	}
}

func (s *SyncerAlter) Stop() {
	s.ticker.Stop()
}

func (s *SyncerAlter) get_last_alter(todo *base.RouterAlter) *base.RouterAlter {
	l := len(s.RouterAlters)
	if l > 0 {
		for i := l - 1; i >= 0; i-- {
			if s.RouterAlters[i].MachineID == todo.MachineID &&
				s.RouterAlters[i].RouterID == todo.RouterID &&
				s.RouterAlters[i].ModbusAddr == todo.ModbusAddr {
				return s.RouterAlters[i]

			}
		}
	} else {
		return nil
	}

	return nil

}

func (s *SyncerAlter) Do() {
	for {
		select {
		case <-s.ticker.C:
			s.DoWork()
		case a := <-s.ChanRouterAlters:
			log.Println("------")
			for _, v := range s.RouterAlters {
				log.Println(v)
			}
			log.Println("------")
			log.Println(a)
			log.Println("++++++")
			exist := false
			pre := s.get_last_alter(a)
			log.Println(pre)
			log.Println("////////")
			if pre == nil {
				exist = true
			} else {
				if !(pre.ModbusAddr == a.ModbusAddr &&
					pre.DataType == a.DataType &&
					pre.DataLen == a.DataLen &&
					pre.Status == a.Status) {
					exist = true
				}
			}

			if exist {
				s.RouterAlters = append(s.RouterAlters, a)
			}
		}
	}
}

var G_SyncerAlter *SyncerAlter
var G_MutuxSyncerAlter sync.Mutex

func GetSyncerAlter() *SyncerAlter {
	defer G_MutuxSyncerAlter.Unlock()
	G_MutuxSyncerAlter.Lock()

	if G_SyncerAlter == nil {
		G_SyncerAlter = NewSyncerAlter()
	}

	return G_SyncerAlter
}
