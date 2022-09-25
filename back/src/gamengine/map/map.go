package _map

import (
	"github.com/diyor28/not-agar/src/gamengine/constants"
	"github.com/diyor28/not-agar/src/gamengine/map/food"
	"github.com/diyor28/not-agar/src/gamengine/map/players"
	"github.com/diyor28/not-agar/src/gamengine/map/spikes"
	"github.com/diyor28/not-agar/src/gamengine/schemas"
	"github.com/diyor28/not-agar/src/randomname"
	"github.com/diyor28/not-agar/src/utils"
)

func New() *Map {
	gameMap := Map{
		Players: &players.Players{},
		Food:    &food.Foods{},
		Spikes:  &spikes.Spikes{},
	}
	return &gameMap
}

type Map struct {
	Players *players.Players
	Food    *food.Foods
	Spikes  *spikes.Spikes
}

func (m *Map) PopulateBots() {
	currentBots := m.Players.BotsCount()
	if currentBots >= constants.MaxBots {
		return
	}
	for i := 0; i < constants.MaxBots; i++ {
		m.CreatePlayer(randomname.GenerateNickname(), true)
	}
}

func (m *Map) PopulateSpikes() {
	for i := 0; i < constants.MaxSpikes; i++ {
		x, y := m.Spikes.RandXY(constants.SpikesSpacing)
		m.Spikes.New(x, y, constants.SpikeWeight)
	}
}

func (m *Map) PopulateFood() []*food.Food {
	var createdFood []*food.Food
	if m.Food.Len() < constants.MaxNumFood {
		for i := 0; i < constants.MaxNumFood-m.Food.Len(); i++ {
			createdFood = append(createdFood, m.createRandomFood())
		}
	}
	return createdFood
}

func (m *Map) GetStats() []*schemas.PlayerStat {
	plrs := m.Players.Largest(constants.StatsNumber)
	var topPlayers []*schemas.PlayerStat
	for i := 0; i < len(plrs); i++ {
		topPlayers = append(topPlayers, &schemas.PlayerStat{
			Nickname: plrs[i].Nickname,
			Weight:   int16(plrs[i].Weight),
		})
	}
	return topPlayers
}

func (m *Map) createFood(x float32, y float32) *food.Food {
	f := m.Food.New(x, y, constants.FoodWeight)
	return f
}

func (m *Map) createRandomFood() *food.Food {
	x, y := utils.RandXY()
	return m.createFood(x, y)
}

func (m *Map) CreatePlayer(nickname string, isBot bool) *players.Player {
	x, y := utils.RandXY()
	pl := m.Players.New(x, y, constants.MinWeight, nickname, isBot)
	return pl
}

func (m *Map) RemoveEatableFood() []*food.Food {
	var eatenFood []*food.Food
	for _, pl := range m.Players.Players {
		var filteredFoods []*food.Food
		for _, f := range m.Food.Food {
			if pl.FoodEatable(f) {
				pl.EatEntity(f)
				eatenFood = append(eatenFood, f)
			} else {
				filteredFoods = append(filteredFoods, f)
			}
		}
		m.Food.Food = filteredFoods
	}
	return eatenFood
}

func (m *Map) SpikeCollisions(pl *players.Player) bool {
	for _, s := range m.Spikes.Spikes {
		if s.Collided(pl) {
			return true
		}
	}
	return false
}

func (m *Map) RemoveDeadPlayers() []*players.Player {
	totalPlayers := len(m.Players.Players)
	for i := 0; i < totalPlayers; i++ {
		p1 := m.Players.Players[i]
		if m.SpikeCollisions(p1) {
			p1.IsDead = true
		}
		for k := i; k < totalPlayers; k++ {
			p2 := m.Players.Players[k]
			if p2.IsDead {
				continue
			}
			if p1.PlayerEatable(p2) {
				p1.EatEntity(p2)
				p2.IsDead = true
			}

			if p2.PlayerEatable(p1) {
				p2.EatEntity(p1)
				p1.IsDead = true
			}
		}
	}
	var newPlayers []*players.Player
	var eatenPlayers []*players.Player
	for index, pl := range m.Players.Players {
		if pl.IsDead {
			if m.Players.Players[index].IsBot {
				continue
			}
			eatenPlayers = append(eatenPlayers, m.Players.Players[index])
		} else {
			newPlayers = append(newPlayers, m.Players.Players[index])
		}
	}
	m.Players.Players = newPlayers
	return eatenPlayers
}
