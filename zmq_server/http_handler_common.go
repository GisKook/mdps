package zmq_server

import (
	"net/http"
	"strconv"
)

func GetPLCIDAndSerial(r *http.Request) (uint64, uint32) {
	plc_id_string := r.Form.Get(HTTP_PLC_ID)
	serial_string := r.Form.Get(HTTP_PLC_SERIAL)

	plc_id, _ := strconv.ParseUint(plc_id_string, 10, 64)
	serial, _ := strconv.ParseUint(serial_string, 10, 32)

	return plc_id, uint32(serial)
}

func GetUint8Value(r *http.Request, key string) uint8 {
	value := r.Form.Get(key)
	v, _ := strconv.ParseUint(value, 10, 8)

	return uint8(v)
}

func CheckParamters(r *http.Request, keys ...string) bool {
	for _, key := range keys {
		value := r.Form.Get(key)
		if value == "" {
			return false
		}
	}

	return true
}

func GenerateKey(id uint64, serial uint32) uint64 {
	var d_serial uint64
	d_serial = uint64(serial)<<32 + uint64(serial)

	return id ^ d_serial
}
