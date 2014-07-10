package main

import (
	"net"
	"github.com/zephyrsun/logserver"
	"flag"
	"time"
)

func main() {

	addr := flag.String("addr", ":8282", "Server address")
	c := flag.Int("c", 2, "Number of server connection to make")
	n := flag.Int("n", 100000, "Number of requests to perform")
	d := flag.String("d", "1=2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 14:14:58|2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||", "Data to be sent")

	num := *c

	pool := make(map[int]net.Conn)

	flag.Parse()

	// send data
	b := []byte(*d)

	// init connection pool
	for i := 0; i < num; i++ {
		conn, err := net.Dial("udp", *addr)
		logserver.PanicOnError(err)

		pool[i] = conn
	}

	// seed
	//rand.Seed(time.Now().UnixNano())

	ci := make(chan int)

	go func() {

		for {
			i := <-ci % num
			_, err := pool[i].Write(b)
			logserver.PanicOnError(err)
		}
	}()

	t := time.Now()
	for i := 0; i < *n; i++ {
		ci<-i
	}

	logserver.Dump("Done! time:%s", time.Now().Sub(t))
}
