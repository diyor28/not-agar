package gamengine

import (
	"github.com/diyor28/not-agar/src/github.com/frankenbeanies/uuid4"
	"github.com/diyor28/not-agar/utils"
)

type Food struct {
	*Entity
	Color [3]int `json:"color"`
}

func NewFood(x float32, y float32, weight float32) *Food {
	food := &Food{
		Entity: &Entity{
			Uuid:   uuid4.New().String(),
			X:      x,
			Y:      y,
			Weight: weight,
		},
		Color: randomColor(),
	}
	return food
}

type Foods []*Food

func (foods Foods) asValues() []Food {
	var result []Food
	for _, f := range foods {
		result = append(result, *f)
	}
	return result
}

func (foods Foods) closest(player *Player, kClosest int) []Food {
	foodCopy := foods.asValues()
	totalNumFood := len(foodCopy)
	foodDistances := make(map[string]float32, totalNumFood)
	for _, f := range foodCopy {
		foodDistances[f.Uuid] = utils.CalcDistance(player.X, f.X, player.Y, f.Y)
	}
	numResults := kClosest
	if kClosest > totalNumFood {
		numResults = totalNumFood
	}
	for i := 0; i < numResults; i++ {
		var minIdx = i
		for j := i + 1; j < totalNumFood; j++ {
			if foodDistances[foodCopy[j].Uuid] < foodDistances[foodCopy[minIdx].Uuid] {
				minIdx = j
			}
		}
		foodCopy[i], foodCopy[minIdx] = foodCopy[minIdx], foodCopy[i]
	}
	return foodCopy[:numResults]
}
