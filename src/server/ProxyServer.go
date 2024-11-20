package server

import (
	"fmt"
	"net"
	mysql "oneProxy/src/mysql"
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
		return
	}

	for {
		c, _ := s.Accept()

		switch p.Type {
		case "mysql":
			break
		default:
			n := mysql.NewConn(c)
			go processConn(n)
			break
		}

	}

}

func processConn(c *mysql.FrontConn) {
	if err := c.Handshake(); err != nil {
		fmt.Println("fail to handshake with client , to close the conn.")
		c.Close()
		return
	}

	for {
		if err := c.HandleData(); err != nil {
			fmt.Println("fail to handle cmd data with client , to close the conn.")
			c.Close()
			return
		}
	}
}
