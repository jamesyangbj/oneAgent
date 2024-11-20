package mysql

import "net"

type Conn interface {
	NewConn(c net.Conn) *Conn
	Close()
	ReadPacket() []byte
	WritePacket(buf []byte)
}

type Result struct {
	Status uint16

	InsertId     uint64
	AffectedRows uint64

	*Resultset
}

type RowData []byte

type Resultset struct {
	Fields     []*Field
	FieldNames map[string]int
	Values     [][]interface{}

	RowDatas []RowData
}

type FieldData []byte

type Field struct {
	Data         FieldData
	Schema       []byte
	Table        []byte
	OrgTable     []byte
	Name         []byte
	OrgName      []byte
	Charset      uint16
	ColumnLength uint32
	Type         uint8
	Flag         uint16
	Decimal      uint8

	DefaultValueLength uint64
	DefaultValue       []byte
}

type SqlError struct {
	Code    uint16
	Message string
	State   string
}
