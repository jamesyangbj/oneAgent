package net

import "net"

type TestConn struct {
	c net.Conn
}

func NewConn(c net.Conn) *FrontConn {
	//return  new(TestConn)
	return nil
}

func (m *TestConn) Close() {

}

func (m *TestConn) Read(c FrontConn, buf []byte, maxSize int) {

}

func (m *TestConn) Write(c FrontConn, buf []byte, maxSize int) {

}
