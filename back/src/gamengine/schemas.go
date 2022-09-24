package gamengine

import (
	"github.com/diyor28/not-agar/src/csbin"
	"github.com/diyor28/not-agar/src/gamengine/constants"
	_map "github.com/diyor28/not-agar/src/gamengine/map"
	"github.com/diyor28/not-agar/src/gamengine/map/entity"
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
	csbin.NewField("id", reflect.Uint32),
	csbin.NewField("x", reflect.Float32),
	csbin.NewField("y", reflect.Float32),
	csbin.NewField("weight", reflect.Float32),
	csbin.NewField("color", reflect.Array).Len(3).SubType(csbin.NewField("color", reflect.Uint8)),
)

var GenericSchema = csbin.New(
	csbin.NewField("event", reflect.Uint8),
)

type GenericEvent struct {
	Event constants.GameEvent
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
	csbin.NewField("spikes", reflect.Slice).MaxLen(255).SubType(spikeField),
	csbin.NewField("food", reflect.Slice).MaxLen(10000).SubType(foodField),
)

type StartedEventPlayer struct {
	X      float32
	Y      float32
	Weight float32
	Color  [3]uint8
	Points []*shell.Point
}

type StartedEvent struct {
	Event  constants.GameEvent
	Player *StartedEventPlayer
	Spikes []*spikes.Spike
	Food   []*food.Food
}

var MoveSchema = GenericSchema.Extends(
	csbin.NewField("newX", reflect.Float32),
	csbin.NewField("newY", reflect.Float32),
)

type MoveEvent struct {
	Event constants.GameEvent
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
	Event     constants.GameEvent
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
	Event      constants.GameEvent
	TopPlayers []_map.PlayerStat
}

var AdminStatsSchema = GenericSchema.Extends(
	csbin.NewField("botsCount", reflect.Uint16),
	csbin.NewField("playersCount", reflect.Uint16),
	csbin.NewField("topsPlayers", reflect.Slice).MaxLen(255).SubType(playerField),
)

var FoodCreatedSchema = GenericSchema.Extends(
	csbin.NewField("food", reflect.Slice).MaxLen(10000).SubType(foodField),
)

type FoodCreatedEvent struct {
	Event constants.GameEvent
	Food  []*food.Food
}

var FoodEatenSchema = GenericSchema.Extends(
	csbin.NewField("id", reflect.Uint32),
)

type FoodEatenEvent struct {
	Event constants.GameEvent
	Id    entity.Id
}

var PlayersUpdatedSchema = GenericSchema.Extends(
	csbin.NewField("players", reflect.Slice).MaxLen(255).SubType(playerField),
)

type PlayersUpdatedEvent struct {
	Event   constants.GameEvent
	Players []*players.Player
}
