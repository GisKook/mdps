package zmq_server

import (
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/pb"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

type Http_server struct {
	Addr             string
	HttpRespones     map[uint64]chan *Report.ControlCommand
	HttpRequestAdd   chan *HttpRequestPair
	HttpRequestDel   chan uint64
	HttpResponseChan chan *HttpRequestPair

	SerialID    uint32
	SerialIDMap map[uint16]uint16
}

var G_MutexHttpServer sync.Mutex
var G_HttpServer *Http_server

func NewHttpServer(config *conf.HttpConf) *Http_server {
	defer G_MutexHttpServer.Unlock()
	G_MutexHttpServer.Lock()
	G_HttpServer = &Http_server{
		Addr:             config.Addr,
		HttpRespones:     make(map[uint64]chan *Report.ControlCommand),
		HttpRequestAdd:   make(chan chan *Report.ControlCommand),
		HttpRequestDel:   make(chan uint64),
		HttpResponseChan: make(chan *HttpResponsePair),
		SerialID:         0,
		SerialIDMap:      make(map[uint16]uint16),
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

	http.HandleFunc(HTTP_TRANSPARENT_TRANSMISSION, TransparentTransmissionHandler)
	http.HandleFunc(HTTP_RELEASE_TRANSPARENT_TRANSMISSION, ReleaseTransparentTransmissionHandler)
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

func (h *Http_server) IncreaseSerial() uint16 {
	return uint16(atomic.AddUint32(&h.SerialID, 1))
}

func (h *Http_server) SetSerialID(origin_serial uint32) uint16 {
	inner_serial := h.IncreaseSerial()
	h.SerialIDMap[inner_serial] = origin_serial

	return inner_serial
}
