package db_socket

import (
	"github.com/giskook/mdps/base"
	"log"
)

const sql_get_monitors string = "SELECT M.ROUTERID, P.ID, P.DATATYPE, P.ADDRESS FROM DMS_DMP P,DMS_MACHINE M,DMS_ROUTER R WHERE M.ROUTERID=R.ID AND P.MID=M.ID"

func (db *DbSocket) GetMonitors() []*base.RouterMonitorDB {
	log.Println("--------select id---------")
	log.Println(sql_get_monitors)
	log.Println("--------select id---------")
	stmt, err := db.Db.Prepare(sql_get_monitors)
	defer stmt.Close()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	var router_id uint32
	var monitor_id uint32
	var data_type uint8
	var modbus_addr uint32

	rows, er := stmt.Query()
	base.CheckError(er)
	defer rows.Close()

	router_monitors := make([]*base.RouterMonitorDB, 0)
	for rows.Next() {
		if e := rows.Scan(&router_id, &monitor_id, &data_type, &modbus_addr); e != nil {
			base.CheckError(e)
		}
		router_monitors = append(router_monitors, &base.RouterMonitorDB{
			RouterID:   router_id,
			MonitorID:  monitor_id,
			Datatype:   data_type,
			ModbusAddr: modbus_addr,
		})
	}

	return router_monitors
}
