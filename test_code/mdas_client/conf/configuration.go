package conf

import (
	"encoding/json"
	"os"
)

type ServerConfiguration struct {
	Addr       string
	ReadLimit  uint16
	WriteLimit uint16
}

type ClientConfiguration struct {
	HeartInterval uint8
}

type Configuration struct {
	Server *ServerConfiguration
	Client *ClientConfiguration
}

var G_conf *Configuration

func ReadConfig(confpath string) (*Configuration, error) {
	file, _ := os.Open(confpath)
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	G_conf = &config

	return &config, err
}

func GetConf() *Configuration {
	return G_conf
}
