package db_socket

import (
	"strconv"
)

const sql_get_plc_id string = "SELECT TOKEN FROM DMS_ROUTER WHERE CPUID=?"

func (db *DbSocket) GetPlcID(cpuid string) uint64 {
	stmt, err := db.Db.Prepare(sql_get_plc_id)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	var plc_id_string string

	err = stmt.QueryRow(cpuid).Scan(&plc_id_string)
	if err != nil {
		panic(err.Error())
	}

	plc_id, _ := strconv.ParseUint(plc_id_string, 10, 64)

	return plc_id
}
