package redis_socket

import (
	//	"bytes"
	//"encoding/binary"
	"github.com/giskook/mdps/conf"
	"log"
	"strconv"
)

const PREFIX_ALTERS string = "TADATA:"
const SEP_ALTERS string = "+"
const SEP_ALTERS_VALUE string = ","

func (socket *RedisSocket) ProccessDataUploadAlters() {
	log.Println("prcccess data upload alters")
	conn := socket.GetConn()
	defer conn.Close()

	log.Println(len(socket.DataUploadAlters))
	for i, data_command := range socket.DataUploadAlters {
		log.Println(len(data_command.Alters))
		for _, alter := range data_command.Alters {
			conn.Send("EXPIRE", PREFIX_ALTERS+
				strconv.FormatUint(data_command.Tid, 10),
				conf.GetConf().Redis.AltersKeyExpire)
			//		reader := bytes.NewReader(alter.Data)
			if alter.DataType == 0 {
				//		var byte_value byte
				//	for i := uint32(0); i < alter.DataLen; i++ {
				//binary.Read(reader, binary.LittleEndian, &byte_value)
				conn.Send("HMSET",
					PREFIX_ALTERS+
						strconv.FormatUint(data_command.Tid, 10),
					strconv.Itoa(int(alter.Id)+i)+
						SEP_ALTERS+
						strconv.Itoa(int(alter.DataType))+
						SEP_ALTERS+
						strconv.Itoa(int(alter.DataLen)),
					append([]byte{byte(alter.Status)}, alter.Data...))

				//	}
			} else if alter.DataType == 1 {
				//			var word_value uint16
				//		for i := uint32(0); i < alter.DataLen; i++ {
				//binary.Read(reader, binary.LittleEndian, &word_value)
				conn.Send("HMSET",
					PREFIX_ALTERS+
						strconv.FormatUint(data_command.Tid, 10),
					strconv.Itoa(int(alter.Id)+i*2)+
						SEP_ALTERS+
						strconv.Itoa(int(alter.DataType))+
						SEP_ALTERS+
						strconv.Itoa(int(alter.DataLen)),
					append([]byte{byte(alter.Status)}, alter.Data...))
				//	}
			}
		}
		data_command = nil
	}
	conn.Do("")
	conn.Close()

	socket.DataUploadAlters = socket.DataUploadAlters[:0]
}
