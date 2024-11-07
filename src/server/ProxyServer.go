package server

import (
	"fmt"
	"net"
	net2 "oneProxy/src/net"
)

type ProxyServer struct {
	Type        string
	ListenHost  string
	Port        int
	MaxClient   int
	ReadBufSize int
}

func (p *ProxyServer) NewServer() *ProxyServer {
	s := new(ProxyServer)
	s.Type = "mysql"
	s.ListenHost = "0.0.0.0"
	s.Port = 9527
	s.MaxClient = 1024
	s.ReadBufSize = 4 * 1024
	return s
}

func (p *ProxyServer) Start() {
	s, err := net.Listen("tcp", fmt.Sprintf("%s:%d", p.ListenHost, p.Port))
	if err != nil {

	}
	for {
		c, _ := s.Accept()

		switch p.Type {
		case "mysql":
			break
		default:
			n := net2.NewConn(c)
			go processConn(n)
			break

		}

	}

}

func processConn(c *net2.Conn) {

}
