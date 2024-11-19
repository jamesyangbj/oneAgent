package net

import (
	"net"
	"time"
)

type FrontConn struct {
	c           net.Conn
	connectTime time.Time
	p           *PacketIO
}

func NewConn(c net.Conn) *FrontConn {
	f := new(FrontConn)
	f.c = c
	f.connectTime = time.Now()
	f.p = NewPacketIO(c)
	return f
}

func (m *FrontConn) Close() {

}

func (m *FrontConn) ReadPacket() []byte {

	return nil

}

func (m *FrontConn) WritePacket(data []byte) {

}
