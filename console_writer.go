package logserver

import (
	"fmt"
)

type ConsoleWriter struct {
}

func (this *ConsoleWriter) Write(k string, b []byte) error {
	_, err := fmt.Println(logType[k] + ":" + string(b))

	ErrorHandler(err)

	return err
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}
