package conn

import (
	"log"
)

func event_handler_server_msg_heart(c *Conn) {
	log.Println("event_handler_server_msg_heart")
	c.UpdateWriteflag()
}
