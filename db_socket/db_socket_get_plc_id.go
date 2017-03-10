package db_socket

import (
	"log"
	"strconv"
)

const sql_get_plc_id string = "SELECT ID FROM DMS_ROUTER WHERE CPUID=?"

func (db *DbSocket) GetPlcID(cpuid string) uint64 {
	log.Println(cpuid)
	stmt, err := db.Db.Prepare(sql_get_plc_id)
	defer stmt.Close()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	var plc_id_string string

	err = stmt.QueryRow(cpuid).Scan(&plc_id_string)
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	plc_id, _ := strconv.ParseUint(plc_id_string, 10, 64)

	return plc_id
}
