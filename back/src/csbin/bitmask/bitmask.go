package bitmask

import (
	"encoding/binary"
	"math/bits"
)

func New() *Bitmask {
	return &Bitmask{}
}

type Bitmask struct {
	Bits uint64
}

func (b *Bitmask) SetBytes(bytes []byte) {
	switch len(bytes) {
	case 1:
		b.Bits = uint64(bytes[0])
	case 2:
		b.Bits = uint64(binary.BigEndian.Uint16(bytes))
	case 4:
		b.Bits = uint64(binary.BigEndian.Uint32(bytes))
	case 8:
		b.Bits = binary.BigEndian.Uint64(bytes)
	}
}

func (b *Bitmask) content() []byte {
	bitsLen := bits.Len64(b.Bits)
	if bitsLen <= 8 {
		return []byte{uint8(b.Bits)}
	}
	if bitsLen <= 16 {
		r := make([]byte, 2)
		binary.BigEndian.PutUint16(r, uint16(b.Bits))
		return r
	}
	if bitsLen <= 32 {
		r := make([]byte, 4)
		binary.BigEndian.PutUint32(r, uint32(b.Bits))
		return r
	}
	r := make([]byte, 8)
	binary.BigEndian.PutUint64(r, b.Bits)
	return r
}

func (b *Bitmask) ToBytes() []byte {
	content := b.content()
	return append([]byte{uint8(len(content))}, content...)
}

func (b *Bitmask) Has(i int) bool {
	return (b.Bits>>uint64(i))&1 == 1
}

func (b *Bitmask) Set(i int, v bool) {
	if v {
		b.Bits |= 1 << uint64(i)
	} else {
		var mask uint64 = ^(1 << uint64(i))
		b.Bits &= mask
	}
}

func (b *Bitmask) Len() int {
	return bits.Len64(b.Bits)
}
