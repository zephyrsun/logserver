package logserver

import (
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
	"os"
	"path"
)

const (
	eol     = "\n"
	bufSize = 1024 * 1024
)

var logType = map[string]string {
	"1":"login",
	"2":"act",
	"3":"pay",
	"4":"item",
	"5":"error",
	"6":"funel",
	"7":"att",
}

type configType map[string]string

var Config = configType{
	"address": ":8282",
	"save_dir": "data/",
	"writer":"",
}

func loadConfig(file string) {
	b, err := ioutil.ReadFile(file)
	DumpError(err, true)

	err = json.Unmarshal(b, &Config)
	DumpError(err, true)
}

func Dump(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func DumpError(err error, exit bool) {
	if err != nil {
		if exit {
			panic(err)
		}else {
			Dump("error: %s", err.Error())
		}
	}
}

func newLogFile(name string, t time.Time) *os.File {
	name = Config["save_dir"] + name + "_" + t.Format("2006-01-02-15") + ".log"

	os.MkdirAll(path.Dir(name), 0755)

	//PanicOnError(err)
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)
	DumpError(err, false)

	return f
}

func Ticker(sec time.Duration, callback func(time.Time)) {
	c := time.Tick(sec)
	for now := range c {
		callback(now)
	}
}

