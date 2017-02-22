package protocol

import (
	"bytes"
	"github.com/giskook/mdas_client/base"
)

type RegisterPacket struct {
	CpuID  []byte
	Serial uint16
}

func (p *RegisterPacket) Serialize() []byte {
	var writer bytes.Buffer
	WriteHeader(&writer, 0,
		PROTOCOL_REQ_REGISTER, 0, p.Serial)
	base.WriteBytes(&writer, p.CpuID)
	base.WriteLength(&writer)

	base.WriteWord(&writer, CalcCRC(writer.Bytes(), uint16(writer.Len())))
	writer.WriteByte(PROTOCOL_END_FLAG)

	return writer.Bytes()
}
