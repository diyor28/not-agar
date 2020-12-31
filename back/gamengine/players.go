package gamengine

import (
	"github.com/diyor28/not-agar/utils"
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
	Uuid         string  `json:"-"`
	Nickname     string  `json:"nickname"`
	VelocityX    float32 `json:"-"`
	VelocityY    float32 `json:"-"`
	Color        [3]int  `json:"color"`
	X            float32 `json:"x"`
	Y            float32 `json:"y"`
	Weight       float32 `json:"weight"`
	Accelerating bool    `json:"-"`
	WeightBurned float32 `json:"-"`
	Speed        float32 `json:"-"`
	Zoom         float32 `json:"-"`
	IsBot        bool    `json:"-"`
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

func (pl *Player) getWeight() float32 {
	return pl.Weight
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
	diff := utils.CalcDistance(pl.X, anotherPlayer.X, pl.Y, anotherPlayer.Y)
	radius1 := float64(pl.Weight) / 2
	radius2 := float64(anotherPlayer.Weight) / 2
	closeEnough := diff < float32(math.Abs(radius1-radius2))
	bigEnough := pl.Weight*0.85 > anotherPlayer.Weight
	return bigEnough && closeEnough
}

func (pl *Player) addWeight(weight float32) {
	sign := float32(math.Copysign(1, float64(weight)))
	nWeight := math.Sqrt(float64(pl.Weight*pl.Weight + weight*weight*sign))
	pl.setWeight(float32(nWeight))
}

func (pl *Player) eatEntity(entity interface{ getWeight() float32 }) {
	pl.setWeight(newWeight(pl.Weight, entity.getWeight()))
}

func (pl *Player) setWeight(weight float32) {
	clippedWeight := utils.Clip(weight, MinWeight, MaxWeight)
	pl.Weight = clippedWeight
	pl.Speed = getSpeedFromWeight(clippedWeight)
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

func (pl *Player) makeMove(gameMap *GameMap) {
	//pls := gameMap.nearByPlayers(gBot.Player)
	foods := gameMap.nearByFood(pl, 1)
	closestFood := foods[0]
	pl.updateDirection(closestFood.X, closestFood.Y)
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
