package constants

const (
	MinXY              = 0
	MaxXY              = 5000
	SurfaceArea        = MaxXY * MaxXY
	FoodWeight         = 20
	SpikeWeight        = 120
	MinSpeed           = 2
	MaxSpeed           = 5
	MaxNumFood         = SurfaceArea / 50000
	MinWeight          = 40
	MaxWeight          = MinWeight * 25
	MinZoom            = 0.7
	MaxZoom            = 1.0
	MaxPlayers         = 20
	MaxSpikes          = MaxPlayers * 2
	StatsNumber        = 10
	SpeedWeightLimit   = 500
	NumFoodResponse    = 30
	NumPlayersResponse = 10
	SpikesSpacing      = MaxXY/MaxSpikes + SpikeWeight
)
