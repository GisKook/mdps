package pkg

type Packet interface {
	Serialize() []byte
}
