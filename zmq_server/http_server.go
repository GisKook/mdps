package zmq_server

import (
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"net/http"
	"sync"
)

type Http_server struct {
	Addr         string
	HttpRespones map[uint64]chan *Report.ControlCommand
	SerialID     uint32
}

var G_MutexHttpServer sync.Mutex
var G_HttpServer *Http_server

func NewHttpServer(config *conf.HttpConf) *Http_server {
	defer G_MutexHttpServer.Unlock()
	G_MutexHttpServer.Lock()
	G_HttpServer = &Http_server{
		Addr:         config.Addr,
		HttpRespones: make(map[uint64]chan *Report.ControlCommand),
		SerialID:     0,
	}

	return G_HttpServer
}

func (server *Http_server) Init() {
	http.HandleFunc(HTTP_RESTART, RestartHandler)
	http.HandleFunc(HTTP_DATA_DOWNLOAD, DataDownloadHandler)
	http.HandleFunc(HTTP_DATA_QUERY, DataQueryHandler)
	http.HandleFunc(HTTP_BATCH_ADD_MONITOR, BatchAddMonitorHandler)
	http.HandleFunc(HTTP_BATCH_ADD_ALTER, BatchAddAlterHandler)

	http.HandleFunc(HTTP_SET_SERVER_ADDR, SetServerAddrHandler)
	http.HandleFunc(HTTP_GET_SERVER_ADDR, GetServerAddrHandler)
	http.HandleFunc(HTTP_RS232_GET_CONFIG, Rs232GetConfigHandler)
	http.HandleFunc(HTTP_RS232_SET_CONFIG, Rs232SetConfigHandler)
	http.HandleFunc(HTTP_RS485_GET_CONFIG, Rs485GetConfigHandler)
	http.HandleFunc(HTTP_RS485_SET_CONFIG, Rs485SetConfigHandler)
}

func (server *Http_server) Start() {
	err := http.ListenAndServe(server.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe :", err)
	}
}

func GetHttpServer() *Http_server {
	defer G_MutexHttpServer.Unlock()
	G_MutexHttpServer.Lock()

	return G_HttpServer
}
