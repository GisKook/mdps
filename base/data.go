package base

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
