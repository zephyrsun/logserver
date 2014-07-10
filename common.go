package logserver

import (
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
)

type configType map[string]string

var config = configType{
	"address": ":8282",
	"sys_log": "syslog/logserver.log",
	"save_dir": "data/",
}

func loadConfig(file string) {
	b, err := ioutil.ReadFile(file)
	//PanicOnError(err)
	if err == nil {
		err = json.Unmarshal(b, &config)
		PanicOnError(err)
	}
}

func PanicOnError(err error) {
	if err != nil {
		//Dump("Error occoured:%s", err.Error())
		//os.Exit(1)
		panic(err)
	}
}

func Dump(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

func ticker(sec time.Duration, callback func(time.Time) bool) {
	c := time.Tick(sec * time.Second)
	for now := range c {
		if callback(now) == false {
			break
		}
	}
}

