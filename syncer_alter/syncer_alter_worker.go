package syncer_alter

import (
	"github.com/giskook/mdps/db_socket"
	"log"
)

func (s *SyncerAlter) DoWork() {
	log.Println("fire")

	db_socket.GetDBSocket().GetAlters(s.RouterAlters)
	db_socket.GetDBSocket().InsertAlters(s.RouterAlters)
	s.RouterAlters = s.RouterAlters[:0]
}
