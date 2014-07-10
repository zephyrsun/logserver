package main

import (
	"net"
	"github.com/zephyrsun/logserver"
	"flag"
	"time"
)

func main() {

	poolNum := 8

	var pool [8]net.Conn


	addr := flag.String("addr", ":8282", "Server address")
	n := flag.Int("n", 100000, "Number of requests to perform")
	d := flag.String("d", "1=2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 14:14:58|2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||", "Data to be sent")

	flag.Parse()

	// send data
	b := []byte(*d)

	// init connection pool
	for i := 0; i < poolNum; i++ {
		conn, err := net.Dial("udp", *addr)
		logserver.PanicOnError(err)

		pool[i] = conn
	}

	// seed
	//rand.Seed(time.Now().UnixNano())

	ci := make(chan int)

	go func() {

		for{
			i:=<-ci % poolNum
			pool[i].Write(b)
		}
	}()

	t := time.Now()
	for i := 0; i < *n; i++ {
		ci<-i
	}

	logserver.Dump("Done! time:%s", time.Now().Sub(t))
}
