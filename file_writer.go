package logserver

import (
	"os"
	"path"
	"time"
)

const (
	eol     = "\n"
	bufSize = 8192
)

type FileWriter struct {
	writers    map[string]*os.File
	//writers map[string]*bufio.Writer
	lastHour int
}

func (this *FileWriter) Write(k string, b []byte) (int, error) {
	b = append(b, eol...)

	return this.writers[k].Write(b)
}

func (this *FileWriter) Rotate(now time.Time) {
	h := now.Hour() //this.timeNow.Format("2006-01-02-03")
	if this.lastHour == h {
		return
	}

	this.lastHour = h
	//print(h)

	for k, v := range logType {

		filename := Config["save_dir"] + v + "_" + now.Format("2006-01-02-15") + ".log"

		//new file
		f, err := newFile(filename)
		DumpError(err, false)

		this.writers[k] = f
	}
}

func newFile(f string) ( *os.File, error) {
	os.MkdirAll(path.Dir(f), 0755)
	//PanicOnError(err)
	return os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)
}

func NewFileWriter() *FileWriter {
	fw := &FileWriter{
		lastHour:-1,
		writers:make(map[string]*os.File),
	}

	//fw.files = make(map[string]*os.File)
	//fw.writers = make(map[string]*bufio.Writer)
	//fw.listenExit()

	fw.Rotate(time.Now())
	go Ticker(1*time.Second, fw.Rotate)

	return fw
}
