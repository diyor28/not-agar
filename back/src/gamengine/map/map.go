package _map

import (
	"github.com/diyor28/not-agar/src/gamengine/constants"
	"github.com/diyor28/not-agar/src/gamengine/map/food"
	"github.com/diyor28/not-agar/src/gamengine/map/players"
	"github.com/diyor28/not-agar/src/gamengine/map/spikes"
	"github.com/diyor28/not-agar/src/randomname"
	"github.com/diyor28/not-agar/src/utils"
)

type PlayerStat struct {
	Nickname string
	Weight   int16
}

func New() *Map {
	gameMap := Map{}
	return &gameMap
}

type Map struct {
	Players players.Players `json:"players"`
	Foods   food.Foods      `json:"foods"`
	Spikes  spikes.Spikes   `json:"spikes"`
}

func (m *Map) PopulateBots() {
	//maxBots := (MaxPlayers - len(m.Players)) / 2
	maxBots := 0
	currentBots := m.Players.BotsCount()
	if currentBots < maxBots {
		for i := 0; i < maxBots; i++ {
			m.CreatePlayer(randomname.GenerateNickname(), true)
		}
	}
}

func (m *Map) PopulateSpikes() {
	for i := 0; i < constants.MaxSpikes; i++ {
		x, y := m.Spikes.RandXY(constants.SpikesSpacing)
		spike := spikes.New(x, y, constants.SpikeWeight)
		m.Spikes = append(m.Spikes, spike)
	}
}

func (m *Map) PopulateFood() {
	if len(m.Foods) < constants.MaxNumFood {
		for i := 0; i < constants.MaxNumFood-len(m.Foods); i++ {
			m.createRandomFood()
		}
	}
}

func (m *Map) GetStats() []PlayerStat {
	plrs := m.Players.Largest(constants.StatsNumber)
	var topPlayers []PlayerStat
	for i := 0; i < len(plrs); i++ {
		topPlayers = append(topPlayers, PlayerStat{
			Nickname: plrs[i].Nickname,
			Weight:   int16(plrs[i].Weight),
		})
	}
	return topPlayers
}

func (m *Map) removePlayerIndex(index int) {
	m.Players = append(m.Players[:index], m.Players[index+1:]...)
}

func (m *Map) removeFoodIndex(index int) {
	m.Foods = append(m.Foods[:index], m.Foods[index+1:]...)
}
func (m *Map) removeFoodUuid(uuid string) {
	for i, p := range m.Foods {
		if p.Uuid == uuid {
			m.Foods = append(m.Foods[:i], m.Foods[i+1:]...)
		}
	}
}

func (m *Map) RemovePlayerUUID(uuid string) {
	for i, p := range m.Players {
		if p.Uuid == uuid {
			m.Players = append(m.Players[:i], m.Players[i+1:]...)
		}
	}
}

func (m *Map) createFood(x float32, y float32) *food.Food {
	food := food.New(x, y, constants.FoodWeight)
	m.Foods = append(m.Foods, food)
	return food
}

func (m *Map) createRandomFood() *food.Food {
	x, y := utils.RandXY()
	return m.createFood(x, y)
}

func (m *Map) CreatePlayer(nickname string, isBot bool) (players.SelfPlayer, spikes.Spikes) {
	x, y := utils.RandXY()
	pl := players.NewPlayer(x, y, constants.MinWeight, nickname, isBot)
	// TODO: better solution
	//if len(m.Players) >= constants.MaxPlayers {
	//	m.removePlayerIndex(0)
	//}
	m.Players = append(m.Players, pl)
	return pl.GetSelfPlayer(), m.Spikes
}

func (m *Map) RemoveEatableFood() {
	for _, pl := range m.Players {
		var filteredFoods []*food.Food
		for _, f := range m.Foods {
			if pl.FoodEatable(f) {
				pl.EatEntity(f)
			} else {
				filteredFoods = append(filteredFoods, f)
			}
		}
		m.Foods = filteredFoods
	}
}

func (m *Map) SpikeCollisions(pl *players.Player) bool {
	for _, s := range m.Spikes {
		if s.Collided(pl) {
			return true
		}
	}
	return false
}
