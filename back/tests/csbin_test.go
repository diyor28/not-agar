package tests

import (
	"bytes"
	"encoding/hex"
	"fmt"
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
	bmask.Set(true)
	bmask.Set(true)
	bmask.Set(false)
	bs := bmask.ToBytes()
	if !bytes.Equal(bs, []byte{1, 6}) {
		t.Error("expected: ", []byte{1, 6}, "got: ", bs)
	}
}

func TestBitmapHas(t *testing.T) {
	bmask := bitmask.New()
	bmask.SetBytes([]byte{6})
	if !bmask.Has(0, 3) {
		t.Error("expected bmask[0] to be true")
	}
	if !bmask.Has(1, 3) {
		t.Error("expected bmask[1] to be true")
	}
	if bmask.Has(2, 3) {
		t.Error("expected bmask[2] to be false")
	}
}

func TestCodecEncodeDecode(t *testing.T) {
	codec := csbin.New(csbin.NewField("event", reflect.String))
	data, err := codec.Encode(&MoveEvent{Event: "update"})
	encoded := data.Bytes()
	if err != nil {
		t.Error(err)
	}
	expected := append([]byte{1, 1, 0, 6}, []byte("update")...)
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
		t.Error(fmt.Sprintf("expected: {\"event\":\"update\"} got: {\"event\": \"%s\"}", decoded.Event))
	}
}

func TestCodecEncodeDecodeMap(t *testing.T) {
	codec := csbin.New(csbin.NewField("event", reflect.String))
	data, err := codec.Encode(&map[string]interface{}{"event": "update"})
	encoded := data.Bytes()
	if err != nil {
		t.Error(err)
		return
	}
	expected := append([]byte{1, 1, 0, 6}, []byte("update")...)
	if !bytes.Equal(encoded, expected) {
		t.Error("expected: ", expected, "got: ", encoded)
	}
	decoded := make(map[string]interface{})
	if err := codec.Decode(encoded, &decoded); err != nil {
		t.Error(err)
		return
	}
	if decoded["event"] != "update" {
		t.Error(fmt.Sprintf("expected: {\"event\":\"update\"} got: {\"event\": \"%s\"}", decoded["event"]))
	}
}

func TestCodecDecodeMapFromHex(t *testing.T) {
	codec := csbin.New(
		csbin.NewField("event", reflect.String),
		csbin.NewField("nickname", reflect.String),
	)
	encoded, err := hex.DecodeString("010300057374617274000464656d6f")
	if err != nil {
		t.Error(err)
		return
	}
	decoded := make(map[string]interface{})
	if err := codec.Decode(encoded, &decoded); err != nil {
		t.Error(err)
		return
	}
	if decoded["event"] != "start" {
		t.Error(fmt.Sprintf("expected: {\"event\":\"start\"} got: {\"event\": \"%s\"}", decoded["event"]))
	}
	if decoded["nickname"] != "demo" {
		t.Error(fmt.Sprintf("expected: {\"nickname\":\"demo\"} got: {\"nickname\": \"%s\"}", decoded["nickname"]))
	}
}

func TestCodecUint64Decode(t *testing.T) {
	var timestamp uint64 = 16000000000000000000
	codec := csbin.New(
		csbin.NewField("event", reflect.String),
		csbin.NewField("timestamp", reflect.Uint64),
	)
	encoded, err := hex.DecodeString("0103000470696e67de0b6b3a76400000")
	if err != nil {
		t.Error(err)
		return
	}
	decoded := make(map[string]interface{})
	if err := codec.Decode(encoded, &decoded); err != nil {
		t.Error(err)
		return
	}
	if decoded["event"] != "ping" {
		t.Error(fmt.Sprintf("expected: {\"event\":\"ping\"} got: {\"event\": \"%s\"}", decoded["event"]))
	}
	if decoded["timestamp"] != timestamp {
		t.Error(fmt.Sprintf("expected: {\"nickname\":\"demo\"} got: {\"nickname\": \"%s\"}", decoded["nickname"]))
	}
}

