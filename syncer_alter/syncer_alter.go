package syncer_alter

import (
	"github.com/giskook/mdps/base"
	"github.com/giskook/mdps/conf"
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

func (s *SyncerAlter) Do() {
	for {
		select {
		case <-s.ticker.C:
			s.DoWork()
		case a := <-s.ChanRouterAlters:
			s.RouterAlters = append(s.RouterAlters, a)
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
