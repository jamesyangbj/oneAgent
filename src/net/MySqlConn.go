package net

import "net"

type MysqlConn struct {
}

func (m MysqlConn) NewConn(c net.Conn) *Conn {
	return nil
}

func (m MysqlConn) Close() {

}

func (m MysqlConn) Read(c Conn, buf []byte, maxSize int) {

}

func (m MysqlConn) Write(c Conn, buf []byte, maxSize int) {

}
