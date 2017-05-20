package db_socket

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/giskook/mdps/base"
	"log"
	"math"
	"time"
)

const (
	TRANS_TABLE_NAME_FMT        string = "DMS_DMP_MONITOR_200601"
	SQL_INSERT_MONITOR_TABLE_EX string = "INSERT %s (MONITOR_ID, DATATYPE, INTBITS, DECIMALBITS, ADDRESS, TIMESTAMP) VALUES(%d, %d, %d, %d, %d, '%s')"
	SQL_INSERT_MONITOR_TABLE    string = "INSERT %s (MONITOR_ID, DATATYPE, DATA, ADDRESS, TIMESTAMP) VALUES(%d, %d, %s, %d, '%s')"
)

func (socket *DbSocket) InsertMonitorsEx(router_db []*base.RouterMonitorDB, router_redis []*base.RouterMonitor) {
	tx, err := socket.Db.Begin()
	base.CheckError(err)
	for _, v := range router_db {
		router_monitor_single, time_stamp := socket.GetValues(router_redis, v)
		if router_monitor_single != nil {
			value, e := socket.ParseValue(v.Datatype, router_monitor_single.DataType, router_monitor_single.DataLen, router_monitor_single.Data)
			if e == nil {
				insert_sql := FmtSQLEx(v, router_monitor_single, time_stamp, value)
				log.Println(insert_sql)
				_, err = tx.Exec(insert_sql)
				base.CheckError(err)

			}
		}
	}
	err = tx.Commit()
	base.CheckError(err)
}

func (socket *DbSocket) InsertMonitors(router_db []*base.RouterMonitorDB, router_redis []*base.RouterMonitor) {
	tx, err := socket.Db.Begin()
	base.CheckError(err)
	for _, v := range router_db {
		router_monitor_single, time_stamp := socket.GetValues(router_redis, v)
		if router_monitor_single != nil {
			insert_sql := FmtSQL(v, router_monitor_single, time_stamp)
			log.Println(insert_sql)
			_, err = tx.Exec(insert_sql)
			base.CheckError(err)

		}
	}
	err = tx.Commit()
	base.CheckError(err)
}

func (socket *DbSocket) GetValues(router_redis []*base.RouterMonitor, router_monitor_db *base.RouterMonitorDB) (*base.RouterMonitorSingle, int64) {
	for _, m := range router_redis {
		if m.RouterID == router_monitor_db.RouterID {
			for _, s := range m.Monitors {
				if s.ModbusAddr == router_monitor_db.ModbusAddr {
					return s, m.TimeStamp
				}
			}
		}
	}

	return nil, 0
}

func (socket *DbSocket) ParseValue(datatype_db uint8, datatype_redis uint8, datalen_redis uint8, data []byte) (*base.Variant, error) {
	err := socket.CheckValue(datatype_db, datatype_redis, datalen_redis)
	if err != nil {
		base.CheckError(err)
		return nil, err
	}

	var data_type uint8
	var data_value uint64
	var data_fvalue float32

	switch datatype_db {
	case base.DATATYPE_DB_SWORD:
		data_type = base.DATATYPE_DB_SWORD
		data_value = uint64(binary.LittleEndian.Uint64(data))
	case base.DATATYPE_DB_UWORD:
		data_type = base.DATATYPE_DB_UWORD
		data_value = uint64(binary.LittleEndian.Uint16(data))
	case base.DATATYPE_DB_SDWORD:
		data_type = base.DATATYPE_DB_SDWORD
		data_value = uint64(binary.LittleEndian.Uint32(data))
	case base.DATATYPE_DB_UDWORD:
		data_type = base.DATATYPE_DB_UDWORD
		data_value = uint64(binary.LittleEndian.Uint32(data))
	case base.DATATYPE_DB_FLOAT:
		data_type = base.DATATYPE_DB_FLOAT
		bits := binary.LittleEndian.Uint64(data)
		data_fvalue = float32(math.Float64frombits(bits))
	}

	return &base.Variant{
		Type:     data_type,
		ValueInt: data_value,
		Float:    data_fvalue,
	}, nil
}

