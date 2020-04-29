package main

import (
	"flag"
	"logserver/config"
	"logserver/reader"
)

func main() {
	file := flag.String("c", "config.json", "configuration file")

	flag.Parse()

	cfg := config.Load(*file)

	server.Start(cfg)
}
