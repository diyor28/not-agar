package schemas

import (
	"github.com/diyor28/not-agar/src/csbin"
	"reflect"
)

var playerField = csbin.NewField("player", reflect.Struct).SubFields(
	csbin.NewField("x", reflect.Uint16),
	csbin.NewField("y", reflect.Uint16),
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
			csbin.NewField("x", reflect.Int16),
			csbin.NewField("y", reflect.Int16),
		)),
	),
	csbin.NewField("spikes", reflect.Slice).MaxLen(255).SubType(spikeField),
	csbin.NewField("food", reflect.Slice).MaxLen(10000).SubType(foodField),
)

var MoveSchema = GenericSchema.Extends(
	csbin.NewField("newX", reflect.Float32),
	csbin.NewField("newY", reflect.Float32),
)

var MovedSchema = GenericSchema.Extends(
	csbin.NewField("x", reflect.Float32),
	csbin.NewField("y", reflect.Float32),
	csbin.NewField("weight", reflect.Float32),
	csbin.NewField("velocityX", reflect.Float32),
	csbin.NewField("velocityY", reflect.Float32),
	csbin.NewField("zoom", reflect.Float32),
	csbin.NewField("points", reflect.Slice).MaxLen(255).SubType(
		csbin.NewField("point", reflect.Struct).SubFields(
			csbin.NewField("x", reflect.Int16),
			csbin.NewField("y", reflect.Int16),
		)),
)

var PlayerStatsSchema = GenericSchema.Extends(
	csbin.NewField("topPlayers", reflect.Slice).MaxLen(255).SubType(csbin.NewField("players", reflect.Struct).SubFields(
		csbin.NewField("nickname", reflect.String).MaxLen(255),
		csbin.NewField("weight", reflect.Int16),
	)),
)

var AdminStatsSchema = GenericSchema.Extends(
	csbin.NewField("botsCount", reflect.Uint16),
	csbin.NewField("playersCount", reflect.Uint16),
	csbin.NewField("topsPlayers", reflect.Slice).MaxLen(255).SubType(playerField),
)

var FoodCreatedSchema = GenericSchema.Extends(
	csbin.NewField("food", reflect.Slice).MaxLen(10000).SubType(foodField),
)

var FoodEatenSchema = GenericSchema.Extends(
	csbin.NewField("id", reflect.Uint32),
)

var PlayersUpdatedSchema = GenericSchema.Extends(
	csbin.NewField("players", reflect.Slice).MaxLen(255).SubType(playerField),
)
