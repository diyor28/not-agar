package gamengine

import (
	"fmt"
	"math"
	"utils"
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
	Uuid       string  `json:"-"`
	Nickname   string  `json:"nickname"`
	DirectionX int     `json:"-"`
	DirectionY int     `json:"-"`
	Color      [3]int  `json:"color"`
	X          float32 `json:"x"`
	Y          float32 `json:"y"`
	Weight     float32 `json:"weight"`
	Speed      float32 `json:"-"`
	Zoom       float32 `json:"-"`
	IsBot      bool    `json:"-"`
}

func (pl Player) getSelfPlayer() SelfPlayer {
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
	if pl.Weight > minWeight*3 {
		pl.Weight = pl.Weight * 0.9999
	}
}

func (pl Player) foodEatable(food *Food) bool {
	diff := utils.CalcDistance(pl.X, food.X, pl.Y, food.Y)
	return diff < pl.Weight/2
}

func (pl *Player) updatePosition() {
	newX := float64(pl.X + pl.Speed*float32(pl.DirectionX))
	newY := float64(pl.Y + pl.Speed*float32(pl.DirectionY))
	pl.X = float32(math.Max(math.Min(newX, maxXY), minXY))
	pl.Y = float32(math.Max(math.Min(newY, maxXY), minXY))
	newZoom := float64((minWeight/pl.Weight)*(maxZoom-minZoom)) + minZoom
	pl.Zoom = float32(math.Max(newZoom, minZoom))
}

func (pl *Player) updateDirection(directionX int, directionY int) {
	if !validDirection[directionX] {
		fmt.Println("Sent an invalid direction", directionX)
		return
	}
	if !validDirection[directionY] {
		fmt.Println("Sent an invalid direction", directionY)
		return
	}
	pl.DirectionX = directionX
	pl.DirectionY = directionY
}

func (pl Player) playerEatable(anotherPlayer *Player) bool {
	diff := utils.CalcDistance(pl.X, anotherPlayer.X, pl.Y, anotherPlayer.Y)
	radius1 := float64(pl.Weight) / 2
	radius2 := float64(anotherPlayer.Weight) / 2
	closeEnough := diff < float32(math.Abs(radius1-radius2))
	bigEnough := pl.Weight*0.81 > anotherPlayer.Weight
	return bigEnough && closeEnough
}

func (pl *Player) eatPlayer(anotherPlayer *Player) {
	pl.Weight = newWeight(pl.Weight, anotherPlayer.Weight)
	pl.Speed = getSpeedFromWeight(anotherPlayer.Weight)
}

func (pl *Player) eatFood(food *Food) {
	pl.Weight = newWeight(pl.Weight, food.Weight)
	pl.Speed = getSpeedFromWeight(pl.Weight)
}

type DirectionEvaluator struct {
	Direction string
	Actions   []int
	Foods     []Food
	Sum       float32
}

func (dEval *DirectionEvaluator) appendFood(food Food) {
	dEval.Foods = append(dEval.Foods, food)
	dEval.Sum = dEval.Sum + food.Weight
}

func (dEval *DirectionEvaluator) isAction(dirX int, dirY int) bool {
	return (dEval.Actions[0] == dirX) && (dEval.Actions[1] == dirY)
}

func (pl *Player) closesFood(gameMap GameMap, k int) []Food {
	var result []Food
	foods := gameMap.Foods
	foodCount := len(foods)
	if foodCount < k {
		k = foodCount
	}
	for i := 0; i < k; i++ {
		var minIdx = i
		for j := i; j < foodCount; j++ {
			if foods[j].Weight < foods[minIdx].Weight {
				minIdx = j
			}
		}
		foods[i], foods[minIdx] = foods[minIdx], foods[i]
	}
	return result
}

func (pl *Player) makeMove(gameMap *GameMap) {
	//pls := gameMap.nearByPlayers(gBot.Player)
	foods := gameMap.nearByFood(pl)
	closestFood := foods[0]
	var closesDist float32 = 0.0
	for _, f := range foods {
		dist := utils.CalcDistance(pl.X, f.X, pl.Y, f.Y)
		if dist < closesDist {
			closestFood = f
			closesDist = dist
		}
	}

	diffX := float64(closestFood.X - pl.X)
	diffY := float64(closestFood.Y - pl.Y)
	directionX := 0
	directionY := 0
	if diffX != 0 {
		directionX = int(math.Copysign(1, diffX))
	}
	if diffY != 0 {
		directionY = int(math.Copysign(1, diffY))
	}
	pl.updateDirection(directionX, directionY)
}

//
//func (pl *Player) makeMove(gameMap *GameMap) {
//	//pls := gameMap.nearByPlayers(gBot.Player)
//	foods := gameMap.nearByFood(*pl)
//	time.Sleep(300)
//	var directions = map[string][]int{
//		"tr": {1, 1},
//		"tl": {-1, 1},
//		"bl": {-1, -1},
//		"br": {1, -1},
//	}
//	var evaluators []DirectionEvaluator
//	for dir, actions := range directions {
//		evaluators = append(evaluators, DirectionEvaluator{
//			Direction: dir,
//			Actions:   actions,
//			Sum:       0,
//		})
//	}
//
//	for _, f := range foods {
//		diffX := float64(f.X - pl.X)
//		diffY := float64(f.Y - pl.Y)
//		dirX := int(diffX / math.Abs(diffX))
//		dirY := int(diffY / math.Abs(diffY))
//		for i, evaluator := range evaluators {
//			if evaluator.isAction(dirX, dirY) {
//				ev := &evaluators[i]
//				ev.appendFood(f)
//			}
//		}
//	}
//
//	bestDirection := evaluators[0]
//	for _, evaluator := range evaluators {
//		if bestDirection.Sum < evaluator.Sum {
//			bestDirection = evaluator
//		}
//	}
//
//	if len(pl.lastAction) == 0 {
//		pl.lastAction = bestDirection.Actions
//	}
//	var directionX = pl.lastAction[0]
//	var directionY = pl.lastAction[1]
//	if pl.count < 5 {
//		pl.count++
//	} else {
//		pl.lastAction = bestDirection.Actions
//		pl.count = 0
//	}
//	gameMap.updatePlayer(ServerRequest{
//		Uuid:       pl.Uuid,
//		DirectionX: directionX,
//		DirectionY: directionY,
//	})
//}
