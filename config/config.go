package config

import (
	"encoding/json"
	"github.com/therahulprasad/go-bulk-mailer/common"
	"os"
	"fmt"
	"strconv"
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
			WaitForMilliSeconds  int `json:"wait_for_x_milliseconds"`
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
var overrideCsvSrc string = ""

// Displays loaded config details in terminal
func PrintDetails() {
	print("Campaign Title", config.Campaign.Title)
	print("Campaign Subject", config.Campaign.Subject)
	if config.Source.Use == "csv" {
		print("Data source - CSV", config.Source.Csv.Src)
	}
	print("Template Source", config.Template.HTMLSrc)
	if config.Dispatcher.Use == "sparkpost" {
		print("Dispatcher", "sparkpost")
	}
	print("Log Folder", config.Logger.FolderPath)
	print("Log Success", strconv.FormatBool(config.Logger.LogSuccess))
}

// For formatted printing
func print(msg, val string) {
	fmt.Println(msg + ": " + val)
}

func ConfirmEarlyUsage() bool {
	// Check if success log filename exists
	success_log_filename := GetSuccessLogFileName()
	_, err := os.Stat(success_log_filename)
	if err == nil {
		// Success file already exists
		return true
	} else {
		return false
	}
}
func GetSuccessLogFileName() string {
	campaign_title := common.LowerAlphaNumericFilter(config.Campaign.Title)
	src_file := GetSourcePathIdentifier()
	return config.Logger.FolderPath + "/" + campaign_title + "-" + src_file + "success.log"
}
func GetSourcePathIdentifier() string {
	if config.Source.Use == "csv" {
		return common.LowerAlphaNumericFilter(config.Source.Csv.Src)
	}
	return "";
}
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
		// Check if CSV is overridden
		if overrideCsvSrc != "" {
			config.Source.Csv.Src = overrideCsvSrc
		}

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

		// Source path must exists
		if !common.FileExists(config.Source.Csv.Src) {
			common.Fail("CSV source file does not exists or is not accessible: " + config.Source.Csv.Src)
		}
	}

	if config.Dispatcher.WorkersCount <= 0 {
		common.Fail("Dispater's worker count should be +ve number in config file")
	}
	if config.Dispatcher.Worker.WaitForMilliSeconds < 0 {
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

func GetIfLogSuccess() bool {
	return config.Logger.LogSuccess
}

func OverrideCsvSrc(csvSrc string) {
	if common.FileExists(csvSrc) {
		overrideCsvSrc = csvSrc
	} else {
		common.Fail("CSV source file does not exists or is not accessible.")
	}
}