package main

import (
	"github.com/therahulprasad/go-bulk-mailer/config"
	"os"
	"log"
	"github.com/therahulprasad/go-bulk-mailer/source"
	"github.com/therahulprasad/go-bulk-mailer/common"
	"github.com/therahulprasad/go-bulk-mailer/dispatcher"
	"sync"
	"github.com/therahulprasad/go-bulk-mailer/logger"
)

const version  = "1.0"

/**
	TODO:
	[ ] Success.log and Error.log are always empty
 */

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./go-bulk-mailer config.json")
	}

	cfp := os.Args[1]

	// Load configurations
	config := config.LoadConfig(cfp)

	// Initiate Loggers
	logger.Init(config)

	// Channel to send emails
	chMail := make(chan common.Mail)

	// Declare a waitgroup to wait for all workers
	var wg sync.WaitGroup

	// Read from source
	// Start worker to process log file and pass it via channel
	// This function returns immediately
	// ITs source responsibility to close channel when work is finished
	source.Init(config, chMail, &wg)

	// Dispatch emails using one of the dispatcher
	// Start worker which receives mail messages and sends it
	// If chMail closes, all dispatcher's workers should finish
	dispatcher.Init(config, chMail, &wg)

	// Wait for all go routines (source and dispatcher's workers) to end
	wg.Wait()

	// Close open files by logger
	logger.Destroy()
}