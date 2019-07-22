package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"net"
)

var (
	ErrOutOfBounds = errors.New("Read out of bounds")
)

type PacketBuilder struct {
	bytes.Buffer
}

func (b *PacketBuilder) WriteCString(s string) {
	b.WriteString(s)
	b.WriteByte(0)
}

func (b *PacketBuilder) WriteBytes(bytes []byte) {
	b.Write(bytes)
}

type PacketReader struct {
	buffer []byte
	pos    int
}

func NewPacketReader(b []byte) *PacketReader {
	return &PacketReader{
		buffer: b,
		pos:    0,
	}
}

func (r *PacketReader) CanRead(size int) error {
	if r.pos+size > len(r.buffer) {
		return ErrOutOfBounds
	}

	return nil
}

func (r *PacketReader) Pos() int {
	return r.pos
}

func (r *PacketReader) SetPos(newPos int) {
	r.pos = newPos
}

func (r *PacketReader) ReadIPv4() (net.IP, error) {
	if err := r.CanRead(net.IPv4len); err != nil {
		return nil, err
	}

	ip := net.IP(r.buffer[r.pos : r.pos+net.IPv4len])

	r.pos += net.IPv4len

	return ip, nil
}

func (r *PacketReader) ReadPort() (uint16, error) {
	if err := r.CanRead(2); err != nil {
		return 0, err
	}

	port := binary.BigEndian.Uint16(r.buffer[r.pos:])

	r.pos += 2

	return port, nil
}

func (r *PacketReader) ReadUint8() uint8 {
	b := r.buffer[r.pos]
	r.pos++
	return b
}

func (r *PacketReader) ReadUint16() uint16 {
	u16 := binary.LittleEndian.Uint16(r.buffer[r.pos:])
	r.pos += 2
	return u16
}

func (r *PacketReader) ReadUint32() uint32 {
	u32 := binary.LittleEndian.Uint32(r.buffer[r.pos:])
	r.pos += 4
	return u32
}

func (r *PacketReader) ReadInt32() int32 {
	return int32(r.ReadUint32())
}

func (r *PacketReader) ReadUint64() uint64 {
	u64 := binary.LittleEndian.Uint64(r.buffer[r.pos:])
	r.pos += 8
	return u64
}

func (r *PacketReader) ReadFloat32() float32 {
	bits := r.ReadUint32()

	return math.Float32frombits(bits)
}

func (r *PacketReader) TryReadString() (string, bool) {
	start := r.pos
	for r.pos < len(r.buffer) {
		if r.buffer[r.pos] == 0 {
			r.pos++
			return string(r.buffer[start : r.pos-1]), true
		}
		r.pos++
	}
	return "", false
}

func (r *PacketReader) ReadString() string {
	start := r.pos
	for {
		if r.buffer[r.pos] == 0 {
			r.pos++
			break
		}
		r.pos++
	}
	return string(r.buffer[start : r.pos-1])
}

func (r *PacketReader) More() bool {
	return r.pos < len(r.buffer)
}
