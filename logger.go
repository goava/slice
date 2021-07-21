package slice

import (
	"log"
	"strings"
)

// this variable need for replace std Logger to mock
var stdLog Logger = stdLogger{}

// Logger
type Logger interface {
	Printf(bundle string, format string, values ...interface{})
	Fatal(err error)
}

type stdLogger struct {
}

// Printf logs message with info level.
func (s stdLogger) Printf(bundle string, format string, values ...interface{}) {
	log.Printf("[%s] "+format, append([]interface{}{strings.ToUpper(bundle)}, values...)...)
}

func (s stdLogger) Fatal(err error) {
	log.Fatal(err.Error())
}
