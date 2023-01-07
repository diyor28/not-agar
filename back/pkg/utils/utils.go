package utils

import (
	"github.com/diyor28/not-agar/pkg/constants"
	"math"
	"math/rand"
)

var colors = [][3]uint8{{255, 21, 21}, {255, 243, 21}, {21, 87, 255}, {21, 255, 208}, {255, 21, 224}}

func RandXY() (float32, float32) {
	return float32(rand.Intn(constants.MaxXY)), float32(rand.Intn(constants.MaxXY))
}

func RandomColor() [3]uint8 {
	return colors[rand.Intn(len(colors))]
}

func Square(x float32) float32 {
	return x * x
}

func Square64(x float64) float64 {
	return x * x
}

func Modulus(x1 float32, x2 float32, y1 float32, y2 float32) float32 {
	return float32(Modulus64(float64(x1), float64(x2), float64(y1), float64(y2)))
}

func Modulus64(x1 float64, x2 float64, y1 float64, y2 float64) float64 {
	return math.Sqrt(Square64(x2-x1) + Square64(y2-y1))
}

func Distance(x1 float32, x2 float32, y1 float32, y2 float32) float32 {
	return Modulus(x1, x2, y1, y2)
}

func Distance64(x1 float64, x2 float64, y1 float64, y2 float64) float64 {
	return Modulus64(x1, x2, y1, y2)
}

func Lerp() {

}

func Clip(v float32, min float32, max float32) float32 {
	return float32(math.Min(math.Max(float64(v), float64(min)), float64(max)))
}

func Clip64(v float64, min float64, max float64) float64 {
	return math.Min(math.Max(v, min), max)
}

func SurfaceArea(r float32) float32 {
	return float32(math.Pi * math.Pow(float64(r), 2))
}

func SurfaceArea64(r float64) float64 {
	return math.Pi * math.Pow(r, 2)
}

func IntersectionArea(r1 float64, r2 float64, dist float64) float64 {
	if dist < math.Abs(r1-r2) {
		return SurfaceArea64(math.Min(r1, r2))
	}
	if dist > math.Abs(r1+r2) {
		return 0
	}
	dist1 := (Square64(r1) - Square64(r2) + Square64(dist)) / (2 * dist)
	dist2 := dist - dist1
	h1 := math.Sqrt(Square64(r1) - Square64(dist1))
	h2 := math.Sqrt(Square64(r2) - Square64(dist2))
	return Square64(r1)*math.Acos(dist1/r1) - dist1*h1 + Square64(r2)*math.Acos(dist2/r2) - dist2*h2
}
