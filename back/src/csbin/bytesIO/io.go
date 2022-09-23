package bytesIO

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/csbin/bitmask"
	"math"
	"strings"
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

type bytesExplanation struct {
	bytes       int
	explanation string
}

type BytesWriter struct {
	bytes        []byte
	explanations []*bytesExplanation
}

func (w *BytesWriter) Bytes() []byte {
	return w.bytes
}

func (w *BytesWriter) formattedHexString() string {
	index := 0
	hexString := hex.EncodeToString(w.bytes)
	for _, explanation := range w.explanations {
		index += explanation.bytes * 2
		hexString = hexString[0:index] + "|" + hexString[index:]
		index++
	}
	return hexString
}

func (w *BytesWriter) maxExpLength() int {
	maxLen := 0
	for _, explanation := range w.explanations {
		if len(explanation.explanation) > maxLen {
			maxLen = len(explanation.explanation)
		}
	}
	return maxLen
}

func (w *BytesWriter) Explain() string {
	hexString := w.formattedHexString()
	maxExpLength := w.maxExpLength()
	explanation := "\n"
	for i := 0; i < maxExpLength; i++ {
		for _, exp := range w.explanations {
			hexBytes := exp.bytes * 2
			if len(exp.explanation) > i {
				explanation += string(exp.explanation[i]) + strings.Repeat(" ", hexBytes-1) + "|"
			} else {
				explanation += strings.Repeat(" ", hexBytes) + "|"
			}
		}
		explanation += "\n"
	}
	return hexString + explanation
}

func (w *BytesWriter) WriteByte(b byte, explanation string) {
	w.bytes = append(w.bytes, b)
	w.explanations = append(w.explanations, &bytesExplanation{1, explanation})
}

func (w *BytesWriter) WriteBytes(b []byte, explanation string) {
	w.bytes = append(w.bytes, b...)
	w.explanations = append(w.explanations, &bytesExplanation{len(b), explanation})
}

func (w *BytesWriter) WriteString(s string, explanation string, length uint64, maxLen uint64) error {
	sLen := uint64(len(s))
	if length > 0 {
		if sLen != length {
			return errors.New(fmt.Sprintf("expected a string of length %d, got %d", length, len(s)))
		}
		w.WriteBytes([]byte(s), explanation)
		return nil
	}
	if maxLen > 0 {
		if sLen > maxLen {
			return errors.New(fmt.Sprintf("expected a string of length <= %d, got %d", maxLen, len(s)))
		}
		w.WriteUint(sLen, explanation)
	} else {
		w.WriteUint16(uint16(sLen), "string length")
	}
	w.WriteBytes([]byte(s), explanation)
	return nil
}

func (w *BytesWriter) WriteNumeric(u interface{}, explanation string) {
	switch t := u.(type) {
	case uint8:
		w.WriteUint8(t, explanation)
	case uint16:
		w.WriteUint16(t, explanation)
	case uint32:
		w.WriteUint32(t, explanation)
	case uint64:
		w.WriteUint64(t, explanation)
	case int8:
		w.WriteInt8(t, explanation)
	case int16:
		w.WriteInt16(t, explanation)
	case int32:
		w.WriteInt32(t, explanation)
	case int64:
		w.WriteInt64(t, explanation)
	case float32:
		w.WriteFloat32(t, explanation)
	case float64:
		w.WriteFloat64(t, explanation)
	}
}

func (w *BytesWriter) WriteBool(b bool, explanation string) {
	if b {
		w.WriteByte(1, explanation)
	} else {
		w.WriteByte(0, explanation)
	}
}

func (w *BytesWriter) WriteUint(u uint64, explanation string) {
	bytesLen := bitmask.MinBytes(u)
	switch bytesLen {
	case 1:
		w.WriteUint8(uint8(u), explanation)
	case 2:
		w.WriteUint16(uint16(u), explanation)
	case 4:
		w.WriteUint32(uint32(u), explanation)
	default:
		w.WriteUint64(u, explanation)
	}
}

func (w *BytesWriter) WriteUint8(u uint8, explanation string) {
	w.WriteByte(u, explanation)
}
func (w *BytesWriter) WriteUint16(u uint16, explanation string) {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, u)
	w.WriteBytes(b, explanation)
}
func (w *BytesWriter) WriteUint32(u uint32, explanation string) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u)
	w.WriteBytes(b, explanation)
}
func (w *BytesWriter) WriteUint64(u uint64, explanation string) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	w.WriteBytes(b, explanation)
}

func (w *BytesWriter) WriteInt8(i int8, explanation string) {
	w.WriteUint8(uint8(i), explanation)
}

func (w *BytesWriter) WriteInt16(i int16, explanation string) {
	w.WriteUint16(uint16(i), explanation)
}

func (w *BytesWriter) WriteInt32(i int32, explanation string) {
	w.WriteUint32(uint32(i), explanation)
}

func (w *BytesWriter) WriteInt64(i int64, explanation string) {
	w.WriteUint64(uint64(i), explanation)
}

func (w *BytesWriter) WriteFloat32(f float32, explanation string) {
	w.WriteUint32(math.Float32bits(f), explanation)
}

func (w *BytesWriter) WriteFloat64(f float64, explanation string) {
	w.WriteUint64(math.Float64bits(f), explanation)
}
