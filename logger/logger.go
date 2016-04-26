package logger

import (
	"os"
	"github.com/therahulprasad/go-bulk-mailer/config"
	"github.com/therahulprasad/go-bulk-mailer/common"
	"log"
)

// File Handler for success and error logs
var fe, fs *os.File
var slogger, elogger *log.Logger

func Init(conf config.Configuration) {
	// Check if Log folder exists, if not create it
	_, err := os.Stat(conf.Logger.FolderPath)

	if os.IsNotExist(err) {
		err = os.Mkdir(conf.Logger.FolderPath, 0775)
		common.FailOnErr(err, "Could not create log directory")
	} else {
		common.FailOnErr(err, "Error while opening log directory")
	}

	// Open success and error file and handle error
	fe, err = os.OpenFile(conf.Logger.FolderPath + "/" + "error.log", os.O_WRONLY|os.O_APPEND, 0666)
	common.FailOnErr(err, "Could not create/open error.log file")

	fs, err = os.OpenFile(conf.Logger.FolderPath + "/" + "success.log", os.O_WRONLY|os.O_APPEND, 0666)
	common.FailOnErr(err, "Could not create/open success.log file")


	elogger = log.New(fe, "", log.LstdFlags | log.Lshortfile)
	slogger = log.New(fs, "", 0)
}

func Destroy() {
	fe.Close()
	fs.Close()
}

func LogError(msg string) {
	elogger.Println(msg)
}
func LogOnError(err error, msg string) {
	if err != nil {
		LogError(msg + " | " + err.Error())
	}
}

func LogSuccess(msg string) {
	slogger.Println(msg)
}
func LogOnSuccess(err error, msg string) {
	if err == nil {
		LogSuccess(msg)
	}
}