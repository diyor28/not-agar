package gamengine

import (
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/randomname"
	"github.com/diyor28/not-agar/sockethub"
	"github.com/diyor28/not-agar/utils"
	"github.com/frankenbeanies/uuid4"
	"github.com/mitchellh/mapstructure"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	MinXY              = 0
	MaxXY              = 5000
	surfaceArea        = MaxXY * MaxXY
	FoodWeight         = 25
	SpikeWeight        = 80
	MinSpeed           = 2
	MaxSpeed           = 8
	MaxNumFood         = surfaceArea / 50000
	MinWeight          = 40
	MaxWeight          = MinWeight * 15
	MinZoom            = 0.8
	MaxZoom            = 1.0
	maxPlayers         = 20
	MaxSpikes          = maxPlayers
	StatsNumber        = 10
	SpeedWeightLimit   = 300
	NumFoodResponse    = 30
	NumPlayersResponse = 10
)

type MoveEvent struct {
	Uuid string  `json:"uuid"`
	NewX float32 `json:"newX"`
	NewY float32 `json:"newY"`
}

type AccelerateEvent struct {
	Uuid string `json:"uuid"`
}

type MovedEvent struct {
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Weight float32 `json:"weight"`
	Zoom   float32 `json:"zoom"`
}

func getRandomColor() [3]int {
	return colors[rand.Intn(len(colors))]
}

func getRandomCoordinate() (float32, float32) {
	return float32(rand.Intn(MaxXY)), float32(rand.Intn(MaxXY))
}

func newWeight(weight1 float32, weight2 float32) float32 {
	nWeight := math.Sqrt(float64(weight1*weight1 + weight2*weight2))
	return float32(math.Min(nWeight, MaxWeight))
}

func getSpeedFromWeight(weight float32) float32 {
	normalized := 1 - (weight-MinWeight)/(SpeedWeightLimit-MinWeight)
	newSpeed := float64(normalized*(MaxSpeed-MinSpeed) + MinSpeed)
	return float32(math.Max(newSpeed, MinSpeed))
}

var colors = [][3]int{{255, 21, 21}, {255, 243, 21}, {21, 87, 255}, {21, 255, 208}, {255, 21, 224}}

type GameMap struct {
	GameId  string
	Players []*Player `json:"players"`
	Foods   []*Food   `json:"foods"`
	Spikes  []*Spike  `json:"spikes"`
	Hub     *sockethub.Hub
}

func NewGameMap() *GameMap {
	hub := sockethub.NewHub()
	gameMap := GameMap{
		GameId: uuid4.New().String(),
		Hub:    hub,
	}
	hub.On("ping", gameMap.sendPong)
	hub.On("accelerate", gameMap.handleAccelerate)
	hub.On("move", gameMap.handleMoveEvent)
	return &gameMap
}

func (gMap *GameMap) handleMoveEvent(data interface{}, client *sockethub.Client) {
	var moveData MoveEvent
	if err := mapstructure.Decode(data, &moveData); err != nil {
		log.Println(err)
		return
	}
	player, err := gMap.UpdatePlayer(moveData.Uuid, moveData.NewX, moveData.NewY)
	if err != nil {
		log.Println(err, moveData.Uuid)
		return
	}
	if err := client.Emit("moved", MovedEvent{player.X, player.Y, player.Weight, player.Zoom}); err != nil {
		log.Println("Socket emit", err)
	}
	players := gMap.nearByPlayers(player, NumPlayersResponse)
	if err := client.Emit("playersUpdated", players); err != nil {
		log.Println("Socket emit: ", err)
	}
	foods := gMap.nearByFood(player, NumFoodResponse)
	if err := client.Emit("foodUpdated", foods); err != nil {
		log.Println("Socket emit: ", err)
	}

	spikes := gMap.nearBySpikes(player, NumFoodResponse)
	if err := client.Emit("spikesUpdated", spikes); err != nil {
		log.Println("Socket emit: ", err)
	}

	//binaryBuffer := make([]byte, 256)
	//binary.BigEndian.PutUint16(binaryBuffer, )
	//
	//if err := client.Emit("foodUpdatedB", foods); err != nil {
	//	log.Println("Socket emit: ", err)
	//}
}

func (gMap *GameMap) handleAccelerate(data interface{}, client *sockethub.Client) {
	var accelerateData AccelerateEvent
	if err := mapstructure.Decode(data, &accelerateData); err != nil {
		log.Println(err)
		return
	}
	player, err := gMap.GetPlayer(accelerateData.Uuid)
	if err != nil {
		log.Println("could not find player ", err)
		return
	}
	if player.Weight > 2*MinWeight {
		player.Accelerating = true
	}
}

func (gMap *GameMap) sendPong(data interface{}, client *sockethub.Client) {
	if err := client.Emit("pong", data); err != nil {
		log.Println(err)
	}
}

