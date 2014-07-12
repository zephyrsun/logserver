package logserver

import (
	"fmt"
	"time"
	"log"
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
	if err == nil {
		err = json.Unmarshal(b, &Config)
		ErrorHandler(err)
	}
}

func ErrorHandler(err error) {
	if err != nil {
		log.Println(err)
	}
}

func Dump(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func DumpError(err error) {
	if err != nil {
		panic(err)
	}
}

func Ticker(sec time.Duration, callback func(time.Time)) {
	c := time.Tick(sec)
	for now := range c {
		callback(now)
	}
}

