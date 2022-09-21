package shell

import (
	"github.com/diyor28/not-agar/src/utils"
	"math"
)

func New(p int, r float32) *Shell {
	var points []*Point
	segLen := r * 2 * math.Pi / float32(p)
	for i := 0; i < p; i++ {
		angle := math.Pi * 2 * float32(i) / float32(p)
		x := r * float32(math.Sin(float64(angle)))
		y := r * float32(math.Cos(float64(angle)))
		point := &Point{
			oX:  x,
			oY:  y,
			X:   x,
			Y:   y,
			Len: segLen,
		}
		points = append(points, point)
	}
	return &Shell{points, r}
}

func CalcMaxXY(x float32, y float32, r float32) (float32, float32) {
	dist := utils.Distance(0, x, 0, y)
	sinA := float32(math.Abs(float64(x / dist)))
	cosA := float32(math.Abs(float64(y / dist)))
	return r * sinA * 1.5, r * cosA * 1.5
}

type Shell struct {
	Points []*Point
	radius float32
}

func (s *Shell) SetRadius(r float32) {
	segLen := r * math.Pi / float32(len(s.Points))
	var ratio float32 = 1
	if s.radius > 0 {
		ratio = r / s.radius
	}
	for i := 0; i < len(s.Points); i++ {
		curr := s.Points[i]
		curr.SetLen(segLen)
		curr.oX *= ratio
		curr.oY *= ratio
		curr.X *= ratio
		curr.Y *= ratio
	}
	s.radius = r
}

func (s *Shell) SmoothPoints() {
	r := s.radius
	for i, p := range s.Points {
		angle := math.Pi * 2 * float64(i) / float64(len(s.Points))
		p.X = p.X*0.99 + float32(math.Sin(angle))*r*0.01
		p.Y = p.Y*0.99 + float32(math.Cos(angle))*r*0.01
	}
}

func (s *Shell) averageAngle(i int, n int) float32 {
	var sum float32 = 0
	start := i - n
	if n > i {
		start = 0
	}
	for _, p := range s.Points[start:i] {
		sum += p.Angle
	}
	return sum / float32(i-start)
}

func (s *Shell) SmoothAngles() {
	for i := 1; i < len(s.Points); i++ {
		avgAngle := s.averageAngle(i, 5)
		curr := s.Points[i]
		//prev := s.Points[i-1]
		curr.Angle = avgAngle
	}
}

func (s *Shell) MovePoint(pIdx int, x float32, y float32) {
	var points []*Point
	for i := pIdx; i < len(s.Points); i++ {
		points = append(points, s.Points[i])
	}
	for i := 0; i < pIdx; i++ {
		points = append(points, s.Points[i])
	}
	head := points[0]
	r := s.radius
	mX, mY := CalcMaxXY(head.X, head.Y, r)
	head.Follow(x, y, mX, mY)
	for i := 1; i < len(points); i++ {
		curr := points[i]
		prev := points[i-1]
		mX, mY := CalcMaxXY(curr.X, curr.Y, r)
		curr.Follow(prev.oX, prev.oY, mX, mY)
	}

	tail := s.Points[len(points)-1]
	tail.SetOrigin(head.X, head.Y, s.radius)
	for i := len(points) - 2; i >= 0; i-- {
		curr := points[i]
		prev := points[i+1]
		curr.SetOrigin(prev.X, prev.Y, s.radius)
	}
}

func (s *Shell) ClosestPoint(x float32, y float32) int {
	closest := 0
	var maxDist float32 = 0
	for i, p := range s.Points {
		dist := utils.Distance(x, p.X, y, p.Y)
		if dist > maxDist {
			closest = i
			maxDist = dist
		}
	}
	return closest
}

func (s *Shell) FarthestApart(pIdx int) (int, int) {
	var points []*Point
	for i := pIdx; i < len(s.Points); i++ {
		points = append(points, s.Points[i])
	}
	for i := 0; i < pIdx; i++ {
		points = append(points, s.Points[i])
	}
	j := 0
	k := len(points) - 1
	sIdx := 0
	eIdx := 0
	var maxDist float32 = 0
	for j < k {
		dist := utils.Distance(points[j].X, points[k].X, points[j].Y, points[k].Y)
		if dist > maxDist {
			maxDist = dist
			sIdx = j
			eIdx = k
		}
		j++
		k--
	}
	sIdx = (sIdx + pIdx) % len(points)
	eIdx = (eIdx + pIdx) % len(points)
	if eIdx > sIdx {
		return sIdx, eIdx
	}
	return eIdx, sIdx
}
