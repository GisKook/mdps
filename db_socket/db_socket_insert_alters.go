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
	TRANS_TABLE_ALTER_NAME_FMT string = "DMS_DAP_ALTER_200601"
	SQL_INSERT_ALTER_TABLE_EX  string = "INSERT %s (ALERT_ID, DATATYPE, INTBITS, DECIMALBITS, ADDRESS, STATUS) VALUES(%d, %d, %d, %d, %d, %d)"
	SQL_INSERT_ALTER_TABLE     string = "INSERT %s (ROUTER_ID, MACHINEID, ALERT_ID, DATATYPE, DATA, ADDRESS, STATUS) VALUES(%d, %s, %d, %d, %s, %d, %d)"
	SQL_CREATE_ALTER_TABLE_FMT string = "CREATE TABLE %s ( ID int(11) NOT NULL AUTO_INCREMENT, ROUTERID int , MACHINEID varchar(50), ALERT_ID int(11) NOT NULL, DATATYPE int(11) NOT NULL, INTBITS int(11) DEFAULT NULL, DECIMALBITS int(11) DEFAULT NULL, ADDRESS int(11) NOT NULL, TIMESTAMP timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, isread int(11) DEFAULT '0', DATA varchar(8) DEFAULT NULL, STATUS int(11) DEFAULT NULL, PRIMARY KEY (ID), KEY ALERT_ID (ALERT_ID)) ENGINE=MyISAM AUTO_INCREMENT=18 DEFAULT CHARSET=utf8"
)

func (socket *DbSocket) create_alter_table(table_name string) bool {
	sql_create_alter_table := fmt.Sprintf(SQL_CREATE_ALTER_TABLE_FMT, table_name)
	log.Println(sql_create_alter_table)
	_, err := socket.Db.Exec(sql_create_alter_table)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (socket *DbSocket) InsertAlters(router_alter []*base.RouterAlter) {
	tx, err := socket.Db.Begin()
	base.CheckError(err)
	table_name := GetAlterTableName(time.Now().Unix())
	if !socket.check_table_if_exist(table_name) {
		if !socket.create_alter_table(table_name) {
			return
		}
	}
	router_done := make([]*base.RouterAlter, 0)
	check_exist := func(done []*base.RouterAlter, todo *base.RouterAlter) bool {
		for _, v := range done {
			if v.RouterID == todo.RouterID &&
				v.SerialPort == todo.SerialPort &&
				v.ModbusAddr == todo.ModbusAddr &&
				v.DataType == todo.DataType &&
				v.DataLen == todo.DataLen &&
				v.Status == todo.Status &&
				v.DataTypeDB == todo.DataTypeDB &&
				v.AlterIDDB == todo.AlterIDDB {
				return true
			}
		}

		return false
	}

	for _, v := range router_alter {
		if !check_exist(router_done, v) {
			insert_sql := FmtInsertAlterSQL(table_name, v)
			log.Println(insert_sql)
			_, err = tx.Exec(insert_sql)
			base.CheckError(err)
		}

	}
	err = tx.Commit()
	base.CheckError(err)
}
func (socket *DbSocket) InsertAltersEx(router_alter []*base.RouterAlter) {
	tx, err := socket.Db.Begin()
	base.CheckError(err)
	for _, v := range router_alter {
		value, e := socket.ParseAlterValue(v)
		if e == nil {
			insert_sql := FmtInsertAlterSQLEx(v, value)
			log.Println(insert_sql)
			_, err = tx.Exec(insert_sql)
			base.CheckError(err)

		}
	}
	err = tx.Commit()
	base.CheckError(err)
}

func (socket *DbSocket) ParseAlterValue(router *base.RouterAlter) (*base.Variant, error) {
	err := socket.CheckAlterValue(router)
	if err != nil {
		base.CheckError(err)
		return nil, err
	}

	var data_type uint8
	var data_value uint64
	var data_fvalue float32

	switch router.DataTypeDB {
	case base.DATATYPE_DB_SWORD:
		data_type = base.DATATYPE_DB_SWORD
		data_value = uint64(binary.LittleEndian.Uint64(router.Data))
	case base.DATATYPE_DB_UWORD:
		data_type = base.DATATYPE_DB_UWORD
		data_value = uint64(binary.LittleEndian.Uint16(router.Data))
	case base.DATATYPE_DB_SDWORD:
		data_type = base.DATATYPE_DB_SDWORD
		data_value = uint64(binary.LittleEndian.Uint32(router.Data))
	case base.DATATYPE_DB_UDWORD:
		data_type = base.DATATYPE_DB_UDWORD
		data_value = uint64(binary.LittleEndian.Uint32(router.Data))
	case base.DATATYPE_DB_FLOAT:
		data_type = base.DATATYPE_DB_FLOAT
		bits := binary.LittleEndian.Uint64(router.Data)
		data_fvalue = float32(math.Float64frombits(bits))
	}

	return &base.Variant{
		Type:     data_type,
		ValueInt: data_value,
		Float:    data_fvalue,
	}, nil
}

func (socket *DbSocket) CheckAlterValue(router *base.RouterAlter) error {
	datatype := router.DataType
	datalen := router.DataLen
	switch router.DataTypeDB {
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

func GetAlterTableName(t int64) string {
	if t == 0 {
		t = time.Now().Unix()
	}
	_t := time.Unix(int64(t), 0)
	return _t.Format(TRANS_TABLE_ALTER_NAME_FMT)
}

//func GetTime(t int64) string {
//	if t == 0 {
//		t = time.Now().Unix()
//	}
//	_t := time.Unix(int64(t), 0)
//	return _t.Format("2006-01-02 15:04:05")
//}

func FmtInsertAlterSQLEx(router *base.RouterAlter, value *base.Variant) string {
	time_stamp := time.Now().Unix()
	var insert_sql string
	if value.Type == base.DATATYPE_DB_FLOAT {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE_EX, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, int32(value.Float), int32(value.Float*100)%100, router.ModbusAddr)
	} else if value.Type == base.DATATYPE_DB_SWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE_EX, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, int16(value.ValueInt), 0, router.ModbusAddr)
	} else if value.Type == base.DATATYPE_DB_UWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE_EX, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, uint16(value.ValueInt), 0, router.ModbusAddr)
	} else if value.Type == base.DATATYPE_DB_SDWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE_EX, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, int32(value.ValueInt), 0, router.ModbusAddr)
	} else if value.Type == base.DATATYPE_DB_UDWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE_EX, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, uint32(value.ValueInt), 0, router.ModbusAddr)
	}

	return insert_sql
}

func FmtInsertAlterSQL(table_name string, router *base.RouterAlter) string {
	insert_sql := fmt.Sprintf(SQL_INSERT_ALTER_TABLE, table_name, router.RouterID, router.MachineID, router.AlterIDDB, router.DataTypeDB, base.GetString(router.Data), router.ModbusAddr, router.Status)

	return insert_sql
}
