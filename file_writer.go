package logserver

import (
	"os"
	"path"
	"time"
)

const (
	eol = "\n"
)

func newFile(name string) *os.File {
	os.MkdirAll(path.Dir(name), 0755)
	//PanicOnError(err)
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)
	DumpError(err, false)

	return f
}

type bufWriter struct{
	name string
	buf  []byte
	wr *os.File
}

func (o *bufWriter) Write(b []byte) {
	o.buf = append(o.buf, b...)
}

func (o *bufWriter) Flush() {
	o.wr.Write(o.buf)

	o.buf = o.buf[:0]
}

func (o *bufWriter) Close() {
	o.Flush()
	o.wr.Close()
}

func (o *bufWriter) Rotate(now time.Time) {
	f := Config["save_dir"] + o.name + "_" + now.Format("2006-01-02-15") + ".log"

	wr := newFile(f)

	o.Close()

	o.wr = wr
}

func newBufWriter(name string) *bufWriter {
	w := &bufWriter{
		name:name,
		buf:make([]byte, 10*bufSize),
	}

	w.Rotate(time.Now())
	return w
}

type FileWriter struct {
	lastHour int
	wr    map[string]*bufWriter
}

func (o *FileWriter) Write(k string, b []byte) {
	o.wr[k].Write(append(b, eol...))
}

func (o *FileWriter) Rotate(now time.Time)  {
	h := now.Hour()
	if o.lastHour == h {
		return
	}

	o.lastHour = h
	//print(h)

	for k, _ := range logType {
		o.wr[k].Rotate(now)
	}
}

func (o *FileWriter) Flush() {
	for k, _ := range logType {
		o.wr[k].Flush()
	}
}

func NewFileWriter() *FileWriter {

	wr := make(map[string]*bufWriter, len(logType))
	for k, v := range logType {
		wr[k] = newBufWriter(v)
	}

	fw := &FileWriter{-1, wr}

	return fw
}

