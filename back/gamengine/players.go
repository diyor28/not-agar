package gamengine

import (
	"errors"
	"github.com/diyor28/not-agar/utils"
	"github.com/frankenbeanies/uuid4"
	"math"
)

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
	*Entity
	Nickname     string  `json:"nickname"`
	VelocityX    float32 `json:"-"`
	VelocityY    float32 `json:"-"`
	Color        [3]int  `json:"color"`
	Accelerating bool    `json:"-"`
	WeightBurned float32 `json:"-"`
	Speed        float32 `json:"-"`
	Zoom         float32 `json:"-"`
	IsBot        bool    `json:"-"`
}

func NewPlayer(x float32, y float32, weight float32, nickname string, isBot bool) *Player {
	player := &Player{
		Entity: &Entity{
			Uuid: uuid4.New().String(),
			X:    x,
			Y:    y,
		},
		Speed:        MaxSpeed,
		Nickname:     nickname,
		IsBot:        isBot,
		Accelerating: false,
		Zoom:         1,
		Color:        randomColor(),
	}
	player.setWeight(weight)
	return player
}

func (pl *Player) getSelfPlayer() SelfPlayer {
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

func (pl *Player) passiveWeightLoss() {
	if pl.Weight > MinWeight*3 {
		pl.setWeight(pl.Weight * 0.99999)
	}
}

func (pl *Player) foodEatable(food *Food) bool {
	diff := utils.CalcDistance(pl.X, food.X, pl.Y, food.Y)
	return diff < pl.Weight/2
}

func (pl *Player) updatePosition(gMap *GameMap) {
	speed := pl.Speed
	if pl.Accelerating {
		speed = utils.Clip(speed*2, MinSpeed, MaxSpeed)
		pl.WeightBurned += pl.Weight / 100
		pl.Accelerating = false
	}
	newX := pl.X + speed*pl.VelocityX
	newY := pl.Y + speed*pl.VelocityY
	newZoom := (MinWeight/pl.Weight)*(MaxZoom-MinZoom) + MinZoom
	pl.X = utils.Clip(newX, MinXY, MaxXY)
	pl.Y = utils.Clip(newY, MinXY, MaxXY)
	pl.Zoom = utils.Clip(newZoom, MinZoom, MaxZoom)
	if pl.WeightBurned >= FoodWeight {
		pl.WeightBurned -= FoodWeight
		pl.addWeight(-FoodWeight)
		x := pl.X - speed*pl.VelocityX
		y := pl.Y - speed*pl.VelocityY
		x -= float32(math.Copysign(float64(pl.Weight), float64(pl.VelocityX))) * (pl.VelocityX * pl.VelocityX)
		y -= float32(math.Copysign(float64(pl.Weight), float64(pl.VelocityY))) * (pl.VelocityY * pl.VelocityY)
		gMap.createFood(x, y)
	}
}

func (pl *Player) updateDirection(newX float32, newY float32) {
	dist := float64(utils.CalcDistance(pl.X, newX, pl.Y, newY))
	diffX := float64(newX - pl.X)
	diffY := float64(newY - pl.Y)
	velocityX := diffX / dist
	velocityY := diffY / dist
	pl.VelocityX = float32(velocityX)
	pl.VelocityY = float32(velocityY)
}

func (pl Player) canEat(anotherPlayer *Player) bool {
	dist := float64(utils.CalcDistance(pl.X, anotherPlayer.X, pl.Y, anotherPlayer.Y))
	radius1 := float64(pl.Weight) / 2
	radius2 := float64(anotherPlayer.Weight) / 2
	interSection := utils.IntersectionArea(radius1, radius2, dist) / utils.SurfaceArea64(radius2)
	closeEnough := interSection > 0.85
	bigEnough := pl.surfaceArea*0.85 > anotherPlayer.surfaceArea
	return bigEnough && closeEnough
}

func (pl *Player) addWeight(weight float32) {
	sign := float32(math.Copysign(1, float64(weight)))
	nWeight := math.Sqrt(float64(pl.Weight*pl.Weight + weight*weight*sign))
	pl.setWeight(float32(nWeight))
}

func (pl *Player) eatEntity(entity interface{ getWeight() float32 }) {
	pl.addWeight(entity.getWeight())
}

func (pl *Player) setWeight(weight float32) {
	pl.Entity.setWeight(weight)
	pl.Speed = getSpeedFromWeight(pl.Weight)
}

func (pl *Player) makeMove(gameMap *GameMap) {1
	foods := gameMap.Foods.closest(pl, 1)
	closestFood := foods[0]
	pl.updateDirection(closestFood.X, closestFood.Y)
}

type Players []*Player

func (players Players) get(uuid string) (*Player, error) {
	for _, p := range players {
		if p.Uuid == uuid {
			return p, nil
		}
	}
	return nil, errors.New("no player found")
}

func (players Players) asValues() []Player {
	var result []Player
	for _, p := range players {
		result = append(result, *p)
	}
	return result
}

func (players Players) exclude(uuid string) Players {
	var result Players
	for _, p := range players {
		if p.Uuid != uuid {
			result = append(result, p)
		}
	}
	return result
}

func (players Players) largest(k int) []Player {
	playersValues := players.asValues()
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

func (players Players) closest(player *Player, kClosest int) Players {
	otherPlayers := players.exclude(player.Uuid)
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

func (players Players) update(uuid string, newX float32, newY float32) (*Player, error) {
	player, err := players.get(uuid)
	if err != nil {
		return nil, err
	}
	player.updateDirection(newX, newY)
	return player, nil
}

func (players Players) real() Players {
	var result Players
	for _, p := range players {
		if !p.IsBot {
			result = append(result, p)
		}
	}
	return result
}

func (players Players) bots() Players {
	var result Players
	for _, p := range players {
		if p.IsBot {
			result = append(result, p)
		}
	}
	return result
}

func (players Players) botsCount() int {
	count := 0
	for _, player := range players {
		if player.IsBot {
			count++
		}
	}
	return count
}
