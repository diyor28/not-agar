package constants

const (
	MaxXY                      = 10000 // meters
	SurfaceArea                = MaxXY * MaxXY
	FoodWeight                 = 20    // kg
	SpikeWeight                = 120   // kg
	Density                    = 1.225 // kg/m3
	DragCoefficient    float32 = 2.5
	Force              float32 = 12000 // N
	MaxNumFood                 = SurfaceArea / 50000
	MinWeight                  = 30             // kg
	MaxWeight                  = MinWeight * 25 // kg
	MinZoom                    = 0.7
	MaxZoom                    = 1.5
	MaxPlayers                 = 20
	MaxSpikes                  = MaxPlayers * 2
	StatsNumber                = 10
	NumFoodResponse            = 30
	NumPlayersResponse         = 10
	MaxBots                    = 20
	SpikesSpacing              = MaxXY/MaxSpikes + SpikeWeight
)

type GameEvent uint8

const (
	Ping GameEvent = iota
	Pong
	Move
	Moved
	Start
	Started
	FoodEaten
	FoodCreated
	PlayersUpdate
	StatsUpdate
	Rip
)
