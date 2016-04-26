# csv-mailer
Its an application written in go to send bulk email.
  
You can choose source from following options
1. csv (working on it)
1. json (todo)

And mail user following services
1. Mandrill (working on it)
1. sendgrid (todo) 
1. SMTP (todo)

## Features

## Installation

## Usage

#### Config
1. `source`
    1. `csv`
        1. `src`: `string` location of csv file
        1. `variables`: `["name", "gender", "email"]` array of strings  
1. `template`
    1. `src`: `string` location of template file

## Development
Pull request are always welcome. Application is made modular to make it easier 
to add new sources and services. Currently I need support for csv as source 
and mandrill as service. You are most welcome to add more services and sources.

#### Todo
1. Implement support for sending attachment

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
