package syncer

import (
	"github.com/giskook/mdps/conf"
	"time"
)

type Syncer struct {
	ticker *time.Ticker
}

func NewSyncer() *Syncer {
	return &Syncer{
		ticker: time.NewTicker(time.Duration(conf.GetConf().Redis.SyncInterval) * time.Second),
	}
}

func (s *Syncer) Stop() {
	s.ticker.Stop()
}

func (s *Syncer) Do() {
	for {
		select {
		case <-s.ticker.C:
			s.DoWork()
		}
	}
}
