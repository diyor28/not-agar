package gamengine

import (
	"encoding/json"
	"github.com/frankenbeanies/uuid4"
	"log"
	"math"
	"math/rand"
	"randomname"
	"sync"
	"time"
)

const windowSize = 800
const minSpeed = 1.1
const maxSpeed = 5
const maxX = 2000
const maxY = 2000
const maxNumFood = (maxX + maxY) / 10
const minWeight = 40
const maxWeight = 300
const minX = 0
const minY = 0
const minZoom = 0.8
const maxZoom = 1.0
const maxPlayers = 100
const maxSpikes = maxPlayers

type ServerResponse struct {
	SelfPlayer SelfPlayer `json:"selfPlayer"`
	Players    []Player   `json:"players"`
	Foods      []Food     `json:"foods"`
	Spikes     []Spike    `json:"spikes"`
}

type ServerRequest struct {
	Uuid       string `json:"uuid"`
	DirectionX int    `json:"directionX"`
	DirectionY int    `json:"directionY"`
}

var validDirection = map[int]bool{1: true, -1: true, 0: true}

func getRandomColor() string {
	return colors[rand.Intn(len(colors))]
}

func getRandomCoordinate() (float32, float32) {
	return float32(rand.Intn(maxX)), float32(rand.Intn(maxY))
}

func newWeight(weight1 float32, weight2 float32) float32 {
	areaSum := math.Pow(float64(weight1), 2) + math.Pow(float64(weight2), 2)
	return float32(math.Sqrt(areaSum))
}

func getSpeedFromWeight(weight float32) float32 {
	normalized := 1 - (weight-minWeight)/(maxWeight-minWeight)
	newSpeed := float64(normalized*(maxSpeed-minSpeed) + minSpeed)
	return float32(math.Max(newSpeed, minSpeed))
}

var colors = []string{"#ff1515", "#fff315", "#1557ff", "#15ffd0", "#ff15e0"}

type GameMap struct {
	GameId      string
	Players     []Player `json:"players"`
	Foods       []Food   `json:"foods"`
	Spikes      []Spike  `json:"spikes"`
	PlayersLock *sync.Mutex
	FoodsLock   *sync.Mutex
}

func (gMap *GameMap) populateBots() {
	//maxBots := (maxPlayers - len(gMap.Players)) / 2
	maxBots := 5
	currentBots := 0
	for _, player := range gMap.Players {
		if player.IsBot {
			currentBots++
		}
	}
	if currentBots < maxBots {
		for i := 0; i < maxBots; i++ {
			gMap.CreatePlayer(randomname.GenerateNickname(), true)
		}
	}
}

func (gMap *GameMap) populateSpikes() {
	for i := 0; i < maxSpikes; i++ {
		gMap.Spikes = append(gMap.Spikes, Spike{
			Uuid:  uuid4.New().String(),
			Color: getRandomColor(),
		})
	}
}

func (gMap *GameMap) populateFood() {
	if len(gMap.Foods) < maxNumFood {
		for i := 0; i < maxNumFood-len(gMap.Foods); i++ {
			gMap.createFood()
		}
	}
}

func (gMap *GameMap) Run() {
	for {
		time.Sleep(15 * time.Millisecond)
		var wg sync.WaitGroup
		gMap.populateFood()
		gMap.populateBots()
		gMap.PlayersLock.Lock()
		for i := range gMap.Players {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				player := &gMap.Players[i]
				if player.IsBot {
					player.makeMove(gMap)
				}
				player.updatePosition()
				gMap.eatFoods(player)
				gMap.eatPlayers(player)
			}(i)
		}
		wg.Wait()
		gMap.PlayersLock.Unlock()
	}
}

func (gMap *GameMap) ServerResponse(player Player) []byte {
	foods := gMap.nearByFood(player)
	players := gMap.nearByPlayers(player)
	serverResponse := ServerResponse{
		SelfPlayer: player.getSelfPlayer(),
		Foods:      foods,
		Players:    players,
	}
	responseMessage, err := json.Marshal(serverResponse)
	if err != nil {
		log.Println(err)
	}
	return responseMessage
}

func (gMap *GameMap) nearByFood(player Player) []Food {
	var foods []Food

	for _, f := range gMap.Foods {
		distX := math.Abs(float64(player.X - f.X))
		distY := math.Abs(float64(player.Y - f.Y))
		if distX < windowSize && distY < windowSize {
			foods = append(foods, f)
		}
	}

	return foods
}

func (gMap *GameMap) nearByPlayers(player Player) []Player {
	var players []Player

	for _, p := range gMap.Players {
		if p.Uuid == player.Uuid {
			continue
		}
		distX := math.Abs(float64(player.X - p.X))
		distY := math.Abs(float64(player.Y - p.Y))
		if distX < windowSize && distY < windowSize {
			players = append(players, p)
		}
	}

	return players
}

