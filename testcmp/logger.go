package testcmp

import (
	"fmt"
	"log"
)

type Log struct {
	PrintLogs []string
	FatalLogs []string
}

func (l *Log) Printf(bundle string, format string, values ...interface{}) {
	s := fmt.Sprintf(format, values...)
	l.PrintLogs = append(l.PrintLogs, s)
	log.Printf(s)
}

func (l *Log) Fatal(err error) {
	l.FatalLogs = append(l.FatalLogs, err.Error())
	log.Printf(err.Error())
	panic("fatal interruption")
}
