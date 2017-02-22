package protocol

import (
	"bytes"
	"github.com/giskook/mdas_client/base"
)

type ServerHeartPacket struct {
	Tid    uint64
	Serial uint16
	Status uint8
}

func (p *ServerHeartPacket) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, 0,
		PROTOCOL_REQ_HEART, p.Tid, p.Serial)
	writer.WriteByte(p.Status)
	base.WriteWord(&writer, 0)

	base.WriteLength(&writer)

	base.WriteWord(&writer, CalcCRC(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(PROTOCOL_END_FLAG)

	return writer.Bytes()
}
