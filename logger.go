package slice

import (
	"log"
	"os"
)

var (
	// this variable need for replace std logger to mock
	stdLog Logger = stdLogger{}
)

// Logger
type Logger interface {
	Error(err error)
	Fatal(err error)
}

type stdLogger struct {
}

func (s stdLogger) Fatal(err error) {
	log.Print("Fatal:", err.Error())
	os.Exit(1)
}

func (s stdLogger) Error(err error) {
	log.Print("Error:", err.Error())
}
