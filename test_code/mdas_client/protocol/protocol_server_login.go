package protocol

import (
	"github.com/giskook/mdas_client/base"
)

type ServerLoginPacket struct {
	Tid       uint64
	Result    uint8
	TimeStamp uint64
	Heart     uint32
	Reserve   uint32
}

func (p *ServerLoginPacket) Serialize() []byte {
	return nil
}

func ParseServerLogin(buffer []byte) *ServerLoginPacket {
	reader, _, _, tid, _ := ParseHeader(buffer)
	result, _ := reader.ReadByte()
	time_stamp := base.ReadQuaWord(reader)
	heart := base.ReadDWord(reader)
	reserve := base.ReadDWord(reader)

	return &ServerLoginPacket{
		Tid:       tid,
		Result:    result,
		TimeStamp: time_stamp,
		Heart:     heart,
		Reserve:   reserve,
	}
}
