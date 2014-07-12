package main

import (
	"net"
	"github.com/zephyrsun/logserver"
	"flag"
	"time"
)

func main() {

	addr := flag.String("addr", ":8282", "Server address")
	r := flag.Int("r", 2, "Number of routing")
	n := flag.Int("n", 100000, "Number of requests to perform")
	d := flag.String("d", "1=2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 14:14:58|2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||", "Data to be sent")

	flag.Parse()

	b := []byte(*d)

	t := time.Now()

	c := make(chan int, *r)

	for i := 0; i < *r; i++ {
		go func() {
			conn, err := net.Dial("udp", *addr)
			logserver.PanicOnError(err)

			i := 0
			for ; i < *n/(*r); i++ {
				_, err := conn.Write(b)
				logserver.PanicOnError(err)
			}

			c<-i
		}()
		logserver.Dump("written: %d", <-c)
	}

	logserver.Dump("Done! time:%s", time.Now().Sub(t))
}
