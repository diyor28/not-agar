package bytesIO

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/csbin/bitmask"
	"math"
)

func NewReader(data []byte) *BytesReader {
	reader := bytes.NewReader(data)
	return &BytesReader{reader: reader}
}

func NewWriter() *BytesWriter {
	return &BytesWriter{}
}

type BytesReader struct {
	Bytes  []byte
	reader *bytes.Reader
}

func (r *BytesReader) ReadByte() (byte, error) {
	return r.reader.ReadByte()
}

func (r BytesReader) ReadBytes(n int) ([]byte, error) {
	result := make([]byte, n)
	_, err := r.reader.Read(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r BytesReader) ReadUint(n int) (uint64, error) {
	switch n {
	case 1:
		u, err := r.ReadUint8()
		if err != nil {
			return 0, err
		}
		return uint64(u), nil
	case 2:
		u, err := r.ReadUint16()
		if err != nil {
			return 0, err
		}
		return uint64(u), nil
	case 4:
		u, err := r.ReadUint32()
		if err != nil {
			return 0, err
		}
		return uint64(u), nil
	case 8:
		return r.ReadUint64()
	}
	return 0, errors.New(fmt.Sprintf("%d is not a valid size", n))
}

func (r *BytesReader) ReadUint8() (uint8, error) {
	return r.ReadByte()
}

func (r *BytesReader) ReadUint16() (uint16, error) {
	uBytes, err := r.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(uBytes), nil
}

func (r *BytesReader) ReadUint32() (uint32, error) {
	uBytes, err := r.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(uBytes), nil
}

func (r *BytesReader) ReadUint64() (uint64, error) {
	uBytes, err := r.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(uBytes), nil
}

func (r *BytesReader) ReadInt(n int) (int64, error) {
	switch n {
	case 1:
		i, err := r.ReadInt8()
		if err != nil {
			return 0, err
		}
		return int64(i), nil
	case 2:
		i, err := r.ReadInt16()
		if err != nil {
			return 0, err
		}
		return int64(i), nil
	case 4:
		i, err := r.ReadInt32()
		if err != nil {
			return 0, err
		}
		return int64(i), nil
	case 8:
		return r.ReadInt64()
	}
	return 0, errors.New(fmt.Sprintf("%d is not a valid size", n))
}

func (r *BytesReader) ReadInt8() (int8, error) {
	v, err := r.ReadUint8()
	return int8(v), err
}

func (r *BytesReader) ReadInt16() (int16, error) {
	v, err := r.ReadUint16()
	return int16(v), err
}

func (r *BytesReader) ReadInt32() (int32, error) {
	v, err := r.ReadUint32()
	return int32(v), err
}

func (r *BytesReader) ReadInt64() (int64, error) {
	v, err := r.ReadUint64()
	return int64(v), err
}

func (r BytesReader) ReadFloat(n int) (float64, error) {
	switch n {
	case 4:
		f, err := r.ReadFloat32()
		if err != nil {
			return 0, err
		}
		return float64(f), nil
	case 8:
		return r.ReadFloat64()
	}
	return 0, errors.New(fmt.Sprintf("%d is not a valid size", n))
}

func (r *BytesReader) ReadFloat32() (float32, error) {
	fBytes, err := r.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	var res float32
	err = binary.Read(bytes.NewReader(fBytes), binary.BigEndian, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *BytesReader) ReadFloat64() (float64, error) {
	fBytes, err := r.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	var res float64
	err = binary.Read(bytes.NewReader(fBytes), binary.BigEndian, &res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *BytesReader) ReadString() (string, error) {
	length, err := r.ReadByte()
	if err != nil {
		return "", nil
	}
	sBytes, err := r.ReadBytes(int(length))
	if err != nil {
		return "", nil
	}
	return string(sBytes), nil
}

func (r *BytesReader) ReadBitmask() (*bitmask.Bitmask, error) {
	bitmaskLen, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	bBytes, err := r.ReadBytes(int(bitmaskLen))
	if err != nil {
		return nil, err
	}
	bMask := bitmask.New()
	bMask.SetBytes(bBytes)
	return bMask, nil
}

type BytesWriter struct {
	Bytes []byte
}

func (w *BytesWriter) WriteByte(b byte) {
	w.Bytes = append(w.Bytes, b)
}

func (w *BytesWriter) WriteBytes(b []byte) {
	w.Bytes = append(w.Bytes, b...)
}

func (w *BytesWriter) WriteString(s string) {
	w.WriteByte(uint8(len(s)))
	w.WriteBytes([]byte(s))
}

func (w *BytesWriter) Write(u interface{}) {
	switch t := u.(type) {
	case uint8:
		w.WriteUint8(t)
	case uint16:
		w.WriteUint16(t)
	case uint32:
		w.WriteUint32(t)
	case uint64:
		w.WriteUint64(t)
	case int8:
		w.WriteInt8(t)
	case int16:
		w.WriteInt16(t)
	case int32:
		w.WriteInt32(t)
	case int64:
		w.WriteInt64(t)
	case float32:
		w.WriteFloat32(t)
	case float64:
		w.WriteFloat64(t)
	case string:
		w.WriteString(t)
	}
}

func (w *BytesWriter) WriteUint8(u uint8) {
	w.WriteByte(u)
}
func (w *BytesWriter) WriteUint16(u uint16) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, u)
	w.WriteBytes(b)
}
func (w *BytesWriter) WriteUint32(u uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u)
	w.WriteBytes(b)
}
func (w *BytesWriter) WriteUint64(u uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	w.WriteBytes(b)
}

func (w *BytesWriter) WriteInt8(i int8) {
	w.WriteUint8(uint8(i))
}

func (w *BytesWriter) WriteInt16(i int16) {
	w.WriteUint16(uint16(i))
}

func (w *BytesWriter) WriteInt32(i int32) {
	w.WriteUint32(uint32(i))
}

func (w *BytesWriter) WriteInt64(i int64) {
	w.WriteUint64(uint64(i))
}

func (w *BytesWriter) WriteFloat32(f float32) {
	w.WriteUint32(math.Float32bits(f))
}

func (w *BytesWriter) WriteFloat64(f float64) {
	w.WriteUint64(math.Float64bits(f))
}
