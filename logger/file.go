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
	writer     *os.File
	timeFlag   int
	fileFormat string
}

const (
	eol         = '\n'
	errorRotate = "log_file_rotate accepts: {daily|hourly}"
)

func NewFileLogger() *FileLogger {
	l := &FileLogger{
		timeFlag: -1,
	}

	l.initFormat()

	go l.initTicker()

	return l
}

func (l *FileLogger) Write(b []byte) {
	if l.writer == nil {
		util.Printf("file handler error")
		return
	}

	_, err := l.writer.Write(append(b, eol))
	util.Error("file write error: %s", err)
}

func (l *FileLogger) Close() {
	if l.writer != nil {
		err := l.writer.Sync()
		util.Error("file sync error: %s", err)

		err = l.writer.Close()
		util.Error("file close error: %s", err)

		l.writer = nil
	}
}

func (l *FileLogger) openFile(t time.Time) {
	l.Close()

	filename := fmt.Sprintf(config.Server.LogFile, t.Format(l.fileFormat))

	err := os.MkdirAll(path.Dir(filename), 0744)
	util.Fatal(err)

	l.writer, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	util.Fatal(err)
}

func (l *FileLogger) initTicker() {
	var c int
	ticker := time.NewTicker(time.Second)

	for {
		t := <-ticker.C

		switch config.Server.LogFileRotate {
		case "hourly":
			c = t.Hour()
		case "daily":
			c = t.Day()
		default:
			util.Fatal(errorRotate)
		}

		if c != l.timeFlag {
			l.timeFlag = c
			l.openFile(t)
		}
	}
}

func (l *FileLogger) initFormat() {
	switch config.Server.LogFileRotate {
	case "hourly":
		l.fileFormat = "2006-01-02-15"
	case "daily":
		l.fileFormat = "2006-01-02"
	default:
		util.Fatal(errorRotate)
	}
}
