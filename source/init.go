package source

import (
	"github.com/therahulprasad/go-bulk-mailer/config"
	c "github.com/therahulprasad/go-bulk-mailer/source/csv"
	"io/ioutil"
	"github.com/therahulprasad/go-bulk-mailer/common"
	"sync"
)

func Init(config config.Configuration, ch chan common.Mail, wg *sync.WaitGroup) {
	// Load HTML template file
	tmp, err := ioutil.ReadFile(config.Template.HTMLSrc)
	common.FailOnErr(err, "Error while reading html template source")
	tpl_html := string(tmp)

	// Load Text template file
	tpl_text := ""
	tmp, err = ioutil.ReadFile(config.Template.HTMLSrc)
	if err == nil {
		tpl_text = string(tmp)
	}

	// Check which source to use from config file
	if config.Source.Use == "csv" {
		wg.Add(1)
		// Process CSV file as goroutine and close channel on exit
		go c.Process(tpl_html, tpl_text, ch, wg)
	}
}

func InitFlags() {
	c.InitFlags()
}

func ProcessFlags() {
	c.ProcessFlags()
}