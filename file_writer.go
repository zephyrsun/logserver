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

/*
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
		buf:make([]byte, 20*bufSize),
	}

	w.Rotate(time.Now())
	return w
}
*/

type FileWriter struct {
	lastHour int
	wr    map[string]*os.File
	old    map[string]*os.File
}

func (o *FileWriter) Write(k string, b []byte) {
	b = append(b, eol...)
	o.wr[k].Write(b)
}

func (o *FileWriter) Rotate(now time.Time) {
	h := now.Hour()
	if o.lastHour == h {
		return
	}

	o.lastHour = h
	//print(h)

	for k, name := range logType {
		//o.wr[k].Rotate(now)

		f := Config["save_dir"] + name + "_" + now.Format("2006-01-02-15") + ".log"

		o.old[k], o.wr[k] = o.wr[k], newFile(f)
	}

	for _, wr := range o.old {
		wr.Close()
	}
}

func NewFileWriter() *FileWriter {

	/*
	wr := make(map[string]*bufWriter, len(logType))
	for k, v := range logType {
		wr[k] = newBufWriter(v)
	}
	*/
	wr := make(map[string]*os.File, len(logType))
	old := make(map[string]*os.File, len(logType))

	fw := &FileWriter{-1, wr, old}

	fw.Rotate(time.Now())

	go Ticker(1*time.Second, fw.Rotate)

	return fw
}

