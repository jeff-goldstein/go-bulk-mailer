package csv

import (
	"github.com/therahulprasad/go-bulk-mailer/config"
	"os"
	"github.com/therahulprasad/go-bulk-mailer/common"
	"encoding/csv"
	"io"
	"strings"
	"sync"
	"github.com/therahulprasad/go-bulk-mailer/logger"
	"fmt"
	"flag"
)

var overriddenCsvSrc string = ""

func Process(tpl_h, tpl_t string, ch chan common.Mail, wg *sync.WaitGroup) {
	defer wg.Done()
	// Load config
	config := config.GetConfig()

	// Load file from config
	csvSrc := ""
	if overriddenCsvSrc == "" {
		csvSrc = config.Source.Csv.Src
	} else {
		csvSrc = overriddenCsvSrc
	}

	fp, err := os.Open(csvSrc)
	common.FailOnErr(err, "Error while reading source file")

	// Process file as csv
	cr := csv.NewReader(fp)

	// Read csv file line by line
	// TODO: Make sure that it reads all the line and does not breaks in between
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Log error and skip
			logger.LogError("Source > CSV > Could not read record because " + err.Error())
			continue
		}

		subject := config.Campaign.Subject
		email, name, from_email, from_name := "", "", "", ""

		// If variable exists in csv config
		// Replace variable name in template by a field from csv
		for i,v := range config.Source.Csv.Variables {
			// // Make sure index exists in record and replace variable in template file
			// If number of columns is more than current variable index
			// Otherwise if number of variables are more than columns in csv then application will panic
			if len(record)  > i {
				// Construct variable name which will be something like {$name} or {$email} in template file
				vname := "{$" + v + "}"

				//  in variable array should correspond to csv record and hence it should be available
				val := record[i]

				// Store value of email variable separately to construct message
				if v == "email" {
					email = val
				} else if v == "name" {
					name = val
				} else if v == "from_email" {
					from_email = val
				} else if v == "from_name" {
					from_name = val
				} else if v == "subject" {
					subject = val
				}

				tpl_h = strings.Replace(tpl_h, vname , val, -1)
				if tpl_h != "" {
					tpl_t = strings.Replace(tpl_t, vname, val, -1)
				}

				// TODO: Check if variable exists in subject to avoid repeating string replace for static subject
				subject = strings.Replace(subject, vname, val, -1)
			}
		}

		m := common.Mail{}
		if email == "" {
			// Log error as email not found and skip
			logger.LogError("CSV > Record > Email not found > " + fmt.Sprintf("%s", record))
			continue
		}

		// Add recipient to message
		m.AddRecipient(email, name, "to")

		// Check if from_name is overriden in csv file otherwise use config's from_name
		if from_name == "" {
			m.FromName = config.Campaign.FromName
		} else {
			m.FromName = from_name
		}

		// Check if from_email is overriden in csv file otherwise use config's from_name
		if from_email == "" {
			m.FromEmail = config.Campaign.FromEmail
		} else {
			m.FromEmail = from_email
		}

		m.Subject = subject
		m.HTML = tpl_h
		m.Text = tpl_t

		ch <- m
	}
	close(ch)
}

var csvSourceFlag *string
func InitFlags() {
	csvSourceFlag = flag.String("csv.source", "no-value", "Override csv source file")
}

func ProcessFlags() {
	// If csv source is provided
	if csvSourceFlag != nil && *csvSourceFlag != "no-value" {
		config.OverrideCsvSrc(*csvSourceFlag)
	}
}
