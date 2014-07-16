package logserver

import (
	"fmt"
)

type ConsoleWriter struct {
}

func (this *ConsoleWriter) Write(k string, b []byte) (int, error) {
	return fmt.Println(logType[k] + ":" + string(b))
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}
