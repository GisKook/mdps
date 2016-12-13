package redis_socket

import (
	"bytes"
	"encoding/binary"
	"github.com/giskook/mdps/conf"
	"strconv"
)

const PREFIX_ALTERS string = "TADATA:"
const SEP_ALTERS string = "+"
const SEP_ALTERS_VALUE string = ","

func (socket *RedisSocket) ProccessDataUploadAlters() {
	conn := socket.GetConn()
	defer conn.Close()

	for _, data_command := range socket.DataUploadAlters {
		for _, alter := range data_command.Alters {
			conn.Send("EXPIRE", alter.Id, conf.GetConf().Redis.AltersKeyExpire)
			reader := bytes.NewReader(alter.Data)
			if alter.DataType == 0 {
				var byte_value byte
				for i := uint32(0); i < alter.DataLen; i++ {
					binary.Read(reader, binary.LittleEndian, &byte_value)
					conn.Send("HMSET",
						PREFIX_ALTERS+
							strconv.FormatUint(data_command.Tid, 10),
						strconv.Itoa(int(alter.Id+i))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataType))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataLen)),
						strconv.Itoa(int(byte_value))+
							SEP_ALTERS_VALUE+
							strconv.Itoa(int(alter.Status)))
				}
			} else if alter.DataType == 1 {
				var word_value uint16
				for i := uint32(0); i < alter.DataLen; i++ {
					binary.Read(reader, binary.LittleEndian, &word_value)
					conn.Send("HMSET",
						PREFIX_ALTERS+
							strconv.FormatUint(data_command.Tid, 10),
						strconv.Itoa(int(alter.Id+i*2))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataType))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataLen)),
						strconv.Itoa(int(word_value))+
							SEP_ALTERS_VALUE+
							strconv.Itoa(int(alter.Status)))
				}
			}
		}
	}
	conn.Do("")
}
