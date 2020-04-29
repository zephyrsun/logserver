package server

import (
	"log"
	"net"
)

type TCPServer struct {
	conn *net.TCPListener
	Log
}

func (s *TCPServer) Listen() {
	la, err := net.ResolveTCPAddr(Config.Network, Config.Address)
	Fatal(err)

	s.conn, err = net.ListenTCP(Config.Network, la)

	defer s.Close()

	go s.Write()

	s.Read()
}

func (s *TCPServer) Read() {
	for {
		c, err := s.conn.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		go s.read(c)
	}
}

func (s *TCPServer) read(c net.Conn) {
	defer c.Close()

	buf := make([]byte, Config.ReadChanSize)
	for {
		n, err := c.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		s.logs <- buf[:n]
	}
}

func (s *TCPServer) Close() {
	err := s.conn.Close()
	Fatal(err)
}
