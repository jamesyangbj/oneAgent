package net

import "net"

type BackendConn struct {
}

func (m BackendConn) NewConn(c net.Conn) *BackendConn {
	return nil
}

func (m BackendConn) Close() {

}

func (m BackendConn) Read(c FrontConn, buf []byte, maxSize int) {

}

func (m BackendConn) Write(c FrontConn, buf []byte, maxSize int) {

}
