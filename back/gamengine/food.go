package gamengine

import (
	"github.com/diyor28/not-agar/src/github.com/frankenbeanies/uuid4"
)

type Food struct {
	*Entity
	Color [3]int `json:"color"`
}

func NewFood(x float32, y float32, weight float32) *Food {
	food := &Food{
		Entity: &Entity{
			Uuid: uuid4.New().String(),
			X:    x,
			Y:    y,
		},
		Color: randomColor(),
	}
	food.setWeight(weight)
	return food
}

