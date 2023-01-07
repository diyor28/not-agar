package schemas

import (
	"github.com/diyor28/not-agar/cmd/gamengine/map/entity"
	"github.com/diyor28/not-agar/pkg/constants"
)

type Color [3]uint8

type Point struct {
	X int16
	Y int16
}

type Spike struct {
	X      float32
	Y      float32
	Weight float32
}

type Player struct {
	X        uint16
	Y        uint16
	Weight   float32
	Nickname string
	Color    Color
}

type Food struct {
	Id     entity.Id
	X      float32
	Y      float32
	Weight float32
	Color  Color
}

type GenericEvent struct {
	Event constants.GameEvent
}

type MovedEvent struct {
	Event     constants.GameEvent
	X         float32
	Y         float32
	VelocityX float32
	VelocityY float32
	Weight    float32
	Zoom      float32
	Points    []*Point
}

type StartedEventPlayer struct {
	X      float32
	Y      float32
	Weight float32
	Color  Color
	Points []*Point
}

type StartedEvent struct {
	Event  constants.GameEvent
	Player *StartedEventPlayer
	Spikes []*Spike
	Food   []*Food
}

type MoveEvent struct {
	Event constants.GameEvent
	NewX  float32
	NewY  float32
}

type PlayerStat struct {
	Nickname string
	Weight   int16
}

type PlayerStatsEvent struct {
	Event      constants.GameEvent
	TopPlayers []*PlayerStat
}

type FoodEatenEvent struct {
	Event constants.GameEvent
	Id    entity.Id
}

type PlayersUpdatedEvent struct {
	Event   constants.GameEvent
	Players []*Player
}

type FoodCreatedEvent struct {
	Event constants.GameEvent
	Food  []*Food
}
