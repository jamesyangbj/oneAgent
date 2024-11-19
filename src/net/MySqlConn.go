package net

import "net"

type BackendConn struct {
}

func (m *BackendConn) NewConn(c net.Conn) *BackendConn {

	return nil
}

func (m *BackendConn) Close() {

}

func (m *BackendConn) Read() []byte {

	return nil
}

func (m *BackendConn) Write(c Conn, data []byte) {

}
