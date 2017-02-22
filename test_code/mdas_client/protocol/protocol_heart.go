package protocol

import (
	"bytes"
	"github.com/giskook/mdas_client/base"
)

type HeartPacket struct {
	Tid uint64
}

func (p *HeartPacket) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, 0,
		PROTOCOL_REQ_HEART, p.Tid, 59)
	writer.WriteByte(0)
	base.WriteWord(&writer, 0)
	base.WriteLength(&writer)

	base.WriteWord(&writer, CalcCRC(writer.Bytes()[1:], uint16(writer.Len())-1))
	writer.WriteByte(PROTOCOL_END_FLAG)

	return writer.Bytes()
}
