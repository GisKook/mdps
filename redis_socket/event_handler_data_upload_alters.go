package redis_socket

import (
	//	"bytes"
	//"encoding/binary"
	//	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/base"
	"log"
	"strconv"
)

const PREFIX_ALTERS string = "TADATA:"
const SEP_ALTERS string = "+"
const SEP_ALTERS_VALUE string = ","
const SEP_ALTERS_KEY string = ":"

func (socket *RedisSocket) ProccessDataUploadAlters() {
	defer socket.MutexAlters.Unlock()
	socket.MutexAlters.Lock()
	if len(socket.DataUploadAlters) > 0 {
		log.Println("prcccess data upload alters")
		conn := socket.GetConn()
		defer conn.Close()

		//log.Println(len(socket.DataUploadAlters))
		for _, data_command := range socket.DataUploadAlters {
			alter_key := PREFIX_ALTERS + strconv.FormatUint(data_command.Tid, 10) + SEP_ALTERS_KEY + strconv.FormatUint(uint64(data_command.SerialPort), 10)
			conn.Send("DEL", alter_key)
			conn.Send("EXPIRE", alter_key)
			//log.Println(len(data_command.Alters))
			for _, alter := range data_command.Alters {
				//		reader := bytes.NewReader(alter.Data)
				if alter.DataType == 0 {
					//		var byte_value byte
					//	for i := uint32(0); i < alter.DataLen; i++ {
					//binary.Read(reader, binary.LittleEndian, &byte_value)
					conn.Send("HMSET",
						alter_key,
						strconv.Itoa(int(alter.ModusAddr))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataType))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataLen)),
						//append([]byte{byte(alter.Status)}, alter.Data...))
						base.GetString(append([]byte{byte(alter.Status)}, alter.Data...)))

					//	}
				} else if alter.DataType == 1 {
					//			var word_value uint16
					//		for i := uint32(0); i < alter.DataLen; i++ {
					//binary.Read(reader, binary.LittleEndian, &word_value)
					conn.Send("HMSET",
						alter_key,
						strconv.Itoa(int(alter.ModusAddr))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataType))+
							SEP_ALTERS+
							strconv.Itoa(int(alter.DataLen)),
						append([]byte{byte(alter.Status)}, alter.Data...))
					//base.GetString(append([]byte{byte(alter.Status)}, alter.Data...)))
					//	}
				}
			}
			data_command = nil
		}
		conn.Do("")
		conn.Close()

		socket.DataUploadAlters = socket.DataUploadAlters[:0]
	}
}
