package server

import (
	"logserver/config"
	"logserver/logger"
	"logserver/util"
	"os"
	"os/signal"
	"syscall"
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
		l = logger.NewFileLogger()
	}

	listenExit(l)
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

	util.Printf("start %s %s", config.Server.Network, config.Server.Address)

	s.Listen()
}

func listenExit(l logger.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		for {
			<-c
			l.Close()
		}
	}()
}

type Log struct {
	// 日志消息
	logs chan []byte
}

func (r *Log) Write() {
	r.logs = make(chan []byte, config.Server.WriteChanSize)

	for {
		l.Write(<-r.logs)
	}
}
