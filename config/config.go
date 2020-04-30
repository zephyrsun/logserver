package config

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Network string `json:"network"`
	Address string `json:"address"`
	//UDP，TCP读取chan大小
	ReadChanSize int `json:"read_chan_size"`
	//日志处理slice大小
	LogChanSize int `json:"log_chan_size"`
	//UDP接收缓冲区大小
	ReadBuffer int `json:"read_buffer"`
	//记录类型：file,cache,
	Logger string `json:"logger"`
	//使用file时，日志目录
	LogFile string `json:"log_file"`
	//日志滚动参数：hourly,daily
	LogFileRotate string `json:"log_file_rotate"`
}

var Server = &config{
	Network:       "udp",
	Address:       ":8282",
	LogChanSize:   1 << 10,
	ReadChanSize:  1 << 10,
	ReadBuffer:    1 << 20,
	Logger:        "file",
	LogFile:       "./logs/log-%s.log",
	LogFileRotate: "daily",
}

func Load(filename string) {

	b, err := ioutil.ReadFile(filename)

	if err == nil {
		err = json.Unmarshal(b, Server)
	}
}
