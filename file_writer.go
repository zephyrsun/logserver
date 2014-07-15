package logserver

import (
	"os"
	"path"
	"time"
	"bufio"
	"os/signal"
	"syscall"
)

const (
	eol     = "\n"
	bufSize = 8192
)

var (

)

type FileWriter struct {
	files    map[string]*os.File
	writers map[string]*bufio.Writer
	lastHour int
}

func (this *FileWriter) Write(k string, b []byte) error {
	b = append(b, eol...)

	_, err := this.files[k].Write(b)

	DumpError(err, false)

	return err
}

func (this *FileWriter) Rotate(now time.Time) {

	//this.Flush()

	h := now.Hour() //this.timeNow.Format("2006-01-02-03")
	if this.lastHour == h {
		return
	}

	this.lastHour = h

	for k, v := range logType {

		filename := Config["save_dir"] + v + "_" + now.Format("2006-01-02-15") + ".log"

		//new file
		f, err := newFile(filename)
		if err != nil {
			continue
		}

		//new writer
		/*
		ow, ok := this.writers[k]
		if ok {
			ow.Flush()
			ow.Reset(f)
		}else {
			this.writers[k] = bufio.NewWriterSize(f, bufSize)
		}
		*/

		// close file
		of, ok := this.files[k]
		if ok {
			of.Close()
		}

		this.files[k] = f
	}
}

func (this *FileWriter) Flush() {
	for k, _ := range logType {
		w, ok := this.writers[k]
		if ok {
			w.Flush()
		}
	}
}

func (this *FileWriter) listenExit() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		syscall.SIGINT)

	go func() {
		<-ch
		this.Flush()

		os.Exit(1)
	}()
}

func newFile(f string) ( *os.File, error) {
	os.MkdirAll(path.Dir(f), 0755)
	//PanicOnError(err)

	file, err := os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)

	DumpError(err, false)

	return file, err
}

func NewFileWriter() *FileWriter {
	fw := &FileWriter{}

	fw.lastHour = -1

	fw.files = make(map[string]*os.File)
	fw.writers = make(map[string]*bufio.Writer)

	//fw.listenExit()

	fw.Rotate(time.Now())
	go Ticker(1*time.Second, fw.Rotate)

	return fw
}
