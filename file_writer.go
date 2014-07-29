package logserver

import (
	"os"
	"time"
)

type FileWriter struct {
	lastHour int
	wr    map[string]*os.File
	oldWr    map[string]*os.File

}

func (o *FileWriter) Write(k string, b []byte) {

	_, err := o.wr[k].Write(b)
	DumpError(err, true)
}

func (o *FileWriter) Rotate(now time.Time) {
	h := now.Hour()
	if o.lastHour == h {
		return
	}

	o.lastHour = h
	//print(h)

	for k, name := range logType {
		o.oldWr[k], o.wr[k] = o.wr[k], newLogFile(name, now)
	}

	for _, wr := range o.oldWr {
		wr.Close()
	}
}

func NewFileWriter() *FileWriter {

	fw := &FileWriter{
		lastHour:-1,
		wr:make(map[string]*os.File, len(logType)),
		oldWr:make(map[string]*os.File, len(logType)),
	}

	fw.Rotate(time.Now())

	go Ticker(1*time.Second, fw.Rotate)

	return fw
}