func (gMap *GameMap) populateBots() {
	//maxBots := (maxPlayers - len(gMap.Players)) / 2
	maxBots := 0
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
	for i := 0; i < MaxSpikes; i++ {
		x, y := getRandomCoordinate()
		spike := Spike{
			Uuid:   uuid4.New().String(),
			Weight: SpikeWeight,
			X:      x,
			Y:      y,
		}
		gMap.Spikes = append(gMap.Spikes, &spike)
	}
}

func (gMap *GameMap) populateFood() {
	var totalWeight float32 = 0
	maxFood := MaxNumFood
	for _, p := range gMap.Players {
		totalWeight += p.Weight * p.Weight * math.Pi
	}
	for _, f := range gMap.Foods {
		totalWeight += f.Weight * f.Weight * math.Pi
	}
	//density := int(20 * totalWeight / surfaceArea)
	//fmt.Println("DENSITY", 20 * totalWeight / surfaceArea)
	//maxFood -= density
	if len(gMap.Foods) < maxFood {
		for i := 0; i < maxFood-len(gMap.Foods); i++ {
			gMap.createRandomFood()
		}
	}
}

func (gMap *GameMap) publishStats() {
	for {
		time.Sleep(2 * time.Second)
		stats := gMap.GetStats()
		gMap.Hub.Publish("stats", stats)
	}
}

func (gMap *GameMap) Run() {
	counter := 0
	go gMap.publishStats()
	go gMap.Hub.Run()
	go gMap.populateSpikes()
	for {
		counter++
		var wg sync.WaitGroup
		if counter > 30 {
			counter = 0
			gMap.populateFood()
			gMap.populateBots()
		}
		for i := range gMap.Players {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				player := gMap.Players[i]
				if player.IsBot {
					player.makeMove(gMap)
				}
				player.updatePosition(gMap)
				player.passiveWeightLoss()
			}(i)
		}
		wg.Wait()
		gMap.removeEatableFood()
		gMap.removeEatablePlayers()
		time.Sleep(15 * time.Millisecond)
	}
}

func (gMap *GameMap) nearByFood(player *Player, kClosest int) []Food {
	totalNumFood := len(gMap.Foods)
	var foods = make([]Food, 0)
	for _, f := range gMap.Foods {
		foods = append(foods, *f)
	}
	foodDistances := make(map[string]float32, totalNumFood)
	for _, f := range foods {
		foodDistances[f.Uuid] = utils.CalcDistance(player.X, f.X, player.Y, f.Y)
	}
	numResults := kClosest
	if kClosest > totalNumFood {
		numResults = totalNumFood
	}
	for i := 0; i < numResults; i++ {
		var minIdx = i
		for j := i + 1; j < totalNumFood; j++ {
			if foodDistances[foods[j].Uuid] < foodDistances[foods[minIdx].Uuid] {
				minIdx = j
			}
		}
		foods[i], foods[minIdx] = foods[minIdx], foods[i]
	}
	return foods[:numResults]
}

func (gMap *GameMap) nearBySpikes(player *Player, kClosest int) []Spike {
	totalSpikes := len(gMap.Spikes)
	var spikes []Spike
	for _, s := range gMap.Spikes {
		spikes = append(spikes, *s)
	}
	spikeDistances := make(map[string]float32, totalSpikes)
	for _, s := range spikes {
		spikeDistances[s.Uuid] = utils.CalcDistance(player.X, s.X, player.Y, s.Y)
	}
	numResults := kClosest
	if kClosest > totalSpikes {
		numResults = totalSpikes
	}
	for i := 0; i < numResults; i++ {
		var minIdx = i
		for j := i + 1; j < totalSpikes; j++ {
			if spikeDistances[spikes[j].Uuid] < spikeDistances[spikes[minIdx].Uuid] {
				minIdx = j
			}
		}
		spikes[i], spikes[minIdx] = spikes[minIdx], spikes[i]
	}
	return spikes[:numResults]
}

type StatsResponse struct {
	Nickname string `json:"nickname"`
	Weight   int    `json:"weight"`
}

func (gMap GameMap) GetStats() []StatsResponse {
	players := gMap.Players
	numOfPlayers := len(players)
	playersInStats := StatsNumber
	if numOfPlayers < StatsNumber {
		playersInStats = numOfPlayers
	}
	var topPlayers []StatsResponse
	for i := 0; i < playersInStats; i++ {
		var maxIdx = i
		for j := i; j < numOfPlayers; j++ {
			if players[j].Weight > players[maxIdx].Weight {
				maxIdx = j
			}
		}
		players[i], players[maxIdx] = players[maxIdx], players[i]
	}
	for i := 0; i < playersInStats; i++ {
		topPlayers = append(topPlayers, StatsResponse{players[i].Nickname, int(players[i].Weight)})
	}
	return topPlayers
}

func (gMap *GameMap) playersExcept(uuid string) []Player {
	var players = make([]Player, 0)
	for _, p := range gMap.Players {
		if p.Uuid != uuid {
			players = append(players, *p)
		}
	}
	return players
}

