package syncer_alter

import (
	"github.com/giskook/mdps/db_socket"
)

func (s *SyncerAlter) DoWork() {
	db_socket.GetDBSocket().GetAlters(s.RouterAlters)
	db_socket.GetDBSocket().InsertAlters(s.RouterAlters)
	s.RouterAlters = s.RouterAlters[:0]
}
