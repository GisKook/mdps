package conf

import (
	"encoding/json"
	"os"
)

type HttpConf struct {
	Addr    string
	Timeout uint8
}

type ZmqConf struct {
	TerminalManageUp    string
	TerminalManageDown  string
	TerminalControlUp   string
	TerminalControlDown string
	TerminalDataUp      string
}

type RedisConf struct {
	Addr              string
	MaxIdle           int
	IdleTimeOut       int
	Passwd            string
	TranInterval      int
	MonitorsKeyExpire uint32
	AltersKeyExpire   uint32
	StatusExpire      uint32
}

type Configuration struct {
	Http  *HttpConf
	Zmq   *ZmqConf
	Redis *RedisConf
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
