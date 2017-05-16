package base

const (
	DATATYPE_DB_BIT    uint8 = 0
	DATATYPE_DB_SWORD  uint8 = 1
	DATATYPE_DB_UWORD  uint8 = 2
	DATATYPE_DB_SDWORD uint8 = 3
	DATATYPE_DB_UDWORD uint8 = 4
	DATATYPE_DB_FLOAT  uint8 = 5

	DATATYPE_REDIS_BYTE uint8 = 0
	DATATYPE_REDIS_WORD uint8 = 1
)

type Variant struct {
	Type     uint8
	ValueInt uint64
	Float    float32
}

type RouterMonitorSingle struct {
	ModbusAddr uint32
	DataType   uint8
	DataLen    uint8
	Data       []byte
}

type RouterMonitor struct {
	RouterID   uint32
	SerialPort uint8
	Monitors   []*RouterMonitorSingle
	TimeStamp  int64
}

type RouterMonitorDB struct {
	RouterID   uint32
	MonitorID  uint32
	Datatype   uint8
	ModbusAddr uint32
}

type Alter struct {
	ModbusAddr uint32
	DataType   uint8
	DateLen    uint8
	Data       []byte
	Status     uint8
}

type RouterAlterRedis struct {
	RouterID   uint64
	SerialPort uint8
	Alters     []Alter
}

type RouterAlterDB struct {
	RouterID   uint64
	SerialPort uint8
	ModbusAddr uint32
	DataType   uint8
	DateLen    uint8
	Data       []byte
	Status     uint8
}
