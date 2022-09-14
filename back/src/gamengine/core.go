package gamengine

import (
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/gamengine/constants"
	"github.com/diyor28/not-agar/src/gamengine/food"
	"github.com/diyor28/not-agar/src/gamengine/player"
	"github.com/diyor28/not-agar/src/randomname"
	"github.com/diyor28/not-agar/src/sockethub"
	"github.com/frankenbeanies/uuid4"
	"log"
	"math/rand"
	"sync"
	"time"
)

type AccelerateEvent struct {
	Uuid string `json:"uuid"`
}

func randXY() (float32, float32) {
	return float32(rand.Intn(constants.MaxXY)), float32(rand.Intn(constants.MaxXY))
}

type GameMap struct {
	GameId     string
	Players    player.Players `json:"players"`
	Foods      food.Foods     `json:"foods"`
	Spikes     Spikes         `json:"spikes"`
	Hub        *sockethub.Hub
	PlayersMap map[*sockethub.Client]string
}

func NewGameMap() *GameMap {
	hub := sockethub.NewHub()
	gameMap := GameMap{
		GameId: uuid4.New().String(),
		Hub:    hub,
	}
	hub.OnMessage(func(data []byte, client *sockethub.Client) {
		var event *GenericEvent
		if err := GenericSchema.Decode(data, event); err != nil {
			log.Println(err)
			return
		}
		switch event.Event {
		case "move":
			var moveEvent *MoveEvent
			if err := MovedSchema.Decode(data, moveEvent); err != nil {
				log.Println(err)
				return
			}
			gameMap.HandleMoveEvent(moveEvent, client)
		case "ping":
			gameMap.SendPong(data, client)
		}
	})
	return &gameMap
}

func (gMap *GameMap) PlayerReverseLookUp(uuid string) (*sockethub.Client, error) {
	for k, pUuid := range gMap.PlayersMap {
		if pUuid == uuid {
			return k, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no player for uuid %s found", uuid))
}

func (gMap *GameMap) HandleMoveEvent(data *MoveEvent, client *sockethub.Client) {
	pl, err := gMap.Players.Update(gMap.PlayersMap[client], data.NewX, data.NewY)
	if err != nil {
		log.Println(err)
		return
	}
	bytes, err := MovedSchema.Encode(MovedEvent{"moved", pl.X, pl.Y, pl.Weight, pl.Zoom})
	if err != nil {
		log.Println(err)
		return
	}
	client.Emit(bytes)
}

func (gMap *GameMap) SendPong(data []byte, client *sockethub.Client) {
	client.Emit(data)
}

func (gMap *GameMap) PopulateBots() {
	//maxBots := (MaxPlayers - len(gMap.Players)) / 2
	maxBots := 0
	currentBots := gMap.Players.BotsCount()
	if currentBots < maxBots {
		for i := 0; i < maxBots; i++ {
			gMap.CreatePlayer(randomname.GenerateNickname(), true)
		}
	}
}

func (gMap *GameMap) PopulateSpikes() {
	for i := 0; i < constants.MaxSpikes; i++ {
		var spikes []Spike
		for _, s := range gMap.Spikes {
			spikes = append(spikes, *s)
		}
		//x, y := randXYNoOverlap(spikes, SpikesSpacing)
		x, y := gMap.Spikes.randXY(constants.SpikesSpacing)
		spike := NewSpike(x, y, constants.SpikeWeight)
		gMap.Spikes = append(gMap.Spikes, spike)
	}
}

func (gMap *GameMap) populateFood() {
	if len(gMap.Foods) < constants.MaxNumFood {
		for i := 0; i < constants.MaxNumFood-len(gMap.Foods); i++ {
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

func (gMap *GameMap) notifyAllPlayers(data []byte) {
	for _, p := range gMap.Players.Real() {
		client, err := gMap.PlayerReverseLookUp(p.Uuid)
		if err != nil {
			log.Println(err)
			continue
		}
		client.Emit(data)
	}
}

func (gMap *GameMap) publishAdminStats() {
	for {
		var result = make(map[string]interface{})
		botsCount := gMap.Players.BotsCount()
		stats := gMap.GetStats()

		result["botsCount"] = botsCount
		result["playersCount"] = len(gMap.Players) - botsCount
		result["topPlayers"] = stats
		//gMap.Hub.Emit("stats", result, "admin")
		time.Sleep(2 * time.Second)
	}
}

func (gMap *GameMap) notifyPlayer(player *player.Player) {
	movedEvent := &MovedEvent{"moved", player.X, player.Y, player.Weight, player.Zoom}
	if bytes, err := MovedSchema.Encode(movedEvent); err != nil {
		log.Println(err)
	} else {
		gMap.Hub.Emit(bytes, player.Uuid)
	}
	players := gMap.Players.Closest(player, constants.NumPlayersResponse)
	plUpdateEvent := &PlayersUpdateEvent{"pUpdated", players}
	if bytes, err := PlayersUpdateSchema.Encode(plUpdateEvent); err != nil {
		log.Println(err)
	} else {
		gMap.Hub.Emit(bytes, player.Uuid)
	}
	foods := gMap.Foods.Closest(player, constants.NumFoodResponse)
	gMap.Hub.Emit("foodUpdated", foods, player.Uuid)
}

func (gMap *GameMap) Run() {
	counter := 0
	go gMap.publishStats()
	go gMap.Hub.Run()
	go gMap.PopulateSpikes()
	go gMap.publishAdminStats()
	for {
		counter++
		var wg sync.WaitGroup
		if counter > 30 {
			counter = 0
			gMap.populateFood()
			gMap.PopulateBots()
		}
		for i := range gMap.Players {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				pl := gMap.Players[i]
				if pl.IsBot {
					pl.MakeMove(gMap)
				}
				pl.UpdatePosition(gMap)
				pl.PassiveWeightLoss()
				if pl.IsBot {
					return
				}
				if counter%2 == 0 {
					gMap.notifyPlayer(pl)
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
	players := gMap.Players.Largest(constants.StatsNumber)
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

func (gMap *GameMap) createFood(x float32, y float32) *food.Food {
	food := food.New(x, y, constants.FoodWeight)
	gMap.Foods = append(gMap.Foods, food)
	return food
}

func (gMap *GameMap) createRandomFood() *food.Food {
	x, y := randXY()
	return gMap.createFood(x, y)
}

func (gMap *GameMap) CreatePlayer(nickname string, isBot bool) map[string]interface{} {
	x, y := randXY()
	pl := player.NewPlayer(x, y, constants.MinWeight, nickname, isBot)
	if len(gMap.Players) >= constants.MaxPlayers {
		gMap.removePlayerIndex(0)
	}
	gMap.Players = append(gMap.Players, pl)
	var result = make(map[string]interface{})
	result["player"] = pl.GetSelfPlayer()
	result["spikes"] = gMap.Spikes
	return result
}

func (gMap *GameMap) removeEatableFood() {
	for _, pl := range gMap.Players {
		var filteredFoods []*food.Food
		for _, f := range gMap.Foods {
			if pl.FoodEatable(f) {
				pl.EatEntity(f)
			} else {
				filteredFoods = append(filteredFoods, f)
			}
		}
		gMap.Foods = filteredFoods
	}
}

func (gMap *GameMap) SpikeCollisions(pl *player.Player) bool {
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
