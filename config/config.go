package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Network string `json:"network"`
	Address string `json:"address"`
	//UDP接收缓冲区大小
	ReceiveBuffer int `json:"receive_buffer"`
	//读取buffer大小（能接收的数据长度）
	ReadBuffer int `json:"read_buffer"`
	//写入chan大小（写入队列长度）
	WriteChanSize int `json:"write_chan_size"`
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
	ReceiveBuffer: 1 << 20,
	ReadBuffer:    1 << 16,
	WriteChanSize: 1 << 20,
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
