package redis_socket

import (
	"github.com/garyburd/redigo/redis"
	"github.com/giskook/mdps/base"
	"github.com/giskook/mdps/pb"
	"strconv"
	"strings"
)

func (socket *RedisSocket) ProccessDataUploadAltersFetch(alters []*Report.DataCommand) []*base.RouterAlterRedis {
	conn := socket.GetConn()
	defer conn.Close()
	router_alters := make([]*base.RouterAlterRedis, 0)

	for _, data_command := range alters {
		alter_key := socket.GenAlterKey(data_command)
		router_alters = append(router_alters, &base.RouterAlterRedis{
			RouterID:   data_command.Tid,
			SerialPort: uint8(data_command.SerialPort),
		})
		conn.Send("HGETALL", alter_key)

	}

	conn.Flush()

	for i := 0; i < len(alters); i++ {
		v_redis, err := conn.Receive()
		if err != nil {
			base.CheckError(err)
			continue
		}
		v, _ := redis.ByteSlices(v_redis, nil)
		socket.PipelineGetAlterValue(v, i, router_alters)
	}

	return router_alters
}

func (socket *RedisSocket) PipelineGetAlterValue(raw [][]byte, index int, router_alters []*base.RouterAlterRedis) {
	alter_count := len(raw)

	for i := 0; i < alter_count; i += 2 {
		key := string(raw[i])
		modbus_datatype_datalen := strings.Split(key, SEP_ALTERS)
		modbus_addr, _ := strconv.Atoi(modbus_datatype_datalen[0])
		datatype, _ := strconv.Atoi(modbus_datatype_datalen[1])
		datalen, _ := strconv.Atoi(modbus_datatype_datalen[2])
		value := raw[i+1]
		status := uint8(value[0])
		data := value[1:]
		router_alters[index].Alters = append(router_alters[index].Alters, &base.Alter{
			ModbusAddr: uint32(modbus_addr),
			DataType:   uint8(datatype),
			DataLen:    uint8(datalen),
			Data:       data,
			Status:     status,
		})
	}
}
