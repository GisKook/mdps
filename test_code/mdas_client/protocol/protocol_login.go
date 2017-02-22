package protocol

import (
	"bytes"
	"github.com/giskook/mdas_client/base"
)

type LoginPacket struct {
	Tid    uint64
	Serial uint16
}

func (p *LoginPacket) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, 0,
		PROTOCOL_REQ_LOGIN, p.Tid, 59)
	writer.WriteByte(100)
	writer.WriteByte(110)
	base.WriteLength(&writer)

	base.WriteWord(&writer, CalcCRC(writer.Bytes()[1:], uint16(writer.Len())-1))
	writer.WriteByte(PROTOCOL_END_FLAG)

	return writer.Bytes()
}