func TestCodecMixedTypesDecode(t *testing.T) {
	type Player struct {
		X uint8
		Y uint8
		Z uint8
	}
	type Stat struct {
		PlayerId string
		Score    int32
	}

	codec := csbin.New(
		csbin.NewField("event", reflect.String),
		csbin.NewField("player", reflect.Struct).UseStruct(&Player{}).SubFields(
			csbin.NewField("x", reflect.Uint8),
			csbin.NewField("y", reflect.Uint8),
			csbin.NewField("z", reflect.Uint8),
		),
		csbin.NewField("stats", reflect.Slice).SubType(
			csbin.NewField("stat", reflect.Struct).UseStruct(&Stat{}).SubFields(
				csbin.NewField("playerId", reflect.String),
				csbin.NewField("score", reflect.Int32),
			)),
	)
	encoded, err := hex.DecodeString("0107000675706461746501060a0c0001010300073433323134326300000020")
	if err != nil {
		t.Error(err)
		return
	}
	decoded := make(map[string]interface{})
	if err := codec.Decode(encoded, &decoded); err != nil {
		t.Error(err)
		return
	}
	player := decoded["player"].(Player)
	stat := decoded["stats"].([]Stat)[0]
	if decoded["event"] != "update" {
		t.Error(fmt.Sprintf("expected: {\"event\":\"update\"} got: {\"event\": %s}", decoded["event"]))
	}
	if player.X != 10 {
		t.Error(fmt.Sprintf("expected: {\"player.X\":10} got: {\"player.X\": %d}", player.X))
	}
	if player.Y != 12 {
		t.Error(fmt.Sprintf("expected: {\"player.Y\":12} got: {\"player.Y\": %d}", player.Y))
	}
	if player.Z != 0 {
		t.Error(fmt.Sprintf("expected: {\"player.Z\":0} got: {\"player.Z\": %d}", player.Z))
	}
	if stat.PlayerId != "432142c" {
		t.Error(fmt.Sprintf("expected: {\"stats[0].playerId\":\"432142c\"} got: {\"stats[0].playerId\": \"%s\"}", stat.PlayerId))
	}
	if stat.Score != 32 {
		t.Error(fmt.Sprintf("expected: {\"stats[0].score\":32} got: {\"stats[0].score\": %d}", stat.Score))
	}
}

func TestCodecArrayTypeEncode(t *testing.T) {
	codec := csbin.New(
		csbin.NewField("colors", reflect.Array).SubType(
			csbin.NewField("color", reflect.Uint8),
		),
	)

	data := map[string]interface{}{
		"colors": [3]uint8{127, 50, 105},
	}
	writer, err := codec.Encode(&data)
	if err != nil {
		t.Error(err)
		return
	}
	encoded := writer.Bytes()
	hexString := hex.EncodeToString(encoded)
	if hexString != "010100037f3269" {
		t.Error(fmt.Sprintf("expected: 010100037f3269 \ngot: %s", hexString))
	}
}

func TestCodecMixedTypesEncode(t *testing.T) {
	type Player struct {
		X uint8
		Y uint8
		Z uint8
	}
	type Stat struct {
		PlayerId string
		Score    int32
	}

	codec := csbin.New(
		csbin.NewField("event", reflect.String),
		csbin.NewField("player", reflect.Struct).UseStruct(&Player{}).SubFields(
			csbin.NewField("x", reflect.Uint8),
			csbin.NewField("y", reflect.Uint8),
			csbin.NewField("z", reflect.Uint8),
		),
		csbin.NewField("stats", reflect.Slice).SubType(
			csbin.NewField("stat", reflect.Struct).UseStruct(&Stat{}).SubFields(
				csbin.NewField("playerId", reflect.String),
				csbin.NewField("score", reflect.Int32),
			)),
	)
	data := map[string]interface{}{
		"event": "update",
		"player": &Player{
			X: 10,
			Y: 12,
			Z: 0,
		},
		"stats": []*Stat{{PlayerId: "432142c", Score: 32}},
	}
	writer, err := codec.Encode(&data)
	encoded := writer.Bytes()
	if err != nil {
		t.Error(err)
		return
	}
	hexString := hex.EncodeToString(encoded)
	if hexString != "0107000675706461746501060a0c0001010300073433323134326300000020" {
		t.Error(fmt.Sprintf("expected: 0107000675706461746501060a0c0001010300073433323134326300000020 \ngot: %s", hexString))
	}
}
