package base

import (
	"errors"
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

var (
	Error_Redis_Monitor_Expired error = errors.New("Redis Monitor expired")
)
