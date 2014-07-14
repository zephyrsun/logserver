package main

import (
	"net"
	"github.com/zephyrsun/logserver"
	"flag"
	"time"
)

func main() {

	addr := flag.String("addr", ":8282", "Server address")
	c := flag.Int("c", 10000, "Number of multiple requests to make")
	n := flag.Int("n", 100000, "Number of requests to perform")
	d := flag.String("d", "1=2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 14:14:58|2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||", "Data to be sent")

	flag.Parse()

	b := []byte(*d)

	max := *n / (*c)

	doSend := func() {
		conn, err := net.Dial("udp", *addr)
		logserver.DumpError(err, true)

		i := 0
		for ; i < *c; i++ {
			_, err := conn.Write(b)
			logserver.DumpError(err, true)

			//logserver.Dump("%s,%s", conn, b)
		}

		logserver.Dump("requested number:%d", i)
	}

	t := time.Now()
	for i := 0; i < max; i++ {
		doSend()
	}
	logserver.Dump("Done! time:%s", time.Now().Sub(t))
}
