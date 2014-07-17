package logserver

import (
	"os"
	"path"
	"time"
	"bufio"
)

const (
	eol     = "\n"
	bufSize = 1024 * 1024
)

func newFile(name string) ( *os.File, error) {
	os.MkdirAll(path.Dir(name), 0755)
	//PanicOnError(err)
	return os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)
}

func NewBufWriter(name string) (*BufWriter, error) {
	f, err := newFile(name)
	if err != nil {
		return nil, err
	}

	return &BufWriter{f, bufio.NewWriterSize(f, bufSize)}, nil
}

type BufWriter struct{
	file *os.File
	writer *bufio.Writer
}

func (this *BufWriter) Write(b []byte) (int, error) {
	return this.writer.Write(b)
}

func (this *BufWriter) Close() error {
	return this.file.Close()
}

func (this *BufWriter) Flush() error {
	return this.writer.Flush()
}

type FileWriter struct {
	lastHour int
	writers    map[string]*BufWriter
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

		new, err := NewBufWriter(filename)
		if err == nil {
			old, ok := this.writers[k]
			if ok {
				old.Flush()
				old.Close()
			}

			this.writers[k] = new
		}
	}
}

func (this *FileWriter) Flush() {
	for k, _ := range logType {
		this.writers[k].Flush()
	}
}

func NewFileWriter() *FileWriter {
	fw := &FileWriter{-1, make(map[string]*BufWriter)}

	//fw.files = make(map[string]*os.File)
	//fw.writers = make(map[string]*bufio.Writer)
	//fw.listenExit()

	fw.Rotate(time.Now())
	go Ticker(1*time.Second, func(now time.Time) {
			fw.Rotate(now)
			fw.Flush()
		})

	return fw
}
