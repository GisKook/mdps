package main

import (
	"fmt"
	//zmq "github.com/pebbe/zmq3"
	"github.com/giskook/mdps/conf"
	"github.com/giskook/mdps/zmq_server"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config, _ := conf.ReadConfig("./conf.json")
	log.Println(config)

	zmq_server := zmq_server.NewZmqServer()
	zmq_server.Init(config.Zmq)
	zmq_server.Run()
	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}
