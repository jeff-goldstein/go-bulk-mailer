package common

import (
	"log"
)

func FailOnErr(err error, msg string) {
	log.SetFlags(log.Ldate | log.Ltime)
	if err != nil {
		log.Fatal(msg + " | " + err.Error())
	}
}

func Fail(msg string) {
	log.Fatal(msg)
}