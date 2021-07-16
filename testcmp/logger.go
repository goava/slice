package testcmp

import (
	"fmt"
	"log"
)

type Logger struct {
	PrintLogs []string
	FatalLogs []string
}

func (l *Logger) Printf(format string, values ...interface{}) {
	s := fmt.Sprintf(format, values...)
	l.PrintLogs = append(l.PrintLogs, s)
	log.Printf(s)
}

func (l *Logger) Fatal(err error) {
	l.FatalLogs = append(l.FatalLogs, err.Error())
	log.Printf(err.Error())
	panic("fatal interruption")
}
