package syncer

import (
	"github.com/giskook/mdps/db_socket"
	"github.com/giskook/mdps/redis_socket"
)

func (s *Syncer) DoWork() {
	all_monitors_redis := redis_socket.GetRedisSocket().QueryAllMonitors()
	all_monitors_db := db_socket.GetDBSocket().GetMonitors()
	db_socket.GetDBSocket().InsertMonitors(all_monitors_db, all_monitors_redis)
}
