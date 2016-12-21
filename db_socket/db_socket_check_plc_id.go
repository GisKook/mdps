package db_socket

import (
	"log"
)

const sql_check_plc_id string = "SELECT EXISTS(SELECT 1 FROM DMS_ROUTER WHERE ID = ?)"

func (db *DbSocket) CheckPlcID(id uint64) uint8 {

	stmt, err := db.Db.Prepare(sql_check_plc_id)
	defer stmt.Close()
	if err != nil {
		log.Println(err.Error())
		return 1
	}

	var check uint8

	err = stmt.QueryRow(id).Scan(&check)
	if err != nil {
		log.Println(err.Error())
		return 1
	}

	return check
}
