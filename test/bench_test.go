package test

import (
	"testing"
	"github.com/zephyrsun/logserver"
)

func BenchmarkParse(b *testing.B) {
	srv := logserver.New()

	str := "1=2014-07-10 10:47:08|kingnetdc|activation|600||||1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 10:47:08|700|1|||||1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||"
	strByte := []byte(str)

	callback := func(string, []byte) {

	}

	for i := 0; i < b.N; i++ { //use b.N for looping
		srv.Parse(strByte, callback)
	}
}

func BenchmarkWrite(b *testing.B) {
	srv := logserver.New()

	k1 := "1"
	str1 := "2014-07-10 10:47:08|kingnetdc|activation|600||||1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||"
	strByte1 := []byte(str1)

	k2 := "2"
	str2 := "2014-07-10 10:47:08|700|1|||||1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||"
	strByte2 := []byte(str2)

	c1 := make(chan int)
	c2 := make(chan int)

	go func() {
		for {
			<-c1
			srv.Write(k1, strByte1)
		}
	}()

	go func() {
		for {
			<-c2
			srv.Write(k2, strByte2)
		}
	}()

	for i := 0; i < b.N; i++ { //use b.N for looping
		c1<-i
		c2<-i
	}
}

/*
func BenchmarkServer2(b *testing.B) {
	srv := logserver.New()
	str :="1=2014-07-10 10:47:08|kingnetdc|activation|600||||1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 10:47:08|700|1|||||1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||"
	strb := []byte(str)

	for i := 0; i < b.N; i++ { //use b.N for looping
		srv.Split(strb)
	}
}

func BenchmarkServer3(b *testing.B) {
	srv := logserver.New()
	str :="3=2014-07-10 10:47:08|400|1|2|3|4|1|1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||"
	strb := []byte(str)

	for i := 0; i < b.N; i++ { //use b.N for looping
		srv.Split(strb)
	}
}
*/

