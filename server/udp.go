package server

import (
	"logserver/config"
	"logserver/util"
	"net"
)

type UDPServer struct {
	conn *net.UDPConn
	Log
}

func (s *UDPServer) Listen() {
	la, err := net.ResolveUDPAddr(config.Server.Network, config.Server.Address)
	util.Fatal(err)

	s.conn, err = net.ListenUDP("udp", la)
	util.Fatal(err)

	if config.Server.ReceiveBuffer > 0 {
		err = s.conn.SetReadBuffer(config.Server.ReceiveBuffer)
		util.Fatal(err)
	}

	defer s.Close()

	go s.Write()

	s.Read()
}

func (s *UDPServer) Read() {
	//var i = 0
	//go func() {
	//	t := time.Tick(time.Second)
	//	for _ = range t {
	//		util.Printf("i: %s", i)
	//	}
	//}()

	buf := make([]byte, config.Server.ReadBuffer)

	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			util.Printf("UDP Read error:%s", err)
			return
		}

		//i++

		if n > 0 {
			s.logs <- buf[:n]
		}
	}
}

func (s *UDPServer) Close() {
	err := s.conn.Close()
	util.Fatal(err)
}
