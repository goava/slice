package slice

import (
	"log"
)

// this variable need for replace std logger to mock
var stdLog Logger = stdLogger{}

// Logger
type Logger interface {
	Printf(format string, values ...interface{})
	Fatal(err error)
}

type stdLogger struct {
}

// Info logs message with info level.
func (s stdLogger) Printf(format string, values ...interface{}) {
	log.Printf(format, values...)
}

func (s stdLogger) Fatal(err error) {
	log.Fatal(err.Error())
}
