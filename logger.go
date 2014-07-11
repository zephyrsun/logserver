package logserver

import (
	"os"
	"path/filepath"
	"fmt"
	"time"
	"sync"
)

const EOL = "\n"

type Logger struct {
	out     *os.File
}

var mutex *sync.Mutex

func (this *Logger) Write(b []byte) error {

	b = append(b, EOL...)

	//Dump("%s", b)

	mutex.Lock()
	_, err := this.out.Write(b)
	mutex.Unlock()

	this.Error(err)//PanicOnError(err)

	return err
}

func (this *Logger) Error(err error) {
	if err != nil {
		this.Log("Error occoured:%s", err.Error())

		os.Exit(1)
	}
}

func (this *Logger) Log(format string, a ...interface{}) {
	this.Write([]byte("[" + time.Now().String() + "]" + fmt.Sprintf(format, a...)))
}

func (this *Logger) Close() {
	this.out.Close()
}

func NewLogger(f string) (*Logger) {
	//Dump("dir:%s", f)

	//create dir first
	os.Mkdir(filepath.Dir(f), os.ModePerm)
	//PanicOnError(err)

	out, err := os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	PanicOnError(err)

	return &Logger{out}
}

func init() {
	mutex = &sync.Mutex{}
}

