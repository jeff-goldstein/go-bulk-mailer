package dispatcher

import (
	"github.com/therahulprasad/go-bulk-mailer/config"
	"github.com/therahulprasad/go-bulk-mailer/common"
	"sync"
	"github.com/therahulprasad/go-bulk-mailer/dispatcher/sparkpost"
)

func Init(conf config.Configuration, ch chan common.Mail, wg *sync.WaitGroup) {
	// If Mandrill dispatcher is used
	if conf.Dispatcher.Use == "sparkpost" {
		// Start specified number of workers
		for i:=0; i<conf.Dispatcher.WorkersCount; i++ {
			// Add number of workers to waitgroup, so it can wait for same number of workers to exit
			wg.Add(1)

			// Start workers
			go sparkpost.Worker(conf, ch, wg)
		}
	}
}