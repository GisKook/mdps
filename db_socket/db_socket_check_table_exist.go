package db_socket

import (
	"fmt"
	"log"
)

const (
	CHECK_TABLE_FMT string = "SELECT 1 FROM %s LIMIT 1;"
)

func (socket *DbSocket) check_table_if_exist(table string) bool {
	sql_check_table := fmt.Sprintf(CHECK_TABLE_FMT, table)
	_, err := socket.Db.Exec(sql_check_table)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
