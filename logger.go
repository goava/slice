package slice

import (
	"log"
	"os"
)

// this variable need for replace std logger to mock
var stdLog Logger = stdLogger{}

// Logger
type Logger interface {
	Info(format string, values ...interface{})
	Error(err error)
	Fatal(err error)
}

type stdLogger struct {
}

// Info logs message with info level.
func (s stdLogger) Info(format string, values ...interface{}) {
	log.Printf(format, values...)
}

func (s stdLogger) Fatal(err error) {
	log.Print("Fatal: ", err.Error())
	os.Exit(1)
}

func (s stdLogger) Error(err error) {
	log.Print("Error: ", err.Error())
}
