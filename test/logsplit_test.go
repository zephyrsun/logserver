package test

import (
	"testing"
	"github.com/zephyrsun/logserver"
)

func BenchmarkServer(b *testing.B) {
	srv := logserver.New()
	str :="1=2014-07-10 10:47:08|kingnetdc|activation|600||||1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||&2=2014-07-10 10:47:08|700|1|||||1|||||||||1111111111|2222222222|3333333333|from|tttttttt||||||||||||||||"
	strb := []byte(str)

	for i := 0; i < b.N; i++ { //use b.N for looping
		srv.Split(strb)
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

