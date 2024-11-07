package net

import "net"

type Conn interface {
	NewConn(c net.Conn)
	Close()
	Read(c Conn, buf []byte, maxSize int)
	Write(c Conn, buf []byte, maxSize int)
}
