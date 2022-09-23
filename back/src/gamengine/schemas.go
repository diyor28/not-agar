package gamengine

import (
	"github.com/diyor28/not-agar/src/csbin"
	_map "github.com/diyor28/not-agar/src/gamengine/map"
	"github.com/diyor28/not-agar/src/gamengine/map/food"
	"github.com/diyor28/not-agar/src/gamengine/map/players"
	"github.com/diyor28/not-agar/src/gamengine/map/players/shell"
	"github.com/diyor28/not-agar/src/gamengine/map/spikes"
	"reflect"
)

var playerField = csbin.NewField("player", reflect.Struct).SubFields(
	csbin.NewField("x", reflect.Float32),
	csbin.NewField("y", reflect.Float32),
	csbin.NewField("weight", reflect.Float32),
	csbin.NewField("nickname", reflect.String).MaxLen(255),
	csbin.NewField("color", reflect.Array).Len(3).SubType(csbin.NewField("color", reflect.Uint8)),
)

var spikeField = csbin.NewField("spike", reflect.Struct).SubFields(
	csbin.NewField("x", reflect.Float32),
	csbin.NewField("y", reflect.Float32),
	csbin.NewField("weight", reflect.Float32),
)

var foodField = csbin.NewField("food", reflect.Struct).SubFields(
	csbin.NewField("x", reflect.Float32),
	csbin.NewField("y", reflect.Float32),
	csbin.NewField("weight", reflect.Float32),
	csbin.NewField("color", reflect.Array).Len(3).SubType(csbin.NewField("color", reflect.Uint8)),
)

var GenericSchema = csbin.New(
	csbin.NewField("event", reflect.String),
)

type GenericEvent struct {
	Event string
}

var PingPongSchema = GenericSchema.Extends(
	csbin.NewField("timestamp", reflect.Uint64),
)

var StartSchema = GenericSchema.Extends(
	csbin.NewField("nickname", reflect.String).MaxLen(255),
)

var StartedSchema = GenericSchema.Extends(
	csbin.NewField("player", reflect.Struct).SubFields(
		csbin.NewField("x", reflect.Float32),
		csbin.NewField("y", reflect.Float32),
		csbin.NewField("weight", reflect.Float32),
		csbin.NewField("color", reflect.Array).Len(3).SubType(csbin.NewField("color", reflect.Uint8)),
		csbin.NewField("points", reflect.Slice).MaxLen(255).SubType(csbin.NewField("point", reflect.Struct).SubFields(
			csbin.NewField("x", reflect.Float32),
			csbin.NewField("y", reflect.Float32),
		)),
	),
	csbin.NewField("spikes", reflect.Slice).SubType(spikeField),
)

type StartedEventPlayer struct {
	X      float32
	Y      float32
	Weight float32
	Color  [3]uint8
	Points []*shell.Point
}

type StartedEvent struct {
	Event  string
	Player *StartedEventPlayer
	Spikes spikes.Spikes
}

var MoveSchema = GenericSchema.Extends(
	csbin.NewField("newX", reflect.Float32),
	csbin.NewField("newY", reflect.Float32),
)

type MoveEvent struct {
	Event string
	NewX  float32 `json:"newX"`
	NewY  float32 `json:"newY"`
}

var MovedSchema = GenericSchema.Extends(
	csbin.NewField("x", reflect.Float32),
	csbin.NewField("y", reflect.Float32),
	csbin.NewField("weight", reflect.Float32),
	csbin.NewField("velocityX", reflect.Float32),
	csbin.NewField("velocityY", reflect.Float32),
	csbin.NewField("zoom", reflect.Float32),
	csbin.NewField("points", reflect.Slice).MaxLen(255).SubType(
		csbin.NewField("point", reflect.Struct).SubFields(
			csbin.NewField("x", reflect.Float32),
			csbin.NewField("y", reflect.Float32),
		)),
)

type MovedEvent struct {
	Event     string
	X         float32
	Y         float32
	VelocityX float32
	VelocityY float32
	Weight    float32
	Zoom      float32
	Points    []*shell.Point
}

var PlayerStatsSchema = GenericSchema.Extends(
	csbin.NewField("topPlayers", reflect.Slice).MaxLen(255).SubType(csbin.NewField("players", reflect.Struct).SubFields(
		csbin.NewField("nickname", reflect.String).MaxLen(255),
		csbin.NewField("weight", reflect.Int16),
	)),
)

type PlayerStatsEvent struct {
	Event      string
	TopPlayers []_map.PlayerStat
}

var AdminStatsSchema = GenericSchema.Extends(
	csbin.NewField("botsCount", reflect.Uint16),
	csbin.NewField("playersCount", reflect.Uint16),
	csbin.NewField("topsPlayers", reflect.Slice).MaxLen(255).SubType(playerField),
)

var FoodUpdateSchema = GenericSchema.Extends(
	csbin.NewField("food", reflect.Slice).MaxLen(255).SubType(foodField),
)

type FoodUpdatedEvent struct {
	Event string
	Food  []*food.Food
}

var PlayersUpdatedSchema = GenericSchema.Extends(
	csbin.NewField("players", reflect.Slice).MaxLen(255).SubType(playerField),
)

type PlayersUpdatedEvent struct {
	Event   string
	Players players.Players
}
