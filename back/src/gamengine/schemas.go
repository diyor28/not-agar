package gamengine

import (
	"github.com/diyor28/not-agar/src/csbin"
	_map "github.com/diyor28/not-agar/src/gamengine/map"
	"github.com/diyor28/not-agar/src/gamengine/map/food"
	"github.com/diyor28/not-agar/src/gamengine/map/players"
	"github.com/diyor28/not-agar/src/gamengine/map/spikes"
	"reflect"
)

var playerField *csbin.Field = csbin.NewField("players", reflect.Struct, nil, csbin.Fields{
	csbin.NewField("x", reflect.Float32, nil, nil),
	csbin.NewField("y", reflect.Float32, nil, nil),
	csbin.NewField("weight", reflect.Float32, nil, nil),
	csbin.NewField("color", reflect.Slice, csbin.NewField("color", reflect.Uint8, nil, nil), nil),
})
var spikeField *csbin.Field = csbin.NewField("spike", reflect.Struct, nil, csbin.Fields{
	csbin.NewField("x", reflect.Float32, nil, nil),
	csbin.NewField("y", reflect.Float32, nil, nil),
	csbin.NewField("weight", reflect.Float32, nil, nil),
})
var foodField *csbin.Field = csbin.NewField("food", reflect.Struct, nil, csbin.Fields{
	csbin.NewField("x", reflect.Float32, nil, nil),
	csbin.NewField("y", reflect.Float32, nil, nil),
	csbin.NewField("weight", reflect.Float32, nil, nil),
	csbin.NewField("color", reflect.Slice, csbin.NewField("color", reflect.Uint8, nil, nil), nil),
})

var GenericSchema *csbin.Schema = csbin.New(csbin.Fields{
	csbin.NewField("event", reflect.String, nil, nil),
})

type GenericEvent struct {
	Event string
}

var StartSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	csbin.NewField("nickname", reflect.String, nil, nil),
})

var StartedSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	playerField,
	csbin.NewField("spikes", reflect.Slice, spikeField, nil),
})

type StartedEvent struct {
	Event  string
	Player players.SelfPlayer
	Spikes spikes.Spikes
}

var MoveSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
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

var PlayerStatsSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	csbin.NewField("topPlayers", reflect.Slice, csbin.NewField("players", reflect.Struct, nil, csbin.Fields{
		csbin.NewField("nickname", reflect.String, nil, nil),
		csbin.NewField("weight", reflect.Int16, nil, nil),
	}), nil),
})

type PlayerStatsEvent struct {
	Event      string
	TopPlayers []_map.PlayerStat
}

var AdminStatsSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	csbin.NewField("botsCount", reflect.Uint16, nil, nil),
	csbin.NewField("playersCount", reflect.Uint16, nil, nil),
	csbin.NewField("topsPlayers", reflect.Slice, playerField, nil),
})

var FoodUpdateSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	csbin.NewField("foods", reflect.Slice, foodField, nil),
})

type FoodUpdateEvent struct {
	Event string
	Food  []food.Food
}

var PlayersUpdatedSchema *csbin.Schema = GenericSchema.Extends(csbin.Fields{
	csbin.NewField("players", reflect.Slice, playerField, nil),
})

type PlayersUpdatedEvent struct {
	Event   string
	Players players.Players
}
