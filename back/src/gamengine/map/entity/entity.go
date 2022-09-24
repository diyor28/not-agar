package entity

import (
	"github.com/diyor28/not-agar/src/gamengine/constants"
	"github.com/diyor28/not-agar/src/utils"
)

type Id uint32

type Interface interface {
	GetX() float32
	GetY() float32
	GetWeight() float32
}

type Entity struct {
	Id          Id
	X           float32
	Y           float32
	Weight      float32
	surfaceArea float32
}

func (e *Entity) SetWeight(weight float32) {
	clippedWeight := utils.Clip(weight, constants.MinWeight, constants.MaxWeight)
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
