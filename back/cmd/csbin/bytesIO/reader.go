package bytesIO

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/cmd/csbin/bitmask"
)

func NewReader(data []byte) *BytesReader {
	reader := bytes.NewReader(data)
	return &BytesReader{reader: reader}
}

type BytesReader struct {
	Bytes  []byte
	reader *bytes.Reader
}

func (r *BytesReader) Decompress() error {
	if len(r.Bytes) == 0 {
		return errors.New("nothing to decompress")
	}
	b := bytes.NewBuffer(r.Bytes)
	gz, err := gzip.NewReader(b)
	if err != nil {
		return err
	}
	var uncompressed []byte
	if _, err := gz.Read(uncompressed); err != nil {
		return err
	}
	r.Bytes = uncompressed
	return nil
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

func (r *BytesReader) ReadBool() (bool, error) {
	u, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	if u == 0 {
		return false, nil
	}
	if u == 1 {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("expected 0 or 1, got: %d", u))
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

func (r *BytesReader) ReadString(len uint64, maxLen uint64) (string, error) {
	var length uint64
	if len > 0 {
		length = len
	} else if maxLen > 0 {
		if l, err := r.ReadUint(bitmask.MinBytes(maxLen)); err == nil {
			length = l
		} else {
			return "", err
		}
	} else {
		if l, err := r.ReadUint16(); err == nil {
			length = uint64(l)
		} else {
			return "", err
		}
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
