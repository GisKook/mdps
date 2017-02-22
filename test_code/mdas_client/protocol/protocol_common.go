package protocol

import (
	"bytes"
	"github.com/giskook/mdas_client/base"
	"log"
)

const (
	PROTOCOL_START_FLAG   byte   = 0x55
	PROTOCOL_END_FLAG     byte   = 0xaa
	PROTOCOL_COMMON_LEN   uint16 = 18
	PROTOCOL_MIN_LEN      uint16 = 18
	PROTOCOL_MAX_LEN      uint16 = 1024
	PROTOCOL_TIME_BCD_LEN uint8  = 6

	PROTOCOL_ILLEGAL   uint16 = 254
	PROTOCOL_HALF_PACK uint16 = 255

	PROTOCOL_REQ_REGISTER uint16 = 0x8000
	PROTOCOL_REP_REGISTER uint16 = 0

	PROTOCOL_REQ_LOGIN uint16 = 0x8001
	PROTOCOL_REP_LOGIN uint16 = 0x0001

	PROTOCOL_REQ_HEART uint16 = 0x8002
	PROTOCOL_REP_HEART uint16 = 0x0002

	PROTOCOL_REQ_RESTART uint16 = 0x8003
	PROTOCOL_REP_RESTART uint16 = 0x0003

	PROTOCOL_REQ_GET_SERVER_ADDR uint16 = 0x8004
	PROTOCOL_REP_GET_SERVER_ADDR uint16 = 0x0004

	PROTOCOL_REQ_SET_SERVER_ADDR uint16 = 0x8005
	PROTOCOL_REP_SET_SERVER_ADDR uint16 = 0x0005

	PROTOCOL_REQ_GET_SERIAL_STATUS uint16 = 0x8011
	PROTOCOL_REP_GET_SERIAL_STATUS uint16 = 0x0011

	PROTOCOL_REQ_SET_SERIAL_STATUS uint16 = 0x8012
	PROTOCOL_REP_SET_SERIAL_STATUS uint16 = 0x0012

	PROTOCOL_REQ_DATA_DOWNLOAD uint16 = 0x8007
	PROTOCOL_REP_DATA_DOWNLOAD uint16 = 0x0007

	PROTOCOL_REQ_DATA_QUERY uint16 = 0x8008
	PROTOCOL_REP_DATA_QUERY uint16 = 0x0008

	PROTOCOL_REQ_MONITOR_DOWNLOAD uint16 = 0x8009
	PROTOCOL_REP_MONITOR_DOWNLOAD uint16 = 0x0009

	PROTOCOL_REQ_ALTER_DOWNLOAD uint16 = 0x800a
	PROTOCOL_REP_ALTER_DOWNLOAD uint16 = 0x000a

	PROTOCOL_REP_DATA_MONITOR_UPLOAD uint16 = 0x800b
	PROTOCOL_REP_DATA_ALTER_UPLOAD   uint16 = 0x800c
)

func ParseHeader(buffer []byte) (*bytes.Reader, uint16, uint16, uint64, uint16) {
	reader := bytes.NewReader(buffer)
	reader.Seek(1, 0)
	length := base.ReadWord(reader)
	protocol_id := base.ReadWord(reader)
	tid := base.ReadQuaWord(reader)
	serial := base.ReadWord(reader)

	return reader, length, protocol_id, tid, serial
}

func WriteHeader(writer *bytes.Buffer, length uint16, cmdid uint16, cpid uint64, serial uint16) {
	writer.WriteByte(PROTOCOL_START_FLAG)
	base.WriteWord(writer, length)
	base.WriteWord(writer, cmdid)
	base.WriteQuaWord(writer, cpid)
	base.WriteWord(writer, serial)
}

const POLY_NOMIAL uint16 = 0x8772
const PRESET_VALUE uint16 = 0xFFFF
const CHECK_VALUE uint16 = 0xF0B8

func CalcCRC(buffer []byte, length uint16) uint16 {
	var i uint32
	var j uint32
	var cur_crc_val uint16

	cur_crc_val = PRESET_VALUE
	for i = 0; i < uint32(length); i++ {
		cur_crc_val = cur_crc_val ^ (uint16(buffer[i]))
		for j = 0; j < 8; j++ {
			if (cur_crc_val & 0x0001) != 0 {
				cur_crc_val = (cur_crc_val >> 1) ^ POLY_NOMIAL
			} else {
				cur_crc_val = (cur_crc_val >> 1)
			}
		}
	}

	return (^cur_crc_val) & 0xFFFF
}

func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	log.Println("check protocol")
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return PROTOCOL_ILLEGAL, 0
	}
	if buffer.Bytes()[0] != PROTOCOL_START_FLAG {
		buffer.ReadByte()
		CheckProtocol(buffer)
	} else if bufferlen > 2 {
		pkglen := base.GetWord(buffer.Bytes()[1:3])
		log.Println(pkglen)
		if pkglen < PROTOCOL_MIN_LEN || pkglen > PROTOCOL_MAX_LEN {
			buffer.ReadByte()
			CheckProtocol(buffer)
		}

		if int(pkglen) > bufferlen {
			return PROTOCOL_HALF_PACK, 0
		} else {
			crc_calc := CalcCRC(buffer.Bytes()[1:], uint16(pkglen-4))
			crc_in_protocol := base.GetWord(buffer.Bytes()[pkglen-3 : pkglen-1])
			log.Println(crc_calc)
			log.Println(crc_in_protocol)
			if crc_calc == crc_in_protocol && buffer.Bytes()[pkglen-1] == PROTOCOL_END_FLAG {
				protocol_id := base.GetWord(buffer.Bytes()[3:5])
				return protocol_id, pkglen
			} else {
				log.Println("check error")
				buffer.ReadByte()
				CheckProtocol(buffer)
			}
		}
	} else {
		return PROTOCOL_HALF_PACK, 0
	}

	return PROTOCOL_HALF_PACK, 0
}
