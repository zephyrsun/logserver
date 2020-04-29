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
	LogFileDir string `json:"log_file_dir"`
	//true:按小时分割，false:按天分割
	LogFileHourly bool `json:"log_file_hourly"`
}

var Server = &config{
	Network:       "udp",
	Address:       ":8282",
	LogChanSize:   1 << 10,
	ReadChanSize:  1 << 10,
	ReadBuffer:    1 << 20,
	Logger:        "file",
	LogFileDir:    "../logs/",
	LogFileHourly: false,
}

func Load(filename string) {

	b, err := ioutil.ReadFile(filename)

	if err == nil {
		err = json.Unmarshal(b, Server)
	}
}
