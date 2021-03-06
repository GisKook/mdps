package redis_socket

import (
	"github.com/giskook/mdps/base"
	"github.com/giskook/mdps/pb"
	"log"
)

func (socket *RedisSocket) FilterAlters(alters []*Report.DataCommand, alters_redis []*base.RouterAlterRedis) []*base.RouterAlter {
	report_alters := len(alters)
	if report_alters != len(alters_redis) || report_alters == 0 {
		log.Println("filteralters return nil")
		return nil
	}

	alter := make([]*base.RouterAlter, 0)

	for i, v := range alters {
		alter = append(alter, CompareAlter(v, alters_redis[i].Alters, v.Tid, uint8(v.SerialPort))...)
	}

	return alter
}

func CompareAlter(data_command *Report.DataCommand, alter_redis []*base.Alter, router_id uint64, serial_port uint8) []*base.RouterAlter {
	alter := make([]*base.RouterAlter, 0)

	alter_rep := data_command.Alters
	for _, a := range alter_rep {
		for _, b := range alter_redis {
			if a.ModusAddr == b.ModbusAddr &&
				a.DataType == uint32(b.DataType) &&
				a.DataLen == uint32(b.DataLen) &&
				a.Status != uint32(b.Status) {
				alter = append(alter, &base.RouterAlter{
					RouterID:   router_id,
					SerialPort: serial_port,
					ModbusAddr: a.ModusAddr,
					DataType:   uint8(a.DataType),
					Data:       a.Data,
					Status:     uint8(a.Status),
					//					ModbusAddr: b.ModbusAddr,
					//					DataType:   b.DataType,
					//					Data:       b.Data,
					//					Status:     b.Status,
					TimeStamp: a.Timestamp,
				})
			}
		}
	}

	return alter
}
