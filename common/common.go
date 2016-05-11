package common

import (
	"log"
	"regexp"
	"strings"
	"os"
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

func LowerAlphaNumericFilter(src string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]")
	safe := reg.ReplaceAllString(src, "")
	return strings.ToLower(safe)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}