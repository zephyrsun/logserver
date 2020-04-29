package server

import (
	"log"
	"net"
)

type UDPServer struct {
	conn *net.UDPConn
	Log
}

func (s *UDPServer) Listen() {
	la, err := net.ResolveUDPAddr(Config.Network, Config.Address)
	Fatal(err)

	s.conn, err = net.ListenUDP("udp", la)
	Fatal(err)

	if Config.ReadBuffer > 0 {
		err = s.conn.SetReadBuffer(Config.ReadBuffer)
		Fatal(err)
	}

	defer s.Close()

	go s.Write()

	s.Read()
}

func (s *UDPServer) Read() {
	for {
		buf := make([]byte, Config.ReadChanSize)
		n, err := s.conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		if n > 0 {
			s.logs <- buf[:n]
		}
	}
}

func (s *UDPServer) Close() {
	err := s.conn.Close()
	Fatal(err)
}
