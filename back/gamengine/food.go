package gamengine

type Food struct {
	Uuid   string  `json:"-"`
	Color  [3]int  `json:"color"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Weight float32 `json:"weight"`
}

func (f *Food) getWeight() float32 {
	return f.Weight
}
