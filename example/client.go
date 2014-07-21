package main

import (
	"net"
	ls "github.com/zephyrsun/logserver"
	"flag"
	"time"
)

func main() {

	addr := flag.String("addr", ":8282", "Server address")
	n := flag.Int("n", 100000, "Number of requests to perform")
	d := flag.String("d", "1=2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 14:14:58|2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||", "Data to be sent")

	flag.Parse()

	b := []byte(*d)

	la, err := net.ResolveUDPAddr("udp", *addr)
	ls.DumpError(err, true)

	conn, err := net.DialUDP("udp", nil, la)
	ls.DumpError(err, true)

	t := time.Now()
	i := 0
	for ; i < *n; i++ {
		_, err := conn.Write(b)
		ls.DumpError(err, true)
	}

	ls.Dump("Done! requested:%d, time:%s", i, time.Now().Sub(t))
}
