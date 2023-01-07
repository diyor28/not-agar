package spikes

import (
	"github.com/diyor28/not-agar/cmd/gamengine/map/entity"
	"github.com/diyor28/not-agar/cmd/gamengine/map/players"
	"github.com/diyor28/not-agar/pkg/utils"
)

type Spike struct {
	*entity.Entity
}

type Spikes struct {
	Spikes []*Spike
	lastId entity.Id
}

func (s *Spikes) New(x float32, y float32, weight float32) *Spike {
	s.lastId++
	id := s.lastId
	spike := &Spike{
		&entity.Entity{
			Id: id,
			X:  x,
			Y:  y,
		},
	}
	spike.SetWeight(weight)
	s.Append(spike)
	return spike
}

func (s *Spikes) Append(spike *Spike) {
	s.Spikes = append(s.Spikes, spike)
}

func (s *Spikes) asValues() []Spike {
	var result []Spike
	for _, p := range s.Spikes {
		result = append(result, *p)
	}
	return result
}

func (s *Spikes) closest(player *players.Player, kClosest int) []Spike {
	spikesCopy := s.asValues()
	totalSpikes := len(spikesCopy)
	spikeDistances := make(map[entity.Id]float32, totalSpikes)
	for _, s := range spikesCopy {
		spikeDistances[s.Id] = utils.Distance(player.X, s.X, player.Y, s.Y)
	}
	numResults := kClosest
	if kClosest > totalSpikes {
		numResults = totalSpikes
	}
	for i := 0; i < numResults; i++ {
		var minIdx = i
		for j := i + 1; j < totalSpikes; j++ {
			if spikeDistances[spikesCopy[j].Id] < spikeDistances[spikesCopy[minIdx].Id] {
				minIdx = j
			}
		}
		spikesCopy[i], spikesCopy[minIdx] = spikesCopy[minIdx], spikesCopy[i]
	}
	return spikesCopy[:numResults]
}

func (s *Spikes) RandXY(minDistance float32) (float32, float32) {
	x, y := utils.RandXY()
	for _, spike := range s.Spikes {
		radius := spike.Weight / 2
		eX, eY := spike.X, spike.Y
		dist := utils.Distance(eX, x, eY, y)
		if dist-radius < minDistance {
			return s.RandXY(minDistance)
		}
	}
	return x, y
}

func (s *Spike) Collided(player *players.Player) bool {
	plRadius := float64(player.Weight / 2)
	sRadius := float64(s.Weight / 2)
	dist := float64(utils.Distance(s.X, player.X, s.Y, player.Y))
	closeEnough := utils.IntersectionArea(plRadius, sRadius, dist)/utils.SurfaceArea64(sRadius) > 0.5
	bigEnough := plRadius > sRadius
	return closeEnough && bigEnough
}
