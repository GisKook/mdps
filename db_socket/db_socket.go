package db_socket

import (
	"database/sql"
	"fmt"
	"github.com/giskook/mdps/conf"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

type DbSocket struct {
	Db *sql.DB
}

var G_MutuxDBSocket sync.Mutex

var G_DBSocket *DbSocket

func NewDbSocket(db_config *conf.DBConf) (*DbSocket, error) {
	defer G_MutuxDBSocket.Unlock()
	G_MutuxDBSocket.Lock()
	//user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
	log.Println(db_config)
	conn_string := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?autocommit=true",
		db_config.User, db_config.Passwd, db_config.Host, db_config.Port, db_config.DbName)

	log.Println(conn_string)

	db, err := sql.Open("mysql", conn_string)
	db.SetMaxOpenConns(100)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("db open success")

	G_DBSocket = &DbSocket{
		Db: db,
	}

	return G_DBSocket, nil
}

func GetDBSocket() *DbSocket {
	defer G_MutuxDBSocket.Unlock()
	G_MutuxDBSocket.Lock()
	return G_DBSocket
}
