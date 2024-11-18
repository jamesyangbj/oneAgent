package net

import "net"

type FrontConn interface {
	NewConn(c net.Conn) *FrontConn
	Close()
	Read(c FrontConn, buf []byte, maxSize int)
	Write(c FrontConn, buf []byte, maxSize int)
}
