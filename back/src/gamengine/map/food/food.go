package food

import (
	"github.com/diyor28/not-agar/src/gamengine/map/entity"
	"github.com/diyor28/not-agar/src/utils"
	"github.com/frankenbeanies/uuid4"
)

type EntityInterface interface {
	getX() float32
	getY() float32
}

type Food struct {
	*entity.Entity
	Color [3]int `json:"color"`
}

func New(x float32, y float32, weight float32) *Food {
	food := &Food{
		Entity: &entity.Entity{
			Uuid:   uuid4.New().String(),
			X:      x,
			Y:      y,
			Weight: weight,
		},
		Color: utils.RandomColor(),
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

func (foods Foods) Closest(player entity.Interface, kClosest int) []Food {
	foodCopy := foods.asValues()
	totalNumFood := len(foodCopy)
	foodDistances := make(map[string]float32, totalNumFood)
	for _, f := range foodCopy {
		foodDistances[f.Uuid] = utils.CalcDistance(player.GetX(), f.X, player.GetY(), f.Y)
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
