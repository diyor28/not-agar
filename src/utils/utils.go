package utils

import "math"

func CalcDistance(x1 float32, x2 float32, y1 float32, y2 float32) float32 {
	return float32(math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2)))
}