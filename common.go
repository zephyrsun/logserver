package logserver

import (
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
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

func FormatLog() {

}

func Ticker(sec time.Duration, callback func(time.Time)) {
	c := time.Tick(sec)
	for now := range c {
		callback(now)
	}
}

