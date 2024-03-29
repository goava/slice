package testcmp

import "fmt"

type FmtLog struct {
	PrintLogs []string
	FatalLogs []string
}

func (l *FmtLog) Printf(bundle string, format string, values ...interface{}) {
	s := fmt.Sprintf(format, values...)
	l.PrintLogs = append(l.PrintLogs, s)
	fmt.Println(s)
}

func (l *FmtLog) Fatal(err error) {
	l.FatalLogs = append(l.FatalLogs, err.Error())
	fmt.Println(err.Error())
	//panic("fatal interruption")
}
