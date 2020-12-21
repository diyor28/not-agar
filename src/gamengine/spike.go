package gamengine

type Spike struct {
	Uuid   string  `json:"-"`
	Color  [3]int  `json:"color"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Weight float32 `json:"weight"`
}

func (s Spike) collided(player Player) {

}
