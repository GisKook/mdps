package conn

import (
	"github.com/giskook/mdas_client/protocol"
	"log"
)

func event_handler_server_msg_common(conn *Conn) {
	for conn.ReadMore {
		cmdid, pkglen := protocol.CheckProtocol(conn.RecieveBuffer)
		log.Printf("protocol id %d\n", cmdid)

		pkgbyte := make([]byte, pkglen)
		conn.RecieveBuffer.Read(pkgbyte)
		switch cmdid {
		//		case protocol.PROTOCOL_REP_REGISTER:
		//			event_handler_server_msg_register(conn, p)
		//			conn.ReadMore = true
		case protocol.PROTOCOL_REP_LOGIN:
			p := protocol.ParseServerLogin(pkgbyte)
			event_handler_server_msg_login(conn, p)
			conn.ReadMore = true
		case protocol.PROTOCOL_ILLEGAL:
			conn.ReadMore = false
		case protocol.PROTOCOL_HALF_PACK:
			conn.ReadMore = false
		}
	}
}
