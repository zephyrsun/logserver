package server

import (
	"log"
	"logserver/config"
	"logserver/logger"
)

func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Server interface {
	Listen()
}

var Config *config.Config

func Start(cfg *config.Config) {

	Config = cfg

	var s Server

	switch Config.Network {
	case "udp", "udp4", "udp6":
		s = &UDPServer{}
	case "tcp", "tcp4", "tcp6":
		s = &TCPServer{}
	case "http":
		s = &HTTPServer{}
	}

	log.Printf("start %s %s", Config.Network, Config.Address)

	s.Listen()
}

type Log struct {
	// 日志消息
	logs chan []byte
}

func (r *Log) Write() {

	var l logger.Logger

	switch Config.Logger {
	case "file":
		l = &logger.File{}
	}

	r.logs = make(chan []byte, Config.LogChanSize)

	for {
		rec := <-r.logs

		//if !utf8.Valid(rec) {
		//	log.Printf("wrong format: %s", rec)
		//	return
		//}

		l.Write(rec)
	}
}
