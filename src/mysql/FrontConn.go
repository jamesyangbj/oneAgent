package mysql

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	config "oneProxy/src/config"
	util "oneProxy/src/util"
	"reflect"
	"time"
	"unsafe"
)

var DEFAULT_CAPABILITY uint32 = config.CLIENT_LONG_PASSWORD | config.CLIENT_LONG_FLAG |
	config.CLIENT_CONNECT_WITH_DB | config.CLIENT_PROTOCOL_41 |
	config.CLIENT_TRANSACTIONS | config.CLIENT_SECURE_CONNECTION

type FrontConn struct {
	conn        net.Conn
	connectTime time.Time
	pkg         *PacketIO

	capability   uint32
	connectionId uint32
	salt         []byte
	status       uint16

	user string
	db   string
}

func NewConn(c net.Conn) *FrontConn {
	f := new(FrontConn)
	f.conn = c
	f.connectTime = time.Now()
	f.pkg = NewPacketIO(c)
	return f
}

func (c *FrontConn) Close() {
	c.conn.Close()

}

func (c *FrontConn) Handshake() error {
	if err := c.writeInitialHandshake(); err != nil {
		return err
	}

	if err := c.readHandshakeResponse(); err != nil {
		return err
	}

	if err := c.writeOK(nil); err != nil {
		return err
	}
	return nil
}

func (c *FrontConn) writeInitialHandshake() error {
	data := make([]byte, 4, 128)

	//min version 10
	data = append(data, 10)

	//server version[00]
	data = append(data, config.ServerVersion...)
	data = append(data, 0)

	//connection id
	data = append(data, byte(c.connectionId), byte(c.connectionId>>8), byte(c.connectionId>>16), byte(c.connectionId>>24))

	//auth-plugin-data-part-1
	data = append(data, c.salt[0:8]...)

	//filter [00]
	data = append(data, 0)

	//capability flag lower 2 bytes, using default capability here
	data = append(data, byte(DEFAULT_CAPABILITY), byte(DEFAULT_CAPABILITY>>8))

	//charset, utf-8 default
	data = append(data, uint8(config.DEFAULT_COLLATION_ID))

	//status
	data = append(data, byte(c.status), byte(c.status>>8))

	//below 13 byte may not be used
	//capability flag upper 2 bytes, using default capability here
	data = append(data, byte(DEFAULT_CAPABILITY>>16), byte(DEFAULT_CAPABILITY>>24))

	//filter [0x15], for wireshark dump, value is 0x15
	data = append(data, 0x15)

	//reserved 10 [00]
	data = append(data, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	//auth-plugin-data-part-2
	data = append(data, c.salt[8:]...)

	//filter [00]
	data = append(data, 0)

	return c.writePacket(data)
}

func (c *FrontConn) readHandshakeResponse() error {
	data, err := c.pkg.ReadPacket()
	if err != nil {
		return err
	}

	pos := 0

	c.capability = binary.LittleEndian.Uint32(data[:4])
	pos += 4

	//skip max packet size
	pos += 4

	//charset, skip, if you want to use another charset, use set names
	//c.collation = CollationId(data[pos])
	pos++

	//skip reserved 23[00]
	pos += 23

	//user name
	c.user = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])

	pos += len(c.user) + 1

	//auth length and auth
	authLen := int(data[pos])
	pos++
	auth := data[pos : pos+authLen]

	pos += authLen

	checkAuth := []byte{}

	if !bytes.Equal(auth, checkAuth) {
		fmt.Println("password error...")
	}

	var db string
	if c.capability&config.CLIENT_CONNECT_WITH_DB > 0 {
		if len(data[pos:]) == 0 {
			return nil
		}

		db = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
		pos += len(c.db) + 1

	}
	c.db = db

	return nil
}

func String(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

func (c *FrontConn) HandleData() error {
	data, err := c.pkg.ReadPacket()
	if err != nil {
		return err
	}
	cmd := data[0]
	data = data[1:]
	switch cmd {
	case config.COM_QUIT:
		c.handleRollback()
		c.Close()
		return nil
	case config.COM_QUERY:
		return c.handleQuery(String(data))
	case config.COM_PING:
		return c.writeOK(nil)
	case config.COM_INIT_DB:
		return c.handleUseDB(String(data))
	case config.COM_FIELD_LIST:
		return c.handleFieldList(data)
	case config.COM_STMT_PREPARE:
		return c.handleStmtPrepare(String(data))
	case config.COM_STMT_EXECUTE:
		return c.handleStmtExecute(data)
	case config.COM_STMT_CLOSE:
		return c.handleStmtClose(data)
	case config.COM_STMT_SEND_LONG_DATA:
		return c.handleStmtSendLongData(data)
	case config.COM_STMT_RESET:
		return c.handleStmtReset(data)
	case config.COM_SET_OPTION:
		return c.writeEOF(0)
	default:
		msg := fmt.Sprintf("command %d not supported now", cmd)
		fmt.Println(msg)

	}

	return nil
}

func (c *FrontConn) handleRollback() {

}

func (c *FrontConn) handleQuery(sql string) error {
	fmt.Println(sql)
	return nil

}

func (c *FrontConn) handleUseDB(sql string) error {
	fmt.Println(sql)
	return nil

}

func (c *FrontConn) handleFieldList(data []byte) error {
	return nil

}

func (c *FrontConn) handleStmtPrepare(sql string) error {
	fmt.Println(sql)
	return nil

}

func (c *FrontConn) handleStmtExecute(data []byte) error {
	return nil

}

func (c *FrontConn) handleStmtClose(data []byte) error {
	return nil

}

func (c *FrontConn) handleStmtSendLongData(data []byte) error {
	return nil

}

func (c *FrontConn) handleStmtReset(data []byte) error {
	return nil

}

func (c *FrontConn) writeEOF(i int) error {
	return nil
}

func (c *FrontConn) writeOK(r *Result) error {
	if r == nil {
		r = &Result{Status: c.status}
	}
	data := make([]byte, 4, 32)

	data = append(data, config.OK_HEADER)

	data = append(data, util.PutLengthEncodedInt(r.AffectedRows)...)
	data = append(data, util.PutLengthEncodedInt(r.InsertId)...)

	if c.capability&config.CLIENT_PROTOCOL_41 > 0 {
		data = append(data, byte(r.Status), byte(r.Status>>8))
		data = append(data, 0, 0)
	}

	return c.writePacket(data)
}

func (c *FrontConn) writePacket(data []byte) error {
	return c.pkg.writeData(data)
}
