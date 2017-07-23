# go-bulk-mailer
__v0.1 beta__
Its an application written in go to send bulk email.
>"Somebody please contribute unit test cases, so that I can remove beta tag."  
>- Rahul Prasad
  
You can choose source from following options
1. csv (working on it)
1. json (todo)

And mail user following services
1. Sparkpost (working on it)  
1. sendgrid (todo) 
1. SMTP (todo)

## Features
1. High throughput because of Go's concurrency model
1. Throttle number of emails send in specific amount of time
1. Modular structure. Its easy to add your service provider such as sendgrid.
1. Interactive CLI mode which can be overridden to be used dynamically
1. Support for overriding csv.source using flag

## Installation
1. 

## Building from source
1. You must have go environment setup first. Learn more at https://golang.org/
1. Run `go get github.com/therahulprasad/go-bulk-mailer`
1. Goto `$GOPATH/src/github.com/therahulprasad/go-bulk-mailer` and run go build
1. It should create an executable of name `go-bulk-mailer`

## Usage
For help use `./go-bulk-mailer -h`  
To send mail, create `config.json` file and run `./go-bulk-mailer --config=config.json`

#### Config
    {
      "campaign": {
        "###title": "Title of the campaign",
        "title": "reinstall mailer",
    
        "###subject": "Subject which will be sent as mail. You can use variable like {$name} within subject if name column comes after subject column in csv.",
        "subject": "{$name} Please come back",
    
        "from_email": "noreply@bobbleapp.me",
        "from_name": "Bobble App"
      },
      "source": {
        "###use": "which source will be used for reading recipient information. Currently only csv is supported",
        "use": "csv",
    
        "csv": {
          "###src": "Path of csv file which contains recipient information",
          "src": "dump.csv",
    
          "###variables": "Array of strings, It must be in same order as csv file. 'email' is a mandatory variable, 'name', 'from_email', 'from_name', 'subject' (If subject if first column, other vaiables may be used within its value) are another static optional variable.",
          "variables": ["subject", "email", "gender", "from_name"]
        }
      },
      "###template": "variables can be used inside template files in this format {$name}. It will be replaced by data found in source",
      "template": {
        "html_src": "template.html",
        "text_src": "template.txt"
      },
      "dispatcher": {
        "###use": "which dispatcher service to use? Currently only sparkpost is available",
        "use": "sparkpost",
    
        "sparkpost": {
          "API_KEY": "xxx",
          "base_url": "https://api.sparkpost.com",
          "api_version": 1
        },
        "###workers_count": "Workers will be working concurrently, it can be increased for better throughput.",
        "workers_count": 1,
    
        "worker": {
          "###wait_for_x_milliseconds": "Number seconds each worker will wait after dispatching certain message. If set to 0 messages will be dispatched continuously. This might cause heavy load on dispatcher.",
          "wait_for_x_milliseconds": 1,
    
          "###after_dispatching_y_msgs": "Number of messages which has to be dispatched by single worker before sleeping for certain seconds. If set to 0 messages will be dispatched continuously.",
          "after_dispatching_y_msgs": 10
        }
      },
      "logger": {
        "###folder_path": "Path of the folder where logs will be stored. If folder does not exists, it will be created.",
        "folder_path": "logs",
    
        "###log_success": "true/false To log successful requests",
        "log_success": true
      }
    }

## Development
> Pull request are most welcome.  
 
Currently I have implemented support for csv as source and sparkpost as service. 
You are most welcome to add more services and sources.

#### Todo
1. Create test mode flag to send test email
1. Implement support for sending attachment
1. Daemon mode
1. Generate sample config file
1. create go get cli just like https://github.com/sideshow/apns2
1. Make config modular, each module should be responsible for loading its config 

#### Folder structure
`main.go` contains bootstapping code  
`init` folder contains business logic  
`config` folder contains configuration management  
`source` contains logic to process email data sources  
`dispatcher` contains service providers  
`common` folder contains Models, Logger, Error Handler  

#### Creating new dispatcher say _sendgrid_
1. Create folder under dispatcher called _sendgrid_
1. Update dispatcher/init.go and code to initiate sendgrid's worker
1. Pass channel for receiving msgs to sendgrid's worker 
1. Worker can range over channel and send emails 
1. Source will close channel if all the messages are processed.
1. Detect if channel is closed and shut down workers accordingly


## Changelog
__v0.1 beta__
This is first version


