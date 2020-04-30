package test

import (
	"net"
	"testing"
)

var data = []byte("1=200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt|||||||||||||||1|&2=200|1|2|3|4|5|6|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||")

func TestUdpClient(t *testing.T) {

	conn, err := net.Dial("udp", ":1982")
	if err != nil {
		t.Fatal(err)
	}

	defer conn.Close()

	for i := 0; i < 1e5; i++ {
		_, err = conn.Write(data)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTcpClient(t *testing.T) {

	for i := 0; i < 1e3; i++ {
		conn, err := net.Dial("tcp", ":1982")
		if err != nil {
			t.Fatal(err)
		}

		_, err = conn.Write(data)
		if err != nil {
			t.Fatal(err)
		}

		_ = conn.Close()
	}
}
