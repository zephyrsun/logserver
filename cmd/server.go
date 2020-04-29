package main

import (
	"flag"
	"logserver/config"
	"logserver/server"
)

func main() {

	file := flag.String("c", "config.json", "configuration file")

	flag.Parse()

	config.Load(*file)

	server.Start()
}
