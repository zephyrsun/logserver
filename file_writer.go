package logserver

import (
	"os"
	"path"
	"time"
)

const (
	eol     = "\n"
	bufSize = 10 * 1024 * 1024
)

func newFile(name string) ( *os.File, error) {
	os.MkdirAll(path.Dir(name), 0755)
	//PanicOnError(err)
	return os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)
}

type FileWriter struct {
	lastHour int
	wr    map[string]*os.File
}

func (this *FileWriter) Write(k string, b []byte) {
	_, err := this.wr[k].Write(append(b, eol...))
	DumpError(err, false)
}

func (this *FileWriter) Rotate(now time.Time) {
	h := now.Hour() //this.timeNow.Format("2006-01-02-03")
	if this.lastHour == h {
		return
	}

	//this.lastHour = h
	//print(h)

	oldMap := make([]*os.File, len(logType))
	for k, v := range logType {

		filename := Config["save_dir"] + v + "_" + now.Format("2006-01-02-15") + ".log"

		new, err := newFile(filename)
		if err == nil {

			old, ok := this.wr[k]

			this.wr[k] = new

			if ok {
				//old.Close()
				oldMap = append(oldMap, old)
			}
		}
	}

	for _, old := range oldMap {
		old.Close()
	}
}

func NewFileWriter() *FileWriter {
	fw := &FileWriter{-1, make(map[string]*os.File)}

	//fw.files = make(map[string]*os.File)
	//fw.writers = make(map[string]*bufio.Writer)
	//fw.listenExit()

	fw.Rotate(time.Now())
	go Ticker(1*time.Second, fw.Rotate)

	return fw
}

