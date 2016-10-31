package main

import (
	"fmt"
	//zmq "github.com/pebbe/zmq3"
	"github.com/giskook/mdps/zmq_server"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	zmq_server := zmq_server.NewZmqServer()
	zmq_server.Init()
	zmq_server.Run()
	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)
}
