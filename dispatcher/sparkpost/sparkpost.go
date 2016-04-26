package sparkpost

import (
	"github.com/therahulprasad/go-bulk-mailer/config"
	"github.com/therahulprasad/go-bulk-mailer/common"
	sp "github.com/SparkPost/gosparkpost"
	"sync"
	"github.com/therahulprasad/go-bulk-mailer/logger"
	"fmt"
)

func Worker(conf config.Configuration, ch chan common.Mail, wg *sync.WaitGroup) {
	defer wg.Done()

	// Initiate a Mandrill client for each worker
	cfg := &sp.Config{
		BaseUrl:    conf.Dispatcher.Sparkpost.BaseUrl,
		ApiKey:     conf.Dispatcher.Sparkpost.API_KEY,
		ApiVersion: conf.Dispatcher.Sparkpost.ApiVersion,
	}
	var client sp.Client
	err := client.Init(cfg)
	common.FailOnErr(err, "Sparkpost client could not be initiated")

	// Start receiving messages from channel
	for msg := range ch {

		// There must be one recipient
		if len(msg.Recipients) == 0 {
			// Log error and skip
			logger.LogError("Mandrill > msg has no recipient" + fmt.Sprintf("%s", msg))
			continue
		}

		// If Body, (HTML and TEXT) is not supplied log error and skip
		if msg.HTML == "" && msg.Text == "" {
			// If both type of body are blank then log error and skip
			logger.LogError("Mandrill > msg has no body > " + fmt.Sprintf("%s", msg))
			continue
		}

		tx := &sp.Transmission{}

		// Add recipients to transmission
		numEmails := 0
		emails := []string{}
		recipients := []sp.Recipient{}
		for _,r := range msg.Recipients {
			re := make(map[string]string)
			if r.Email == "" {
				// Email can not be empty
				// Log it and skip
				logger.LogError(fmt.Sprintf("SparkPost > Email can not be empty %s", msg))
				continue
			}
			numEmails++
			if r.Name != "" {
				re["name"] = r.Name
			}
			re["email"] = r.Email
			reci := sp.Recipient{Address: re}
			recipients = append(recipients, reci)

			// Keeping a list of email for logging on success
			emails = append(emails, r.Email)
		}

		// If no emails found, log error and skip
		if numEmails == 0 {
			logger.LogError(fmt.Sprintf("SparkPost > No recipient contain email address %s", msg))
		}

		tx.Recipients = recipients

		content := sp.Content{}
		content.HTML = msg.HTML

		// If Text email is present, then add it to content object
		if msg.Text != "" {
			content.Text = msg.Text
		}

		// If msg does not contain necessary field take it from config
		from := make(map[string]string)
		if msg.FromEmail == "" {
			from["email"] = conf.Campaign.FromEmail
		} else {
			from["email"] = msg.FromEmail
		}
		if msg.FromName == "" {
			from["name"] = conf.Campaign.FromName
		} else {
			from["name"] = msg.FromName
		}
		content.From = from

		// If msg does not contain necessary field take it from config
		if msg.Subject == "" {
			content.Subject = conf.Campaign.Subject
		} else {
			content.Subject = msg.Subject
		}

		tx.Content = content

		id, _, err := client.Send(tx)

		// Remove msg body before loggin
		msg.HTML = ""
		msg.Text = ""
		logger.LogOnError(err, "Could not send email: " + fmt.Sprintf("%s", msg))
		logger.LogOnSuccess(err, fmt.Sprintf("%s - %s", id, emails))
	}
}
