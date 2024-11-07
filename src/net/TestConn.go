package net

import "net"

type TestConn struct {
	c net.Conn
}

func NewConn(c net.Conn) *Conn {
	//return  new(TestConn)
	return nil
}

func (m *TestConn) Close() {

}

func (m *TestConn) Read(c Conn, buf []byte, maxSize int) {

}

func (m *TestConn) Write(c Conn, buf []byte, maxSize int) {

}
