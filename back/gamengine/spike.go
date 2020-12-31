package gamengine

import "github.com/diyor28/not-agar/utils"

type Spike struct {
	Uuid   string  `json:"-"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Weight float32 `json:"weight"`
}

func (s *Spike) collided(player *Player) bool {
	diff := utils.CalcDistance(s.X, player.X, s.Y, player.Y)
	closeEnough := diff < s.Weight/2
	bigEnough := player.Weight > s.Weight*0.85
	return closeEnough && bigEnough
}
