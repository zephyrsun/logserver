package main

import (
	"net"
	"github.com/zephyrsun/logserver"
	"flag"
)

func main() {

	n := flag.Int("n", 100000, "Number of requests to perform")
	d := flag.String("d", "1=2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 14:14:58|2014-07-10 13:57:40|200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||", "Data to be sent")

	flag.Parse()

	conn, err := net.Dial("udp", ":8282")
	logserver.PanicOnError(err)

	b := []byte(d)

	for i := 0; i < *n; i++ {

		d, _ := conn.Write(b)

		//logserver.PanicOnError(err)
		logserver.Dump("Bytes was sent:%d", d)
	}

	logserver.Dump("Done!")
}
