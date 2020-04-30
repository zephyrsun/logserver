package logger

import (
	"fmt"
	"logserver/config"
	"logserver/util"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	handler *os.File
}

const (
	eol         = '\n'
	errorRotate = "log_file_rotate accepts: {daily|hourly}"
)

var (
	lastT  = -1
	format string
)

func initFormat() {
	switch config.Server.LogFileRotate {
	case "hourly":
		format = "2006-01-02-15"
	case "daily":
		format = "2006-01-02"
	default:
		util.Fatal(errorRotate)
	}
}

func (l *FileLogger) initTicker() {
	var t int
	ticker := time.NewTicker(time.Second)

	for {
		c := <-ticker.C

		switch config.Server.LogFileRotate {
		case "hourly":
			t = c.Hour()
		case "daily":
			t = c.Day()
		default:
			util.Fatal(errorRotate)
		}

		if t == lastT {
			return
		}

		lastT = t

		l.rotateFile()
	}
}

func (l *FileLogger) Init() {
	initFormat()

	go l.initTicker()
}

func (l *FileLogger) Write(b []byte) {
	if l.handler == nil {
		util.Print("file handler error")
		return
	}

	_, err := l.handler.Write(append(b, eol))
	util.Error("file write error: %s", err)
}

func (l *FileLogger) rotateFile() {
	if l.handler != nil {
		err := l.handler.Close()
		util.Error("file close error: %s", err)
		l.handler = nil
	}

	filename := fmt.Sprintf(config.Server.LogFile, time.Now().Format(format))

	err := os.MkdirAll(path.Dir(filename), 0744)
	util.Fatal(err)

	l.handler, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	util.Fatal(err)
}
