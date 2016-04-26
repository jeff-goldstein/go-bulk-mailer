package config

import (
	"encoding/json"
	"github.com/therahulprasad/go-bulk-mailer/common"
	"os"
)

type Configuration struct {
	Campaign struct {
		Title     string `json:"title"`
		Subject   string `json:"subject"`
		FromName  string `json:"from_name"`
		FromEmail string `json:"from_email"`
	} `json:"campaign"`
	Source struct {
		Use string `json:"use"`
		Csv struct {
			Src       string   `json:"src"`
			Variables []string `json:"variables"`
		} `json:"csv"`
	} `json:"source"`
	Template struct {
		HTMLSrc string `json:"html_src"`
		TextSrc string `json:"text_src"`
	} `json:"template"`
	Dispatcher struct {
		Use       string `json:"use"`
		Sparkpost struct {
			API_KEY    string `json:"API_KEY"`
			BaseUrl    string `json:"base_url"`
			ApiVersion int    `json:"api_version"`
		} `json:"sparkpost"`
		WorkersCount int `json:"workers_count"`
		Worker       struct {
			WaitForMicroSeconds  int `json:"wait_for_x_microseconds"`
			AfterDispatchingMsgs int `json:"after_dispatching_y_msgs"`
		} `json:"worker"`
	} `json:"dispatcher"`
	Logger struct {
		FolderPath string `json:"folder_path"`
		LogSuccess bool   `json:"log_success"`
	} `json:"logger"`
}

var config Configuration
var configLoaded bool = false

func LoadConfig(path string) Configuration {
	fc, err := os.Open(path)
	common.FailOnErr(err, "Could not open config file")
	defer fc.Close()

	cd := json.NewDecoder(fc)
	err = cd.Decode(&config)
	common.FailOnErr(err, "Could not decode configuration json")

	// Make sure that mandatory fields exists
	// 1. config.Template.HTMLSrc
	fp, err := os.Open(config.Template.HTMLSrc)
	common.FailOnErr(err, "Could not open HTML Template file")
	fp.Close()

	// If Text template file path is specified then it must exist
	if config.Template.TextSrc != "" {
		fp, err = os.Open(config.Template.TextSrc)
		common.FailOnErr(err, "Could not open Text Template file")
		fp.Close()
	}

	// Check if Campaign title is present
	if config.Campaign.Title == "" {
		common.Fail("Campaign title should not be empty")
	}

	// TODO: Also validate that its a proper email address
	if config.Campaign.FromEmail == "" {
		common.Fail("Campaign must have a from_email")
	}

	if config.Campaign.FromName == "" {
		common.Fail("Campiagn must have a from_name")
	}

	// Check if Campaign subject is provided
	if config.Campaign.Subject == "" {
		common.Fail("Campaign subject should not be empty")
	}

	// If configuration has CSV source
	if config.Source.Use == "csv" {
		// One of the variable must be email
		if len(config.Source.Csv.Variables) > 0 {
			emailFound := false
			for _, v := range config.Source.Csv.Variables {
				if v == "email" {
					emailFound = true
				}
			}
			if emailFound == false {
				common.Fail("There must be one variable called email")
			}
		} else {
			common.Fail("There must be one variable called email")
		}
	}

	if config.Dispatcher.WorkersCount <= 0 {
		common.Fail("Dispater's worker count should be +ve number in config file")
	}
	if config.Dispatcher.Worker.WaitForMicroSeconds < 0 {
		common.Fail("Dispatcher's worker's wait_for_x_microseconds should be 0 or more")
	}
	if config.Dispatcher.Worker.AfterDispatchingMsgs < 0 {
		common.Fail("Dispatcher's worker's after_dispatching_y_msgs should be 0 or more")
	}
	if config.Dispatcher.Use == "sparkpost" {
		if config.Dispatcher.Sparkpost.API_KEY == "" {
			common.Fail("Mandrill's API key not found in config file")
		}
	}

	configLoaded = true

	return config
}

func IsConfigLoaded() bool {
	return configLoaded
}

func GetConfig() Configuration {
	return config
}
