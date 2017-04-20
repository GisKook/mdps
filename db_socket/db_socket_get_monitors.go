package db_socket

import (
	"log"
	"strconv"
)

const sql_get_monitors string = "SELECT ID,DATATYPE,ADDRESS FROM DMS_DMP WHERE MID=?"

func (db *DbSocket) GetMonitors(router_id uint32) uint64 {
	log.Println("--------select id---------")
	log.Println(router_id)
	log.Println(sql_get_plc_id)
	log.Println("--------select id---------")
	stmt, err := db.Db.Prepare(sql_get_monitors)
	defer stmt.Close()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	return plc_id
}