func (gMap *GameMap) removePlayerIndex(index int) {
	gMap.PlayersLock.Lock()
	gMap.Players = append(gMap.Players[:index], gMap.Players[index+1:]...)
	gMap.PlayersLock.Unlock()
}

func (gMap *GameMap) removeFoodIndex(index int) {
	gMap.FoodsLock.Lock()
	gMap.Foods = append(gMap.Foods[:index], gMap.Foods[index+1:]...)
	gMap.FoodsLock.Unlock()
}
func (gMap *GameMap) removeFoodUuid(uuid string) {
	for i, p := range gMap.Foods {
		if p.Uuid == uuid {
			gMap.FoodsLock.Lock()
			gMap.Foods = append(gMap.Foods[:i], gMap.Foods[i+1:]...)
			gMap.FoodsLock.Unlock()
		}
	}
}

func (gMap *GameMap) RemovePlayerUUID(uuid string) {
	for i, p := range gMap.Players {
		if p.Uuid == uuid {
			gMap.PlayersLock.Lock()
			gMap.Players = append(gMap.Players[:i], gMap.Players[i+1:]...)
			gMap.PlayersLock.Unlock()
		}
	}
}

func (gMap *GameMap) createFood() Food {
	x, y := getRandomCoordinate()
	food := Food{
		Uuid:   uuid4.New().String(),
		X:      x,
		Y:      y,
		Color:  getRandomColor(),
		Weight: 25,
	}
	gMap.FoodsLock.Lock()
	gMap.Foods = append(gMap.Foods, food)
	gMap.FoodsLock.Unlock()
	return food
}

func (gMap *GameMap) CreatePlayer(nickname string, isBot bool) SelfPlayer {
	x, y := getRandomCoordinate()
	player := Player{
		Uuid:     uuid4.New().String(),
		X:        x,
		Y:        y,
		Color:    getRandomColor(),
		Weight:   minWeight,
		Speed:    maxSpeed,
		Zoom:     1,
		Nickname: nickname,
		IsBot:    isBot,
	}
	if len(gMap.Players) >= maxPlayers {
		gMap.removePlayerIndex(0)
	}
	gMap.PlayersLock.Lock()
	gMap.Players = append(gMap.Players, player)
	gMap.PlayersLock.Unlock()
	return player.getSelfPlayer()
}

func (gMap *GameMap) GetPlayer(uuid string) Player {
	for _, p := range gMap.Players {
		if p.Uuid == uuid {
			return p
		}
	}
	return Player{}
}

func (gMap *GameMap) eatableFood(player Player) []Food {
	var eatable []Food
	for _, f := range gMap.nearByFood(player) {
		if player.foodEatable(f) {
			eatable = append(eatable, f)
		}
	}
	return eatable
}

func (gMap *GameMap) eatablePlayers(player Player) []Player {
	var eatable []Player
	for _, p := range gMap.nearByPlayers(player) {
		if player.playerEatable(p) {
			eatable = append(eatable, p)
		}
	}
	return eatable
}

func (gMap *GameMap) eatFood(player *Player, food Food) {
	player.Weight = newWeight(player.Weight, food.Weight)
	player.Speed = getSpeedFromWeight(player.Weight)
	gMap.removeFoodUuid(food.Uuid)
}

func (gMap *GameMap) eatPlayer(player1 *Player, player2 *Player) {
	var removePlayerUUID string
	if player1.Weight > player2.Weight {
		player1.Weight = newWeight(player1.Weight, player2.Weight)
		player1.Speed = getSpeedFromWeight(player1.Weight)
		removePlayerUUID = player2.Uuid
	} else {
		player2.Weight = newWeight(player1.Weight, player2.Weight)
		player2.Speed = getSpeedFromWeight(player2.Weight)
		removePlayerUUID = player1.Uuid
	}
	gMap.RemovePlayerUUID(removePlayerUUID)
}

func (gMap *GameMap) eatFoods(player *Player) {
	foods := gMap.eatableFood(*player)
	for _, f := range foods {
		gMap.eatFood(player, f)
	}
}

func (gMap *GameMap) eatPlayers(player *Player) {
	players := gMap.eatablePlayers(*player)
	for _, p := range players {
		gMap.eatPlayer(player, &p)
	}
}

func (gMap *GameMap) UpdatePlayer(request ServerRequest) *Player {
	for i, p := range gMap.Players {
		if p.Uuid == request.Uuid {
			player := &gMap.Players[i]
			player.updateDirection(request.DirectionX, request.DirectionY)
			return player
		}
	}
	return &Player{}
}
