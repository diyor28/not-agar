package gamengine

import (
	"github.com/diyor28/not-agar/utils"
	"github.com/frankenbeanies/uuid4"
)

type Spike struct {
	*Entity
}

func NewSpike(x float32, y float32, weight float32) *Spike {
	spike := &Spike{
		&Entity{
			Uuid: uuid4.New().String(),
			X:    x,
			Y:    y,
		},
	}
	spike.setWeight(weight)
	return spike
}

func (s *Spike) collided(player *Player) bool {
	plRadius := float64(player.Weight / 2)
	sRadius := float64(s.Weight / 2)
	dist := float64(utils.CalcDistance(s.X, player.X, s.Y, player.Y))
	closeEnough := utils.IntersectionArea(plRadius, sRadius, dist)/utils.SurfaceArea64(sRadius) > 0.5
	bigEnough := plRadius > sRadius
	return closeEnough && bigEnough
}
