package testcmp

import "fmt"

type FmtLog struct {
	PrintLogs []string
	FatalLogs []string
}

func (l *FmtLog) Printf(format string, values ...interface{}) {
	s := fmt.Sprintf(format, values...)
	l.PrintLogs = append(l.PrintLogs, s)
	fmt.Printf(s)
}

func (l *FmtLog) Fatal(err error) {
	l.FatalLogs = append(l.FatalLogs, err.Error())
	fmt.Printf(err.Error())
	panic("fatal interruption")
}
