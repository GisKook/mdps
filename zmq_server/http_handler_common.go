package zmq_server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync/atomic"
)

func GetPLCIDAndSerial(r *http.Request) (uint64, uint32) {
	plc_id_string := r.Form.Get(HTTP_PLC_ID)
	//serial_string := r.Form.Get(HTTP_PLC_SERIAL)

	plc_id, _ := strconv.ParseUint(plc_id_string, 10, 64)
	//serial, _ := strconv.ParseUint(serial_string, 10, 32)
	serial_id := atomic.AddUint32(&GetHttpServer().SerialID, 1)

	return plc_id, serial_id
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
	//	var d_serial uint64
	//	d_serial = uint64(serial)<<32 + uint64(serial)
	//
	//	return id ^ d_serial

	return uint64(serial)
}

type GeneralResponse struct {
	Result uint8  `json:"result"`
	Desc   string `json:"desc"`
}

func EncodingGeneralResponse(result uint8) string {
	general_response := &GeneralResponse{
		Result: result,
		Desc:   HTTP_RESULT[result],
	}

	response, _ := json.Marshal(general_response)

	return string(response)
}
