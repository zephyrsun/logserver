package logserver

import (
	"fmt"
)

type ConsoleWriter struct {
}

func (this *ConsoleWriter) Write(k string, b []byte) {
	fmt.Println(logType[k] + ":" + string(b))
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}