func (gMap *GameMap) nearByPlayers(player *Player, kClosest int) []Player {
	players := gMap.playersExcept(player.Uuid)
	totalPlayers := len(players)
	playersDistances := make(map[string]float32, totalPlayers)
	for _, p := range players {
		playersDistances[p.Uuid] = utils.CalcDistance(player.X, p.X, player.Y, p.Y)
	}

	numResults := kClosest
	if kClosest > totalPlayers {
		numResults = totalPlayers
	}
	for i := 0; i < numResults; i++ {
		var minIdx = i
		for j := i + 1; j < totalPlayers; j++ {
			if playersDistances[players[j].Uuid] < playersDistances[players[minIdx].Uuid] {
				minIdx = j
			}
		}
		players[i], players[minIdx] = players[minIdx], players[i]
	}
	return players[:numResults]
}

func (gMap *GameMap) removePlayerIndex(index int) {
	gMap.Players = append(gMap.Players[:index], gMap.Players[index+1:]...)
}

func (gMap *GameMap) removeFoodIndex(index int) {
	gMap.Foods = append(gMap.Foods[:index], gMap.Foods[index+1:]...)
}
func (gMap *GameMap) removeFoodUuid(uuid string) {
	for i, p := range gMap.Foods {
		if p.Uuid == uuid {
			gMap.Foods = append(gMap.Foods[:i], gMap.Foods[i+1:]...)
		}
	}
}

func (gMap *GameMap) RemovePlayerUUID(uuid string) {
	for i, p := range gMap.Players {
		if p.Uuid == uuid {
			gMap.Players = append(gMap.Players[:i], gMap.Players[i+1:]...)
		}
	}
}

func (gMap *GameMap) createFood(x float32, y float32) Food {
	food := Food{
		Uuid:   uuid4.New().String(),
		X:      x,
		Y:      y,
		Color:  getRandomColor(),
		Weight: FoodWeight,
	}
	gMap.Foods = append(gMap.Foods, &food)
	return food
}

func (gMap *GameMap) createRandomFood() Food {
	x, y := getRandomCoordinate()
	return gMap.createFood(x, y)
}

func (gMap *GameMap) GetPlayer(uuid string) (*Player, error) {
	for i := range gMap.Players {
		if gMap.Players[i].Uuid == uuid {
			return gMap.Players[i], nil
		}
	}
	return nil, errors.New("no player found")
}

func (gMap *GameMap) CreatePlayer(nickname string, isBot bool) SelfPlayer {
	x, y := getRandomCoordinate()
	player := Player{
		Uuid:         uuid4.New().String(),
		X:            x,
		Y:            y,
		Color:        getRandomColor(),
		Weight:       MinWeight,
		Accelerating: false,
		Speed:        MaxSpeed,
		Zoom:         1,
		Nickname:     nickname,
		IsBot:        isBot,
	}
	if len(gMap.Players) >= maxPlayers {
		gMap.removePlayerIndex(0)
	}
	gMap.Players = append(gMap.Players, &player)
	return player.getSelfPlayer()
}

func (gMap *GameMap) UpdatePlayer(uuid string, newX float32, newY float32) (*Player, error) {
	player, err := gMap.GetPlayer(uuid)
	if err != nil {
		return nil, err
	}
	player.updateDirection(newX, newY)
	return player, nil
}

func (gMap *GameMap) removeEatableFood() {
	for _, player := range gMap.Players {
		var filteredFoods []*Food
		for _, food := range gMap.Foods {
			if player.foodEatable(food) {
				player.eatEntity(food)
			} else {
				filteredFoods = append(filteredFoods, food)
			}
		}
		gMap.Foods = filteredFoods
	}
}

func (gMap *GameMap) spikeCollisions(pl *Player) bool {
	for _, s := range gMap.Spikes {
		if s.collided(pl) {
			return true
		}
	}
	return false
}

func (gMap *GameMap) removeEatablePlayers() {
	// take first player, compare it to every other player after it
	// get a new array of players that
	//fmt.Println("Removing eatable players")
	totalPlayers := len(gMap.Players)
	var willBeEaten = make(map[int]bool, totalPlayers)
	for i := 0; i < totalPlayers; i++ {
		willBeEaten[i] = false
	}

	for i := 0; i < totalPlayers; i++ {
		p1 := gMap.Players[i]
		fmt.Println("spikeCollisions", p1)
		if gMap.spikeCollisions(p1) {
			willBeEaten[i] = true
		}
		for k := i; k < totalPlayers; k++ {
			if willBeEaten[k] {
				continue
			}
			p2 := gMap.Players[k]
			if p1.canEat(p2) {
				p1.eatEntity(p2)
				willBeEaten[k] = true
			}

			if p2.canEat(p1) {
				p2.eatEntity(p1)
				willBeEaten[i] = true
			}
		}
	}
	var newPlayers []*Player
	var eatenPlayers []Player
	for index, value := range willBeEaten {
		if value {
			if gMap.Players[index].IsBot {
				continue
			}
			eatenPlayers = append(eatenPlayers, *gMap.Players[index])
		} else {
			newPlayers = append(newPlayers, gMap.Players[index])
		}
	}

	for _, p := range eatenPlayers {
		gMap.Hub.Emit("rip", "", p.Uuid)
	}
	gMap.Players = newPlayers
}
