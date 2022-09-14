package spikes

import (
	"github.com/diyor28/not-agar/src/gamengine/map/entity"
	"github.com/diyor28/not-agar/src/gamengine/map/players"
	"github.com/diyor28/not-agar/src/utils"
	"github.com/frankenbeanies/uuid4"
)

type Spike struct {
	*entity.Entity
}

func New(x float32, y float32, weight float32) *Spike {
	spike := &Spike{
		&entity.Entity{
			Uuid: uuid4.New().String(),
			X:    x,
			Y:    y,
		},
	}
	spike.SetWeight(weight)
	return spike
}

func (s *Spike) Collided(player *players.Player) bool {
	plRadius := float64(player.Weight / 2)
	sRadius := float64(s.Weight / 2)
	dist := float64(utils.CalcDistance(s.X, player.X, s.Y, player.Y))
	closeEnough := utils.IntersectionArea(plRadius, sRadius, dist)/utils.SurfaceArea64(sRadius) > 0.5
	bigEnough := plRadius > sRadius
	return closeEnough && bigEnough
}

type Spikes []*Spike

func (spikes Spikes) asValues() []Spike {
	var result []Spike
	for _, p := range spikes {
		result = append(result, *p)
	}
	return result
}

func (spikes Spikes) closest(player *players.Player, kClosest int) []Spike {
	spikesCopy := spikes.asValues()
	totalSpikes := len(spikesCopy)
	spikeDistances := make(map[string]float32, totalSpikes)
	for _, s := range spikesCopy {
		spikeDistances[s.Uuid] = utils.CalcDistance(player.X, s.X, player.Y, s.Y)
	}
	numResults := kClosest
	if kClosest > totalSpikes {
		numResults = totalSpikes
	}
	for i := 0; i < numResults; i++ {
		var minIdx = i
		for j := i + 1; j < totalSpikes; j++ {
			if spikeDistances[spikesCopy[j].Uuid] < spikeDistances[spikesCopy[minIdx].Uuid] {
				minIdx = j
			}
		}
		spikesCopy[i], spikesCopy[minIdx] = spikesCopy[minIdx], spikesCopy[i]
	}
	return spikesCopy[:numResults]
}

func (spikes Spikes) RandXY(minDistance float32) (float32, float32) {
	x, y := utils.RandXY()
	for _, s := range spikes {
		radius := s.Weight / 2
		eX, eY := s.X, s.Y
		dist := utils.CalcDistance(eX, x, eY, y)
		if dist-radius < minDistance {
			return spikes.RandXY(minDistance)
		}
	}
	return x, y
}
