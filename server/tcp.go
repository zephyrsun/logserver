package server

import (
	"logserver/config"
	"logserver/util"
	"net"
)

type TCPServer struct {
	conn *net.TCPListener
	Log
}

func (s *TCPServer) Listen() {
	la, err := net.ResolveTCPAddr(config.Server.Network, config.Server.Address)
	util.Fatal(err)

	s.conn, err = net.ListenTCP(config.Server.Network, la)

	defer s.Close()

	go s.Write()

	s.Read()
}

func (s *TCPServer) Read() {
	for {
		c, err := s.conn.Accept()
		if err != nil {
			util.Print("TCP Accept error:%s", err)
			return
		}

		go s.read(c)
	}
}

func (s *TCPServer) read(c net.Conn) {
	defer c.Close()

	buf := make([]byte, config.Server.ReadChanSize)
	for {
		n, err := c.Read(buf)
		if err != nil {
			util.Print("TCP Read error:%s", err)
			return
		}

		s.logs <- buf[:n]
	}
}

func (s *TCPServer) Close() {
	err := s.conn.Close()
	util.Fatal(err)
}
