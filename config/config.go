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
	ReadBuffer int    `json:"read_buffer"`
	Logger     string `json:"logger"`
}

func Load(filename string) *Config {

	cfg := &Config{
		Network:      "udp",
		Address:      ":8282",
		LogChanSize:  1 << 10,
		ReadChanSize: 1 << 10,
		ReadBuffer:   1 << 20,
		Logger:       "file",
	}

	b, err := ioutil.ReadFile(filename)

	if err == nil {
		err = json.Unmarshal(b, cfg)
	}

	return cfg
}
