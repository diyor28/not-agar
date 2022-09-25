package bytesIO

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/csbin/bitmask"
	"math"
	"strings"
)

func NewWriter() *BytesWriter {
	return &BytesWriter{}
}

type bytesExplanation struct {
	bytes       int
	explanation string
}

type BytesWriter struct {
	bytes        []byte
	explanations []*bytesExplanation
}

func (w *BytesWriter) Compress() error {
	if len(w.bytes) == 0 {
		return errors.New("nothing to compress")
	}
	var b bytes.Buffer
	flt, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return err
	}
	if _, err := flt.Write(w.bytes); err != nil {
		return err
	}
	if err := flt.Close(); err != nil {
		return err
	}
	w.bytes = b.Bytes()
	return nil
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
		w.WriteUint(sLen, bitmask.MinBytes(maxLen), explanation)
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

func (w *BytesWriter) WriteUint(u uint64, size int, explanation string) {
	switch size {
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
