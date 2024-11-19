package net

import "net"

type Conn interface {
	NewConn(c net.Conn) *Conn
	Close()
	ReadPacket() []byte
	WritePacket(buf []byte)
}
