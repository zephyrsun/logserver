package server

import (
	"logserver/config"
	"logserver/logger"
	"logserver/util"
)

type Server interface {
	Listen()
}

var l logger.Logger

func Start() {
	initLogger()
	initServer()
}

func initLogger() {
	switch config.Server.Logger {
	case "file":
		l = &logger.FileLogger{}
		l.Init()
	}
}

func initServer() {

	var s Server

	switch config.Server.Network {
	case "udp", "udp4", "udp6":
		s = &UDPServer{}
	case "tcp", "tcp4", "tcp6":
		s = &TCPServer{}
	case "http":
		s = &HTTPServer{}
	}

	util.Print("start %s %s", config.Server.Network, config.Server.Address)

	s.Listen()
}

type Log struct {
	// 日志消息
	logs chan []byte
}

func (r *Log) Write() {

	r.logs = make(chan []byte, config.Server.LogChanSize)

	for {
		rec := <-r.logs

		//if !utf8.Valid(rec) {
		//	log.Printf("wrong format: %s", rec)
		//	return
		//}

		l.Write(rec)
	}
}
