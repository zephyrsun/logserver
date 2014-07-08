package logserver

import (
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
)

func panicOnError(err error) {
	if err != nil {
		//dump("Error occoured:%s", err.Error())
		//os.Exit(1)

		panic(err)
	}
}

func dump(format string, a ...interface{}) {
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

func loadConfig(file string) (config jsonConfigType) {

	b, err := ioutil.ReadFile(file)
	panicOnError(err)

	err = json.Unmarshal(b, &config)
	panicOnError(err)

	return
}

