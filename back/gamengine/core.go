package gamengine

import (
	"github.com/diyor28/not-agar/randomname"
	"github.com/diyor28/not-agar/sockethub"
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
	FoodWeight         = 20
	SpikeWeight        = 120
	MinSpeed           = 2
	MaxSpeed           = 5
	MaxNumFood         = SurfaceArea / 50000
	MinWeight          = 40
	MaxWeight          = MinWeight * 25
	MinZoom            = 0.7
	MaxZoom            = 1.0
	MaxPlayers         = 20
	MaxSpikes          = MaxPlayers * 2
	StatsNumber        = 10
	SpeedWeightLimit   = 500
	NumFoodResponse    = 30
	NumPlayersResponse = 10
	SpikesSpacing      = MaxXY/MaxSpikes + SpikeWeight
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
	Players Players `json:"players"`
	Foods   Foods   `json:"foods"`
	Spikes  Spikes  `json:"spikes"`
	Hub     *sockethub.Hub
}

func NewGameMap() *GameMap {
	hub := sockethub.NewHub()
	gameMap := GameMap{
		GameId: uuid4.New().String(),
		Hub:    hub,
	}
	hub.On("ping", gameMap.sendPong)
	hub.On("shoot", gameMap.handleShoot)
	hub.On("move", gameMap.handleMoveEvent)
	return &gameMap
}

func (gMap *GameMap) handleMoveEvent(data interface{}, playerId string) {
	var moveData MoveEvent
	if err := mapstructure.Decode(data, &moveData); err != nil {
		log.Println(err)
		return
	}
	player, err := gMap.Players.update(playerId, moveData.NewX, moveData.NewY)
	if err != nil {
		log.Println(err, moveData.Uuid)
		return
	}
	gMap.Hub.Emit("moved", MovedEvent{player.X, player.Y, player.Weight, player.Zoom}, player.Uuid)

	//binaryBuffer := make([]byte, 256)
	//binary.BigEndian.PutUint16(binaryBuffer, )
	//
	//if err := client.Emit("foodUpdatedB", foods); err != nil {
	//	log.Println("socket emit: ", err)
	//}
}

func (gMap *GameMap) handleShoot(data interface{}, playerId string) {
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

func (gMap *GameMap) sendPong(data interface{}, playerId string) {
	gMap.Hub.Emit("pong", data, playerId)
}

func (gMap *GameMap) populateBots() {
	//maxBots := (MaxPlayers - len(gMap.Players)) / 2
	maxBots := 0
	currentBots := gMap.Players.botsCount()
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
		x, y := gMap.Spikes.randXY(SpikesSpacing)
		spike := NewSpike(x, y, SpikeWeight)
		gMap.Spikes = append(gMap.Spikes, spike)
	}
}

func (gMap *GameMap) populateFood() {
	if len(gMap.Foods) < MaxNumFood {
		for i := 0; i < MaxNumFood-len(gMap.Foods); i++ {
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
	for _, p := range gMap.Players.real() {
		gMap.Hub.Emit(event, data, p.Uuid)
	}
}

func (gMap *GameMap) publishAdminStats() {
	for {
		var result = make(map[string]interface{})
		botsCount := gMap.Players.botsCount()
		stats := gMap.GetStats()

		result["botsCount"] = botsCount
		result["playersCount"] = len(gMap.Players) - botsCount
		result["topPlayers"] = stats
		gMap.Hub.Emit("stats", result, "admin")
		time.Sleep(2 * time.Second)
	}
}

func (gMap *GameMap) notifyPlayer(player *Player) {
	gMap.Hub.Emit("moved", MovedEvent{player.X, player.Y, player.Weight, player.Zoom}, player.Uuid)
	players := gMap.Players.closest(player, NumPlayersResponse)
	gMap.Hub.Emit("playersUpdated", players, player.Uuid)
	foods := gMap.Foods.closest(player, NumFoodResponse)
	gMap.Hub.Emit("foodUpdated", foods, player.Uuid)
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
				if player.IsBot {
					return
				}
				if counter%2 == 0 {
					gMap.notifyPlayer(player)
				}
			}(i)
		}
		wg.Wait()
		gMap.removeEatableFood()
		gMap.removeDeadPlayers()
		time.Sleep(15 * time.Millisecond)
		//fmt.Println("running loop")
	}
}

func (gMap GameMap) GetStats() []map[string]interface{} {
	players := gMap.Players.largest(StatsNumber)
	var topPlayers []map[string]interface{}
	for i := 0; i < len(players); i++ {
		stat := map[string]interface{}{}
		stat["nickname"] = players[i].Nickname
		stat["weight"] = int(players[i].Weight)
		topPlayers = append(topPlayers, stat)
	}
	return topPlayers
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

func (gMap *GameMap) CreatePlayer(nickname string, isBot bool) map[string]interface{} {
	x, y := randXY()
	player := NewPlayer(x, y, MinWeight, nickname, isBot)
	if len(gMap.Players) >= MaxPlayers {
		gMap.removePlayerIndex(0)
	}
	gMap.Players = append(gMap.Players, player)
	var result = make(map[string]interface{})
	result["player"] = player.getSelfPlayer()
	result["spikes"] = gMap.Spikes
	return result
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

func (gMap *GameMap) removeDeadPlayers() {
	// take first player, compare it to every other player after it
	// get a new array of players that
	//for _, p := range eatenPlayers {
	//	gMap.Hub.Emit("rip", "", p.Uuid)
	//}
	//gMap.Players = newPlayers
}
