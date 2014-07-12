package logserver

import (
	"os"
	"path"
	"time"
	"bufio"
)

const (
	eol     = "\n"
	bufSize = 8192
)

var (
	lastHour int
)

type FileWriter struct {
	files    map[string]*os.File
	writers map[string]*bufio.Writer
	timeNow time.Time
}

func (this *FileWriter) Write(k string, b []byte) error {

	b = append(b, eol...)

	_, err := this.files[k].Write(b)

	ErrorHandler(err)

	return err
}

func (this *FileWriter) rotate(now time.Time) {
	h := now.Hour() //this.timeNow.Format("2006-01-02-03")
	if lastHour != h {

		lastHour = h

		for k, v := range logType {
			file := Config["save_dir"] + v + "_" + this.timeNow.Format("2006-01-02-15") + ".log"

			//new file
			f, err := newFile(file)
			if err != nil {
				continue
			}

			//new writer
			ow, ok := this.writers[k]
			if ok {
				ow.Flush()
				ow.Reset(f)
			}else {
				this.writers[k] = bufio.NewWriterSize(f, bufSize)
			}

			// close file
			of, ok := this.files[k]
			if ok {
				of.Close()
			}

			this.files[k] = f
		}
	}
}

func newFile(f string) ( *os.File, error) {
	os.MkdirAll(path.Dir(f), 0755)
	//PanicOnError(err)

	file, err := os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)

	ErrorHandler(err)

	return file, err
}

func NewFileWriter() *FileWriter {
	fw := &FileWriter{}

	fw.files = make(map[string]*os.File)
	fw.writers = make(map[string]*bufio.Writer)

	fw.rotate(time.Now())

	go Ticker(1*time.Second, fw.rotate)

	return fw
}
