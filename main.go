package main

import (
	"github.com/therahulprasad/go-bulk-mailer/config"
	"os"
	"github.com/therahulprasad/go-bulk-mailer/source"
	"github.com/therahulprasad/go-bulk-mailer/common"
	"github.com/therahulprasad/go-bulk-mailer/dispatcher"
	"sync"
	"github.com/therahulprasad/go-bulk-mailer/logger"
	"flag"
	"fmt"
	"bufio"
	"log"
)

const version  = "0.1.0"
const appName = "go-bulk-mailer"

// TODO: Move flag logic to a separate function or a file
func main() {
	// Loaf fag details from command line arguments
	versionFlag := flag.Bool("version", false, "Displays version and exist")
	noDetailsFlag := flag.String("no-details", "no-value", "Do not print details for review.")
	configFlag := flag.String("config", "no-value", "Path of config file")
	noWarnFlag := flag.String("no-warning", "no-value", "Do warn if same source is already used for same campaign")
	source.InitFlags()
	flag.Parse()

	// If --version flag is set then just display version and exit
	if(*versionFlag == true) {
		// Display application name
		fmt.Println(appName)
		// Print version
		fmt.Println("Version: " + version)
		os.Exit(0)
	}
	source.ProcessFlags()

	if (*configFlag == "no-value") {
		log.Fatal("Usage: ./go-bulk-mailer --config=/path/to/config.json")
	}

	// Load configurations
	conf := config.LoadConfig(*configFlag)

	// Show important details of config before beginning if user has not overwritten it
	if (*noDetailsFlag != "true") {
		// Print config details
		config.PrintDetails()

		// Ask user if they want to continue
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Looks Good (yes/no): ")
		text, _ := reader.ReadString('\n')

		// If user does not enter yes to confirm that details are correct. Exit.
		if text != "yes\n" {
			fmt.Println("No! Exiting")
			os.Exit(1)
		} else {
			fmt.Println("Yes")
		}
	}


	if *noWarnFlag != "true" && config.ConfirmEarlyUsage() == true {
		fmt.Printf("Campaign \"%s\" has already been run for source \"%s\" \nAre you sure you want to proceed ? (yes/no) ", conf.Campaign.Title, config.GetSourcePathIdentifier())

		/// Ask user for confirmation
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		// If user does not enter yes to confirm that details are correct. Exit.
		if text != "yes\n" {
			fmt.Println("No! Exiting")
			os.Exit(1)
		} else {
			fmt.Println("Yes")
		}
	}

	// Initiate Loggers
	logger.Init(conf)

	// Channel to send emails
	chMail := make(chan common.Mail)

	// Declare a waitgroup to wait for all workers
	var wg sync.WaitGroup

	// Read from source
	// Start worker to process log file and pass it via channel
	// This function returns immediately
	// ITs source responsibility to close channel when work is finished
	source.Init(conf, chMail, &wg)

	// Dispatch emails using one of the dispatcher
	// Start worker which receives mail messages and sends it
	// If chMail closes, all dispatcher's workers should finish
	dispatcher.Init(conf, chMail, &wg)

	// Wait for all go routines (source and dispatcher's workers) to end
	wg.Wait()

	// Close open files by logger
	logger.Destroy()
}