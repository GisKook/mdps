package zmq_server

import (
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"net/http"
)

type Http_server struct {
	Addr         string
	HttpRespones map[uint64]chan *Report.ControlCommand
}

var G_HttpServer *Http_server

func NewHttpServer(config *conf.HttpConf) *Http_server {
	G_HttpServer = &Http_server{
		Addr:         config.Addr,
		HttpRespones: make(map[uint64]chan *Report.ControlCommand),
	}

	return G_HttpServer
}

func (server *Http_server) Init() {
	http.HandleFunc(HTTP_RESTART, RestartHandler)
}

func (server *Http_server) Start() {
	err := http.ListenAndServe(server.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe :", err)
	}
}

func GetHttpServer() *Http_server {
	return G_HttpServer
}
