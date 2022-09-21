package players

import (
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/gamengine/constants"
	"github.com/diyor28/not-agar/src/gamengine/map/entity"
	"github.com/diyor28/not-agar/src/gamengine/map/players/shell"
	"github.com/diyor28/not-agar/src/utils"
	"github.com/frankenbeanies/uuid4"
	"math"
	"time"
)

type Player struct {
	*entity.Entity
	Nickname     string
	VelocityX    float32
	VelocityY    float32
	DirectionX   float32
	DirectionY   float32
	Color        [3]uint8
	Shell        *shell.Shell
	Accelerating bool
	Speed        float32
	Zoom         float32
	IsBot        bool
	IsDead       bool
}

func NewPlayer(x float32, y float32, weight float32, nickname string, isBot bool) *Player {
	sh := shell.New(50, weight/2)
	player := &Player{
		Entity: &entity.Entity{
			Uuid: uuid4.New().String(),
			X:    x,
			Y:    y,
		},
		Speed:        0,
		Nickname:     nickname,
		IsBot:        isBot,
		Accelerating: false,
		Zoom:         1,
		Color:        utils.RandomColor(),
		Shell:        sh,
	}
	player.SetWeight(weight)
	return player
}

func (pl *Player) PassiveWeightLoss() {
	if pl.Weight > constants.MinWeight*3 {
		pl.SetWeight(pl.Weight * 0.99999)
	}
}

func (pl *Player) FoodEatable(food entity.Interface) bool {
	diff := utils.Distance(pl.X, food.GetX(), pl.Y, food.GetY())
	return diff < pl.Weight/2
}

func (pl *Player) PlayerEatable(player entity.Interface) bool {
	radius1 := float64(pl.GetWeight()) / 2
	radius2 := float64(player.GetWeight()) / 2
	surfArea1 := utils.SurfaceArea64(radius1)
	surfArea2 := utils.SurfaceArea64(radius2)
	dist := float64(utils.Distance(pl.X, player.GetX(), pl.Y, player.GetY()))
	interSection := utils.IntersectionArea(radius1, radius2, dist) / surfArea2
	closeEnough := interSection > 0.85
	bigEnough := surfArea1*0.85 > surfArea2
	return bigEnough && closeEnough
}

func (pl *Player) UpdatePosition(delta time.Duration) {
	dT := float32(delta.Seconds())
	speed := utils.Modulus(0, pl.VelocityX, 0, pl.VelocityY)
	drag := 0.5 * constants.Density * utils.Square(speed) * constants.DragCoefficient * utils.SurfaceArea(pl.Weight/1000)
	fX := pl.DirectionX * constants.Force
	fY := pl.DirectionY * constants.Force
	var cosA float32 = 0
	var sinA float32 = 0
	if speed != 0 {
		cosA = pl.VelocityX / speed
		sinA = pl.VelocityY / speed
	}
	dragX := -1 * cosA * drag
	dragY := -1 * sinA * drag
	aX := (fX + dragX) / pl.Weight
	aY := (fY + dragY) / pl.Weight
	pl.VelocityX += aX * dT
	pl.VelocityY += aY * dT
	newX := pl.X + pl.VelocityX*dT
	newY := pl.Y + pl.VelocityY*dT
	newZoom := (constants.MinWeight/pl.Weight)*(constants.MaxZoom-constants.MinZoom) + constants.MinZoom
	if newX >= constants.MaxXY-5 {
		pl.VelocityX *= 0.05
	}
	if newY >= constants.MaxXY-5 {
		pl.VelocityY *= 0.05
	}
	pl.X = utils.Clip(newX, 1, constants.MaxXY)
	pl.Y = utils.Clip(newY, 1, constants.MaxXY)
	pl.Zoom = utils.Clip(newZoom, constants.MinZoom, constants.MaxZoom)

	r := pl.Weight / 2
	closest := pl.Shell.ClosestPoint(r*cosA, r*sinA)
	j, k := pl.Shell.FarthestApart(closest)
	for i := j; i < k; i++ {
		p := pl.Shell.Points[i]
		mX, mY := shell.CalcMaxXY(p.X, p.Y, r)
		p.Follow(p.X-pl.VelocityX*dT*0.05, p.Y-pl.VelocityY*dT*0.05, mX, mY)
	}
	pl.Shell.SmoothAngles()
	pl.Shell.SmoothPoints()
	for _, p := range pl.Shell.Points {
		p.X = utils.Clip(pl.X+p.X, 0, constants.MaxXY) - pl.X
		p.Y = utils.Clip(pl.Y+p.Y, 0, constants.MaxXY) - pl.Y
	}
}

func (pl *Player) UpdateDirection(newX float32, newY float32) {
	dist := float64(utils.Distance(pl.X, newX, pl.Y, newY))
	diffX := float64(newX - pl.X)
	diffY := float64(newY - pl.Y)
	pl.DirectionX = float32(diffX / dist)
	pl.DirectionY = float32(diffY / dist)
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
	pl.Shell.SetRadius(weight / 2)
	pl.Entity.SetWeight(weight)
}

type Players []*Player

func (players Players) Get(uuid string) (*Player, error) {
	for _, p := range players {
		if p.Uuid == uuid {
			return p, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no player for uuid %s found", uuid))
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
		playersDistances[p.Uuid] = utils.Distance(player.X, p.X, player.Y, p.Y)
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
