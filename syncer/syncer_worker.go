package syncer

import (
	"github.com/giskook/mdps/redis_socket"
	"log"
)

func (s *Syncer) DoWork() {
	all_monitors := redis_socket.GetRedisSocket().QueryAllMonitors()
	log.Println(all_monitors)
}
