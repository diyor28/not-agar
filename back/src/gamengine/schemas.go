package gamengine

import (
	"github.com/diyor28/not-agar/src/csbin"
	"github.com/diyor28/not-agar/src/gamengine/player"
	"reflect"
)

var eventField *csbin.Field = csbin.NewField("event", reflect.String, nil, nil)
var playerField *csbin.Field = csbin.NewField("player", reflect.Struct, nil, csbin.Fields{
	csbin.NewField("x", reflect.Float32, nil, nil),
	csbin.NewField("y", reflect.Float32, nil, nil),
	csbin.NewField("weight", reflect.Float32, nil, nil),
	csbin.NewField("weight", reflect.Float32, nil, nil),
	csbin.NewField("nickname", reflect.String, nil, nil),
	csbin.NewField("color", reflect.Slice, csbin.NewField("color", reflect.Uint8, nil, nil), nil),
})

var GenericSchema *csbin.Schema = csbin.New(csbin.Fields{eventField})

type GenericEvent struct {
	Event string
}

var MoveSchema *csbin.Schema = csbin.New(csbin.Fields{
	eventField,
	csbin.NewField("newX", reflect.Float32, nil, nil),
	csbin.NewField("newY", reflect.Float32, nil, nil),
})

type MoveEvent struct {
	Event string
	NewX  float32 `json:"newX"`
	NewY  float32 `json:"newY"`
}

var MovedSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	csbin.NewField("x", reflect.Float32, nil, nil),
	csbin.NewField("y", reflect.Float32, nil, nil),
	csbin.NewField("weight", reflect.Float32, nil, nil),
	csbin.NewField("zoom", reflect.Float32, nil, nil),
})

type MovedEvent struct {
	Event  string
	X      float32
	Y      float32
	Weight float32
	Zoom   float32
}

var StatsSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	csbin.NewField("botsCount", reflect.Uint16, nil, nil),
	csbin.NewField("playersCount", reflect.Uint16, nil, nil),
	csbin.NewField("topsPlayers", reflect.Slice, playerField, nil),
})

var PlayersUpdateSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	csbin.NewField("players", reflect.Slice, playerField, nil),
})

type PlayersUpdateEvent struct {
	Event   string
	Players player.Players
}
