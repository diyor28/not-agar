package shell

import (
	"github.com/diyor28/not-agar/src/utils"
	"math"
)

type Point struct {
	oX    float32
	oY    float32
	X     float32
	Y     float32
	Len   float32
	Angle float32
}

func (p *Point) PointAt(x float32, y float32) {
	p.Angle = utils.Clip(float32(math.Atan2(float64(y-p.oY), float64(x-p.oX))), 0, math.Pi/18)
}

func (p *Point) MoveTo(x float32, y float32, mX float32, mY float32) {
	p.X = x
	p.Y = y
	p.oX = utils.Clip(x-p.Len*float32(math.Cos(float64(p.Angle))), -mX, mY)
	p.oY = utils.Clip(y-p.Len*float32(math.Sin(float64(p.Angle))), -mX, mY)
}

func (p *Point) Follow(x float32, y float32, mX float32, mY float32) {
	x = utils.Clip(x, -mX, mX)
	y = utils.Clip(y, -mY, mY)
	p.PointAt(x, y)
	p.MoveTo(x, y, mX, mY)
}

func (p *Point) SetOrigin(x float32, y float32, r float32) {
	x = utils.Clip(x, -r*1.2, r*1.2)
	y = utils.Clip(y, -r*1.2, r*1.2)
	p.Angle = float32(math.Atan2(float64(p.Y-y), float64(p.X-x)))
	p.oX = x
	p.oY = y
	p.X = utils.Clip(x+p.Len*float32(math.Cos(float64(p.Angle))), -r*1.2, r*1.2)
	p.Y = utils.Clip(y+p.Len*float32(math.Sin(float64(p.Angle))), -r*1.2, r*1.2)
}

func (p *Point) SetLen(newLen float32) {
	p.Len = newLen
}
