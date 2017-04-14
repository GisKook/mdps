package main

import (
	"fmt"
	//zmq "github.com/pebbe/zmq3"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/db_socket"
	"github.com/giskook/mdps/redis_socket"
	"github.com/giskook/mdps/zmq_server"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func check_auth() bool {
	if time.Now().Unix() > 1496246399 {
		return false
	}
	return true
}

func main() {
	if !check_auth() {
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	config, _ := conf.ReadConfig("./conf.json")
	log.Println(config)

	zeromq_server := zmq_server.NewZmqServer()
	zeromq_server.Init(config.Zmq)
	go zeromq_server.Run()

	redis_server, _ := redis_socket.NewRedisSocket(config.Redis)
	redis_server.LoadAll()
	go redis_server.DoWork()

	db_socket.NewDbSocket(config.DB)

	http_server := zmq_server.NewHttpServer(config.Http)
	http_server.Init()
	http_server.Start()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}
