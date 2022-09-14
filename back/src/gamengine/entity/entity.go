package entity

import (
	"errors"
	"github.com/diyor28/not-agar/src/gamengine"
	"github.com/diyor28/not-agar/src/utils"
)

type Entity struct {
	Uuid        string  `json:"-"`
	X           float32 `json:"x"`
	Y           float32 `json:"y"`
	Weight      float32 `json:"weight"`
	surfaceArea float32
}

func (e *Entity) SetWeight(weight float32) {
	clippedWeight := utils.Clip(weight, gamengine.MinWeight, gamengine.MaxWeight)
	e.Weight = clippedWeight
	e.surfaceArea = utils.SurfaceArea(weight / 2)
}

func (e *Entity) GetWeight() float32 {
	return e.Weight
}

func (e *Entity) GetX() float32 {
	return e.X
}

func (e *Entity) GetY() float32 {
	return e.Y
}

func (e *Entity) GetXY() (float32, float32) {
	return e.X, e.Y
}

type Entities struct {
	Items []*Entity
}

func (entities *Entities) GetItem(uuid string) (*Entity, error) {
	for i := range entities.Items {
		if entities.Items[i].Uuid == uuid {
			return entities.Items[i], nil
		}
	}
	return nil, errors.New("no player found")
}
