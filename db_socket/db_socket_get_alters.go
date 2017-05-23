package db_socket

import (
	"fmt"
	"github.com/giskook/mdps/base"
	"log"
)

const sql_get_alters_fmt string = "SELECT P.ID, P.DATATYPE FROM DMS_DAP P,DMS_MACHINE M,DMS_ROUTER R WHERE M.ROUTERID=R.ID AND P.MID=M.ID and M.ROUTERID=%d and P.ADDRESS=%d"

func (db *DbSocket) GetAlters(alters []*base.RouterAlter) {
	for i, v := range alters {
		alter_id, data_type := db.GetAlterDataTypeID(v)
		alters[i].DataTypeDB = data_type
		alters[i].AlterIDDB = alter_id
	}
}

func (db *DbSocket) GetAlterDataTypeID(alter *base.RouterAlter) (uint32, uint8) {
	_sql := db.FmtSelectAlterDataTypeSQL(alter)
	log.Println(_sql)
	stmt, err := db.Db.Prepare(_sql)
	defer stmt.Close()
	if err != nil {
		log.Println(err.Error())
		return 0, 0
	}

	rows, er := stmt.Query()
	base.CheckError(er)
	defer rows.Close()

	var alter_id uint32
	var data_type uint8

	for rows.Next() {
		if e := rows.Scan(&alter_id, &data_type); e != nil {
			base.CheckError(e)
		}
	}

	return alter_id, data_type
}

func (Db *DbSocket) FmtSelectAlterDataTypeSQL(alter *base.RouterAlter) string {
	_sql := fmt.Sprintf(sql_get_alters_fmt, alter.RouterID, alter.ModbusAddr)

	return _sql
}