func (socket *DbSocket) CheckValue(datatype_db uint8, datatype_redis uint8, datalen_redis uint8) error {
	datatype := datatype_redis
	datalen := datalen_redis
	switch datatype_db {
	case base.DATATYPE_DB_BIT:
		return errors.New("should not set to bit")
	case base.DATATYPE_DB_SWORD:
		if (datatype == 0 && datalen == 2) ||
			(datatype == 1 && datalen == 1) {
			return nil
		} else {
			return errors.New("db type and redis type do not match")
		}

	case base.DATATYPE_DB_UWORD:
		if (datatype == 0 && datalen == 2) ||
			(datatype == 1 && datalen == 1) {
			return nil
		} else {
			return errors.New("db type and redis type do not match")
		}
	case base.DATATYPE_DB_SDWORD:
		if (datatype == 0 && datalen == 4) ||
			(datatype == 1 && datalen == 2) {
			return nil
		} else {
			return errors.New("db type and redis type do not match")
		}
	case base.DATATYPE_DB_UDWORD:
		if (datatype == 0 && datalen == 4) ||
			(datatype == 1 && datalen == 2) {
			return nil
		} else {
			return errors.New("db type and redis type do not match")
		}
	case base.DATATYPE_DB_FLOAT:
		if (datatype == 0 && datalen == 4) ||
			(datatype == 1 && datalen == 2) {
			return nil
		} else {
			return errors.New("db type and redis type do not match")
		}
	}

	return errors.New("unrecogniced db datatype")
}

func GetTableName(t int64) string {
	if t == 0 {
		t = time.Now().Unix()
	}
	_t := time.Unix(int64(t), 0)
	return _t.Format(TRANS_TABLE_NAME_FMT)
}

func GetTime(t int64) string {
	if t == 0 {
		t = time.Now().Unix()
	}
	_t := time.Unix(int64(t), 0)
	return _t.Format("2006-01-02 15:04:05")
}

func FmtSQLEx(monitor_db *base.RouterMonitorDB, monitor_redis *base.RouterMonitorSingle, time_stamp int64, value *base.Variant) string {
	var insert_sql string
	if value.Type == base.DATATYPE_DB_FLOAT {
		insert_sql = fmt.Sprintf(SQL_INSERT_MONITOR_TABLE_EX, GetTableName(time_stamp), monitor_db.MonitorID, monitor_db.Datatype, int32(value.Float), int32(value.Float)*100%100, monitor_db.ModbusAddr, GetTime(time_stamp))
	} else if value.Type == base.DATATYPE_DB_SWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_MONITOR_TABLE_EX, GetTableName(time_stamp), monitor_db.MonitorID, monitor_db.Datatype, int16(value.ValueInt), 0, monitor_db.ModbusAddr, GetTime(time_stamp))
	} else if value.Type == base.DATATYPE_DB_UWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_MONITOR_TABLE_EX, GetTableName(time_stamp), monitor_db.MonitorID, monitor_db.Datatype, uint16(value.ValueInt), 0, monitor_db.ModbusAddr, GetTime(time_stamp))
	} else if value.Type == base.DATATYPE_DB_SDWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_MONITOR_TABLE_EX, GetTableName(time_stamp), monitor_db.MonitorID, monitor_db.Datatype, int32(value.ValueInt), 0, monitor_db.ModbusAddr, GetTime(time_stamp))
	} else if value.Type == base.DATATYPE_DB_UDWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_MONITOR_TABLE_EX, GetTableName(time_stamp), monitor_db.MonitorID, monitor_db.Datatype, uint32(value.ValueInt), 0, monitor_db.ModbusAddr, GetTime(time_stamp))
	}

	return insert_sql
}

func FmtSQL(monitor_db *base.RouterMonitorDB, monitor_redis *base.RouterMonitorSingle, time_stamp int64) string {
	insert_sql := fmt.Sprintf(SQL_INSERT_MONITOR_TABLE, GetTableName(time_stamp), monitor_db.MonitorID, monitor_db.Datatype, base.GetString(monitor_redis.Data), monitor_db.ModbusAddr, GetTime(time_stamp))

	return insert_sql
}
