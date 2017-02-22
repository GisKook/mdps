package main

import (
	"bufio"
	"fmt"
	"github.com/giskook/mdas_client/conf"
	"github.com/giskook/mdas_client/conn"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// read configuration
	configuration, err := conf.ReadConfig("./conf.json")

	checkError(err)

	file, _ := os.Open("./id.txt")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_id := scanner.Text()
		id, _ := strconv.ParseUint(_id, 10, 64)
		go conn.NewConn(id, configuration).Start()
	}
	//	conn := conn.NewConn(12, configuration)
	//	go conn.Start()
	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
