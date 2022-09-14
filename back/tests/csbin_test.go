package tests

import (
	"bytes"
	"github.com/diyor28/not-agar/src/csbin"
	"github.com/diyor28/not-agar/src/csbin/bitmask"
	"reflect"
	"testing"
)

type MoveEvent struct {
	Event string
}

func TestBitmapToBytes(t *testing.T) {
	bmask := bitmask.New()
	bmask.Set(0, true)
	bmask.Set(1, false)
	bmask.Set(2, true)
	bs := bmask.ToBytes()
	if !bytes.Equal(bs, []byte{1, 5}) {
		t.Error("expected: ", []byte{1, 5}, "got: ", bs)
	}
}

func TestBitmapHas(t *testing.T) {
	bmask := bitmask.New()
	bmask.SetBytes([]byte{5})
	if !bmask.Has(0) {
		t.Error("expected bmask[0] to be true")
	}
	if bmask.Has(1) {
		t.Error("expected bmask[1] to be false")
	}
	if !bmask.Has(2) {
		t.Error("expected bmask[2] to be true")
	}
}

func TestCodecEncodeDecode(t *testing.T) {
	fields := csbin.Fields{
		csbin.NewField("event", reflect.String, nil, nil),
	}
	codec := csbin.New(fields)
	encoded, err := codec.Encode(&MoveEvent{Event: "update"})
	if err != nil {
		t.Error(err)
	}
	expected := append([]byte{1, 1, 6}, []byte("update")...)
	if !bytes.Equal(encoded, expected) {
		t.Error("expected: ", expected, "got: ", encoded)
		return
	}
	decoded := &MoveEvent{}
	if err := codec.Decode(encoded, decoded); err != nil {
		t.Error(err)
		return
	}
	if decoded.Event != "update" {
		t.Error("expected: \"update\" got: ", decoded.Event)
	}
}

func TestCodecEncodeDecodeMap(t *testing.T) {
	fields := csbin.Fields{
		csbin.NewField("event", reflect.String, nil, nil),
	}
	codec := csbin.New(fields)
	encoded, err := codec.Encode(&map[string]interface{}{"event": "update"})
	if err != nil {
		t.Error(err)
	}
	expected := append([]byte{1, 1, 6}, []byte("update")...)
	if !bytes.Equal(encoded, expected) {
		t.Error("expected: ", expected, "got: ", encoded)
	}
	decoded := make(map[string]interface{})
	if err := codec.Decode(encoded, &decoded); err != nil {
		t.Error(err)
		return
	}
	if decoded["event"] != "update" {
		t.Error("expected: \"update\" got: ", decoded["event"])
	}
}
