package mysql

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"oneProxy/src/config"
	"oneProxy/src/util"
)

type BackendConn struct {
	Host       string
	Port       int
	conn       net.Conn
	user       string
	password   string
	salt       []byte
	charset    string
	db         string
	pkg        *PacketIO
	pkgErr     error
	capability uint32
	status     uint16
	collation  config.CollationId
}

func (c *BackendConn) Connect(host string, port int, user string, password string, db string) {
	c.Host = host
	c.Port = port
	c.user = user
	c.password = password
	c.db = db

	c.ReConnect()
}

func NewBackendConn() *BackendConn {
	m := new(BackendConn)
	return m
}

func (c *BackendConn) ReConnect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return err
	}

	if c.conn != nil {
		oldConn := c.conn
		c.conn = conn
		oldConn.Close()
	} else {
		c.conn = conn
	}

	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetNoDelay(true)

	c.pkg = NewPacketIO(conn)

	if err := c.readInitialHandshake(); err != nil {
		c.conn.Close()
		return err
	}

	if err := c.writeAuthHandshake(); err != nil {
		c.conn.Close()

		return err
	}

	if _, err := c.readOK(); err != nil {
		c.conn.Close()

		return err
	}

	return nil

}

func (c *BackendConn) Close() {
	c.conn.Close()
}

func (c *BackendConn) readInitialHandshake() error {
	data, err := c.readPacket()
	if err != nil {
		return err
	}

	if data[0] == config.ERR_HEADER {
		return errors.New("read initial handshake error")
	}

	if data[0] < config.MinProtocolVersion {
		return fmt.Errorf("invalid protocol version %d, must >= 10", data[0])
	}

	//skip mysql version and connection id
	//mysql version end with 0x00
	//connection id length is 4
	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1 + 4

	c.salt = append(c.salt, data[pos:pos+8]...)

	//skip filter
	pos += 8 + 1

	//capability lower 2 bytes
	c.capability = uint32(binary.LittleEndian.Uint16(data[pos : pos+2]))

	pos += 2

	if len(data) > pos {
		//skip server charset
		//c.charset = data[pos]
		pos += 1

		c.status = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2

		c.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | c.capability

		pos += 2

		//skip auth data len or [00]
		//skip reserved (all [00])
		pos += 10 + 1

		// The documentation is ambiguous about the length.
		// The official Python library uses the fixed length 12
		// mysql-proxy also use 12
		// which is not documented but seems to work.
		c.salt = append(c.salt, data[pos:pos+12]...)
	}

	return nil
}

func (c *BackendConn) writeAuthHandshake() error {
	capability := config.CLIENT_PROTOCOL_41 | config.CLIENT_SECURE_CONNECTION |
		config.CLIENT_LONG_PASSWORD | config.CLIENT_TRANSACTIONS | config.CLIENT_LONG_FLAG

	capability &= c.capability

	//packet length
	//capbility 4
	//max-packet size 4
	//charset 1
	//reserved all[0] 23
	length := 4 + 4 + 1 + 23

	//username
	length += len(c.user) + 1

	//we only support secure connection
	auth := util.CalcPassword(c.salt, []byte(c.password))

	length += 1 + len(auth)

	if len(c.db) > 0 {
		capability |= config.CLIENT_CONNECT_WITH_DB

		length += len(c.db) + 1
	}

	c.capability = capability

	data := make([]byte, length+4)

	//capability [32 bit]
	data[4] = byte(capability)
	data[5] = byte(capability >> 8)
	data[6] = byte(capability >> 16)
	data[7] = byte(capability >> 24)

	//MaxPacketSize [32 bit] (none)
	//data[8] = 0x00
	//data[9] = 0x00
	//data[10] = 0x00
	//data[11] = 0x00

	//Charset [1 byte]
	data[12] = byte(c.collation)

	//Filler [23 bytes] (all 0x00)
	pos := 13 + 23

	//User [null terminated string]
	if len(c.user) > 0 {
		pos += copy(data[pos:], c.user)
	}
	//data[pos] = 0x00
	pos++

	// auth [length encoded integer]
	data[pos] = byte(len(auth))
	pos += 1 + copy(data[pos+1:], auth)

	// db [null terminated string]
	if len(c.db) > 0 {
		pos += copy(data[pos:], c.db)
		//data[pos] = 0x00
	}

	return c.writePacket(data)
}

func (c *BackendConn) readOK() (*Result, error) {
	data, err := c.readPacket()
	if err != nil {
		return nil, err
	}

	if data[0] == config.OK_HEADER {
		return c.handleOKPacket(data)
	} else if data[0] == config.ERR_HEADER {
		return nil, c.handleErrorPacket(data)
	} else {
		return nil, errors.New("invalid ok packet")
	}
	return nil, nil
}

func (c *BackendConn) readPacket() ([]byte, error) {
	d, err := c.pkg.ReadPacket()
	c.pkgErr = err
	return d, err
}

func (c *BackendConn) writePacket(data []byte) error {
	err := c.pkg.WritePacket(data)
	c.pkgErr = err
	return err
}

func (c *BackendConn) handleOKPacket(data []byte) (*Result, error) {
	var n int
	var pos int = 1

	r := new(Result)

	r.AffectedRows, _, n = util.LengthEncodedInt(data[pos:])
	pos += n
	r.InsertId, _, n = util.LengthEncodedInt(data[pos:])
	pos += n

	if c.capability&config.CLIENT_PROTOCOL_41 > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2

		//todo:strict_mode, check warnings as error
		//Warnings := binary.LittleEndian.Uint16(data[pos:])
		//pos += 2
	} else if c.capability&config.CLIENT_TRANSACTIONS > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		c.status = r.Status
		pos += 2
	}

	//info
	return r, nil

}

func (c *BackendConn) handleErrorPacket(data []byte) error {

	e := new(SqlError)

	var pos int = 1

	e.Code = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if c.capability&config.CLIENT_PROTOCOL_41 > 0 {
		//skip '#'
		pos++
		e.State = string(data[pos : pos+5])
		pos += 5
	}

	e.Message = string(data[pos:])

	return nil
}
