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
	SQL_INSERT_ALTER_TABLE     string = "INSERT %s (ALTER_ID, DATATYPE, INTBITS, DECIMALBITS, ADDRESS) VALUES(%d, %d, %d, %d, %d)"
)

func (socket *DbSocket) InsertAlters(router_alter []*base.RouterAlter) {
	tx, err := socket.Db.Begin()
	base.CheckError(err)
	for _, v := range router_alter {
		value, e := socket.ParseAlterValue(v)
		if e == nil {
			insert_sql := FmtInsertAlterSQL(v, value)
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

func FmtInsertAlterSQL(router *base.RouterAlter, value *base.Variant) string {
	time_stamp := time.Now().Unix()
	var insert_sql string
	if value.Type == base.DATATYPE_DB_FLOAT {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, int32(value.Float), int32(value.Float*100)%100, router.ModbusAddr)
	} else if value.Type == base.DATATYPE_DB_SWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, int16(value.ValueInt), 0, router.ModbusAddr)
	} else if value.Type == base.DATATYPE_DB_UWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, uint16(value.ValueInt), 0, router.ModbusAddr)
	} else if value.Type == base.DATATYPE_DB_SDWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, int32(value.ValueInt), 0, router.ModbusAddr)
	} else if value.Type == base.DATATYPE_DB_UDWORD {
		insert_sql = fmt.Sprintf(SQL_INSERT_ALTER_TABLE, GetAlterTableName(time_stamp), router.AlterIDDB, router.DataTypeDB, uint32(value.ValueInt), 0, router.ModbusAddr)
	}

	return insert_sql
}
