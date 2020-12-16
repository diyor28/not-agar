package gamengine

type Food struct {
	Uuid   string  `json:"-"`
	Color  string  `json:"color"`
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Weight float32 `json:"weight"`
}
