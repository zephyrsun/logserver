package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Network string `json:"network"`
	Address string `json:"address"`
	//UDP，TCP读取chan大小
	ReadChanSize int `json:"read_chan_size"`
	//日志处理slice大小
	LogChanSize int `json:"log_chan_size"`
	//UDP接收缓冲区大小
	NetReadBuffer int `json:"net_read_buffer"`
	//记录类型：file,cache,
	Logger string `json:"logger"`
	//使用file时，日志目录
	LogFile string `json:"log_file"`
	//日志滚动参数：hourly,daily
	LogFileRotate string `json:"log_file_rotate"`
}

var Server = &Config{
	Network:       "udp",
	Address:       ":1982",
	ReadChanSize:  1 << 10,
	LogChanSize:   1 << 16,
	NetReadBuffer: 1 << 20,
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

func SetConfig(cfg *Config) {
	Server = cfg
}
