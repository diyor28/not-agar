package gamengine

import (
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
	SurfaceArea        = MaxXY * MaxXY
	FoodWeight         = 25
	SpikeWeight        = 120
	MinSpeed           = 2
	MaxSpeed           = 8
	MaxNumFood         = SurfaceArea / 50000
	MinWeight          = 40
	MaxWeight          = MinWeight * 25
	MinZoom            = 0.7
	MaxZoom            = 1.0
	MaxPlayers         = 20
	MaxSpikes          = MaxPlayers * 5
	StatsNumber        = 10
	SpeedWeightLimit   = 500
	NumFoodResponse    = 30
	NumPlayersResponse = 10
	SpikesSpacing      = MaxXY / MaxSpikes
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

func randomColor() [3]int {
	return colors[rand.Intn(len(colors))]
}

func randXYNoOverlap(data interface{}, minDistance float32) (float32, float32) {
	x, y := randXY()
	entities := data.([]interface {
		getWeight() float32
		getXY() (float32, float32)
	})
	for _, entity := range entities {
		radius := entity.getWeight() / 2
		eX, eY := entity.getXY()
		dist := utils.CalcDistance(eX, x, eY, y)
		if dist-radius < minDistance {
			return randXYNoOverlap(data, minDistance)
		}
	}
	return x, y
}

func randXY() (float32, float32) {
	return float32(rand.Intn(MaxXY)), float32(rand.Intn(MaxXY))
}

func getSpeedFromWeight(weight float32) float32 {
	normalized := 1 - (weight-MinWeight)/(SpeedWeightLimit-MinWeight)
	newSpeed := float64(normalized*(MaxSpeed-MinSpeed) + MinSpeed)
	return float32(math.Max(newSpeed, MinSpeed))
}

var colors = [][3]int{{255, 21, 21}, {255, 243, 21}, {21, 87, 255}, {21, 255, 208}, {255, 21, 224}}

type GameMap struct {
	GameId  string
	Players Players  `json:"players"`
	Foods   []*Food  `json:"foods"`
	Spikes  []*Spike `json:"spikes"`
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
	player, err := gMap.Players.update(moveData.Uuid, moveData.NewX, moveData.NewY)
	if err != nil {
		log.Println(err, moveData.Uuid)
		return
	}
	if err := client.Emit("moved", MovedEvent{player.X, player.Y, player.Weight, player.Zoom}); err != nil {
		log.Println("Socket emit", err)
	}
	players := gMap.Players.closest(player, NumPlayersResponse)
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
	player, err := gMap.Players.get(accelerateData.Uuid)
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

func (gMap *GameMap) getBotsCount() int {
	currentBots := 0
	for _, player := range gMap.Players {
		if player.IsBot {
			currentBots++
		}
	}
	return currentBots
}

func (gMap *GameMap) populateBots() {
	maxBots := (MaxPlayers - len(gMap.Players)) / 2
	//maxBots := 0
	currentBots := gMap.getBotsCount()
	if currentBots < maxBots {
		for i := 0; i < maxBots; i++ {
			gMap.CreatePlayer(randomname.GenerateNickname(), true)
		}
	}
}

func (gMap *GameMap) populateSpikes() {
	for i := 0; i < MaxSpikes; i++ {
		var spikes []Spike
		for _, s := range gMap.Spikes {
			spikes = append(spikes, *s)
		}
		//x, y := randXYNoOverlap(spikes, SpikesSpacing)
		x, y := randXY()
		spike := NewSpike(x, y, SpikeWeight)
		gMap.Spikes = append(gMap.Spikes, spike)
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
	//density := int(20 * totalWeight / SurfaceArea)
	//fmt.Println("DENSITY", 20 * totalWeight / SurfaceArea)
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
		gMap.notifyAllPlayers("stats", stats)
	}
}

func (gMap *GameMap) notifyAllPlayers(event string, data interface{}) {
	for _, p := range gMap.Players {
		if p.IsBot {
			continue
		}
		gMap.Hub.Emit(event, data, p.Uuid)
	}
}

func (gMap *GameMap) publishAdminStats() {
	for {
		var result = make(map[string]interface{})
		botsCount := gMap.getBotsCount()
		stats := gMap.GetStats()

		result["botsCount"] = botsCount
		result["playersCount"] = len(gMap.Players) - botsCount
		result["topPlayers"] = stats
		gMap.Hub.Emit("stats", result, "admin")
		time.Sleep(2 * time.Second)
	}
}

func (gMap *GameMap) Run() {
	counter := 0
	go gMap.publishStats()
	go gMap.Hub.Run()
	go gMap.populateSpikes()
	go gMap.publishAdminStats()
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

func (gMap GameMap) GetStats() []map[string]interface{} {
	players := gMap.Players.largest(StatsNumber)
	var topPlayers []map[string]interface{}
	for i := 0; i < len(players); i++ {
		stat := map[string]interface{}{}
		stat["nickname"] = players[i].Nickname
		stat["weight"] = players[i].Weight
		topPlayers = append(topPlayers, stat)
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

func (gMap *GameMap) createFood(x float32, y float32) *Food {
	food := NewFood(x, y, FoodWeight)
	gMap.Foods = append(gMap.Foods, food)
	return food
}

func (gMap *GameMap) createRandomFood() *Food {
	x, y := randXY()
	return gMap.createFood(x, y)
}

func (gMap *GameMap) CreatePlayer(nickname string, isBot bool) SelfPlayer {
	x, y := randXY()
	player := NewPlayer(x, y, MinWeight, nickname, isBot)
	if len(gMap.Players) >= MaxPlayers {
		gMap.removePlayerIndex(0)
	}
	gMap.Players = append(gMap.Players, player)
	return player.getSelfPlayer()
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
