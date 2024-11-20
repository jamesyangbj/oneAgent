package mysql

import (
	"bufio"
	"io"
	"net"
	"oneProxy/src/config"
)

const (
	DEFAULT_READER_SIZE = 4 * 1024
)

type PacketIO struct {
	Sequence uint8
	rb       *bufio.Reader
	wb       io.Writer
}

func NewPacketIO(conn net.Conn) *PacketIO {
	p := new(PacketIO)
	p.Sequence = 0
	p.rb = bufio.NewReaderSize(conn, DEFAULT_READER_SIZE)
	p.wb = conn
	return p
}

func (p *PacketIO) ReadPacket() ([]byte, error) {
	header := []byte{0, 0, 0, 0}
	if _, err := io.ReadFull(p.rb, header); err != nil {
		return nil, err
	}

	length := uint32(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)
	sequence := uint8(header[3])
	if sequence != p.Sequence {

	}
	data := make([]byte, length)

	if _, err := io.ReadFull(p.rb, data); err != nil {
		return nil, err
	} else {
		if int(length) < config.MaxPayloadLen {
			return data, err
		} else {
			var buf []byte
			var err error
			if buf, err = p.ReadPacket(); err != nil {
				return nil, err
			}
			return append(data, buf...), nil
		}
	}
}

func (p *PacketIO) WritePacket(data []byte) error {
	header := []byte{0, 0, 0, 0}
	for {
		allData := []byte{}
		length := len(data)
		if int(length) <= config.MaxPayloadLen {
			break
		}

		header[0] = 0xff
		header[1] = 0xff
		header[2] = 0xff
		header[3] = p.Sequence
		allData = append(allData, header...)
		allData = append(allData, data[:config.MaxPayloadLen]...)
		data = data[config.MaxPayloadLen:]
		p.Sequence++
		p.writeData(allData)
	}

	allData := []byte{}
	length := len(data)
	header[0] = byte(length)
	header[1] = byte(length >> 8)
	header[2] = byte(length >> 16)
	header[3] = p.Sequence
	allData = append(allData, header...)
	allData = append(allData, data...)

	p.Sequence++
	p.writeData(allData)
	return nil
}

func (p *PacketIO) writeData(data []byte) error {
	if _, err := p.wb.Write(data); err != nil {
		return err
	}
	return nil
}
