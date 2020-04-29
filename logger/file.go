package logger

import (
	"logserver/config"
	"logserver/util"
	"os"
	"strings"
	"time"
)

type FileLogger struct {
	handler *os.File
}

const EOL = '\n'

var (
	lastT  = -1
	rotate = make(chan bool)
)

func (l *FileLogger) Init() {
	l.rotateFile()

	t := time.NewTicker(time.Second)

	go func() { //按小时滚动日志
		for {
			c := <-t.C

			t := c.Day()
			if config.Server.LogFileHourly {
				t = c.Hour()
			}

			if t == lastT {
				return
			}

			l.rotateFile()
		}

	}()
}

func (l *FileLogger) Write(b []byte) {
	if l.handler == nil {
		util.Print("file handler error")
		return
	}

	_, err := l.handler.Write(append(b, EOL))
	util.Error("file write error: %s", err)
}

func (l *FileLogger) rotateFile() {
	if l.handler != nil {
		err := l.handler.Close()
		util.Error("file close error: %s", err)
		l.handler = nil
	}

	dir := strings.TrimSuffix(config.Server.LogFileDir, "/")

	format := "2006-01-02"
	if config.Server.LogFileHourly {
		format = "2006-01-02-15"
	}

	filename := dir + "/" + time.Now().Format(format) + ".log"

	err := os.MkdirAll(dir, 0744)
	util.Fatal(err)

	l.handler, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	util.Fatal(err)
}
