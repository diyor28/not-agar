package players

import (
	"errors"
	"github.com/diyor28/not-agar/src/gamengine/constants"
	"github.com/diyor28/not-agar/src/gamengine/map/entity"
	"github.com/diyor28/not-agar/src/utils"
	"github.com/frankenbeanies/uuid4"
	"math"
)

func getSpeedFromWeight(weight float32) float32 {
	normalized := 1 - (weight-constants.MinWeight)/(constants.SpeedWeightLimit-constants.MinWeight)
	newSpeed := float64(normalized*(constants.MaxSpeed-constants.MinSpeed) + constants.MinSpeed)
	return float32(math.Max(newSpeed, constants.MinSpeed))
}

type SelfPlayer struct {
	Uuid     string  `json:"uuid"`
	Nickname string  `json:"nickname"`
	Color    [3]int  `json:"color"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
	Weight   float32 `json:"weight"`
	Speed    float32 `json:"speed"`
	Zoom     float32 `json:"zoom"`
	IsBot    bool    `json:"-"`
}

type Player struct {
	*entity.Entity
	Nickname     string  `json:"nickname"`
	VelocityX    float32 `json:"-"`
	VelocityY    float32 `json:"-"`
	Color        [3]int  `json:"color"`
	Accelerating bool    `json:"-"`
	Speed        float32 `json:"-"`
	Zoom         float32 `json:"-"`
	IsBot        bool    `json:"-"`
}

func NewPlayer(x float32, y float32, weight float32, nickname string, isBot bool) *Player {
	player := &Player{
		Entity: &entity.Entity{
			Uuid: uuid4.New().String(),
			X:    x,
			Y:    y,
		},
		Speed:        constants.MaxSpeed,
		Nickname:     nickname,
		IsBot:        isBot,
		Accelerating: false,
		Zoom:         1,
		Color:        utils.RandomColor(),
	}
	player.SetWeight(weight)
	return player
}

func (pl *Player) GetSelfPlayer() SelfPlayer {
	return SelfPlayer{
		Uuid:     pl.Uuid,
		Nickname: pl.Nickname,
		Color:    pl.Color,
		X:        pl.X,
		Y:        pl.Y,
		Weight:   pl.Weight,
		Speed:    pl.Speed,
		Zoom:     pl.Zoom,
		IsBot:    pl.IsBot,
	}
}

func (pl *Player) PassiveWeightLoss() {
	if pl.Weight > constants.MinWeight*3 {
		pl.SetWeight(pl.Weight * 0.99999)
	}
}

func (pl *Player) FoodEatable(food entity.Interface) bool {
	diff := utils.CalcDistance(pl.X, food.GetX(), pl.Y, food.GetY())
	return diff < pl.Weight/2
}

func (pl *Player) UpdatePosition() {
	speed := pl.Speed
	newX := pl.X + speed*pl.VelocityX
	newY := pl.Y + speed*pl.VelocityY
	newZoom := (constants.MinWeight/pl.Weight)*(constants.MaxZoom-constants.MinZoom) + constants.MinZoom
	pl.X = utils.Clip(newX, constants.MinXY, constants.MaxXY)
	pl.Y = utils.Clip(newY, constants.MinXY, constants.MaxXY)
	pl.Zoom = utils.Clip(newZoom, constants.MinZoom, constants.MaxZoom)
}

func (pl *Player) UpdateDirection(newX float32, newY float32) {
	dist := float64(utils.CalcDistance(pl.X, newX, pl.Y, newY))
	diffX := float64(newX - pl.X)
	diffY := float64(newY - pl.Y)
	velocityX := diffX / dist
	velocityY := diffY / dist
	pl.VelocityX = float32(velocityX)
	pl.VelocityY = float32(velocityY)
}

func (pl *Player) AddWeight(weight float32) {
	sign := float32(math.Copysign(1, float64(weight)))
	nWeight := math.Sqrt(float64(pl.Weight*pl.Weight + weight*weight*sign))
	pl.SetWeight(float32(nWeight))
}

func (pl *Player) EatEntity(entity interface{ GetWeight() float32 }) {
	pl.AddWeight(entity.GetWeight())
}

func (pl *Player) SetWeight(weight float32) {
	pl.Entity.SetWeight(weight)
	pl.Speed = getSpeedFromWeight(pl.Weight)
}

type Players []*Player

func (players Players) Get(uuid string) (*Player, error) {
	for _, p := range players {
		if p.Uuid == uuid {
			return p, nil
		}
	}
	return nil, errors.New("no players found")
}

func (players Players) AsValues() []Player {
	var result []Player
	for _, p := range players {
		result = append(result, *p)
	}
	return result
}

func (players Players) Exclude(uuid string) Players {
	var result Players
	for _, p := range players {
		if p.Uuid != uuid {
			result = append(result, p)
		}
	}
	return result
}

func (players Players) Largest(k int) []Player {
	playersValues := players.AsValues()
	numOfPlayers := len(playersValues)
	resultLength := k
	if numOfPlayers < k {
		resultLength = numOfPlayers
	}
	for i := 0; i < resultLength; i++ {
		var maxIdx = i
		for j := i; j < numOfPlayers; j++ {
			if playersValues[j].Weight > playersValues[maxIdx].Weight {
				maxIdx = j
			}
		}
		playersValues[i], playersValues[maxIdx] = playersValues[maxIdx], playersValues[i]
	}
	return playersValues[:resultLength]
}

func (players Players) Closest(player *Player, kClosest int) Players {
	otherPlayers := players.Exclude(player.Uuid)
	totalPlayers := len(otherPlayers)
	playersDistances := make(map[string]float32, totalPlayers)
	for _, p := range otherPlayers {
		playersDistances[p.Uuid] = utils.CalcDistance(player.X, p.X, player.Y, p.Y)
	}

	numResults := kClosest
	if kClosest > totalPlayers {
		numResults = totalPlayers
	}
	for i := 0; i < numResults; i++ {
		var minIdx = i
		for j := i + 1; j < totalPlayers; j++ {
			if playersDistances[otherPlayers[j].Uuid] < playersDistances[otherPlayers[minIdx].Uuid] {
				minIdx = j
			}
		}
		otherPlayers[i], otherPlayers[minIdx] = otherPlayers[minIdx], otherPlayers[i]
	}
	return otherPlayers[:numResults]
}

func (players Players) Update(uuid string, newX float32, newY float32) (*Player, error) {
	player, err := players.Get(uuid)
	if err != nil {
		return nil, err
	}
	player.UpdateDirection(newX, newY)
	return player, nil
}

func (players Players) Real() Players {
	var result Players
	for _, p := range players {
		if !p.IsBot {
			result = append(result, p)
		}
	}
	return result
}

func (players Players) Bots() Players {
	var result Players
	for _, p := range players {
		if p.IsBot {
			result = append(result, p)
		}
	}
	return result
}

func (players Players) BotsCount() int {
	count := 0
	for _, player := range players {
		if player.IsBot {
			count++
		}
	}
	return count
}
