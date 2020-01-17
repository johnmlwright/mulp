# mulp

MUlti Log Parser is a log parser made to read log files in a more accessible format.
This program was built with the intention of reading JFrog log files (artifactory, catalina, request, etc).

## Getting Started

This program will work out of the box as long as go is on your machine and your GOPATH is set. For additional help running this program, please utilize the `-h` flag.

## Example Uses

Here are some example inputs:
```
$ go run logParser.go
// Uses the default arguments.
$ go run logParser.go -v artifactory.1.log
// Shows verbose output (for debugging) and executes on file artifactory.1.log
$ go run logParser.go -i INFO
// Looks for lines with input "INFO" case sensitive.
```

## Built With

* Golang

## To-Do

* Convert verbose to debug
* * Re-dd verbose as shortcut to unmerge and show time.
* Case Insensitivity flag (for inputs)
* Dynamic timestamp removal 
* * This will deprecate -http
* * Also need to account for text in timestamps (Jan, Feb, etc, not always at start of line, see: catalina.out and localhost.*.out)
* Split program into multiple files
* * src folder?
* Create executables
* Create duration array for http requests.