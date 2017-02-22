package conn

import (
	"github.com/giskook/mdas_client/pkg"
	"github.com/giskook/mdas_client/protocol"
	"log"
)

func event_handler_server_msg_login(c *Conn, p pkg.Packet) {
	log.Println("event_handler_server_msg_login")
	login_pkg := p.(*protocol.ServerLoginPacket)
	c.Status = login_pkg.Result
	if c.Status == 0 {
		go c.heart()
	}
}
