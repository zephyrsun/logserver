package logserver

import (
	"os"
	"path/filepath"
	"fmt"
	"time"
)

const EOL = "\n"

type Logger struct {
	out     *os.File
}

func (this *Logger) Write(b []byte) error {

	b = append(b, EOL...)

	//dump("%s", b)

	_, err := this.out.Write(b)

	this.Error(err)//panicOnError(err)

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
	//dump("dir:%s", f)

	//create dir first
	os.Mkdir(filepath.Dir(f), os.ModePerm)
	//panicOnError(err)

	out, err := os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	panicOnError(err)

	return &Logger{out}
}
