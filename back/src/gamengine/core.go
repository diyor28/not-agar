package gamengine

import (
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/gamengine/constants"
	_map "github.com/diyor28/not-agar/src/gamengine/map"
	"github.com/diyor28/not-agar/src/gamengine/map/entity"
	"github.com/diyor28/not-agar/src/gamengine/map/players"
	"github.com/diyor28/not-agar/src/sockethub"
	"log"
	"sync"
	"time"
)

type GameEngine struct {
	Hub        *sockethub.Hub
	Map        *_map.Map
	PlayersMap map[*sockethub.Client]entity.Id
	framerate  int
	runEvery   time.Duration
}

func NewGameMap(framerate int) *GameEngine {
	hub := sockethub.NewHub()
	gameMap := _map.New()
	delta := time.Duration(1000/framerate) * time.Millisecond
	engine := GameEngine{
		Hub:        hub,
		Map:        gameMap,
		PlayersMap: make(map[*sockethub.Client]entity.Id),
		framerate:  framerate,
		runEvery:   delta,
	}
	return &engine
}

func (eng *GameEngine) Loop() {
	eng.Map.PopulateBots()
	eng.populateFood()
	var wg sync.WaitGroup
	for i := range eng.Map.Players.Players {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pl := eng.Map.Players.Players[i]
			if pl.IsBot {
				eng.MakeMove(pl)
			}
			pl.UpdatePosition(eng.runEvery)
			pl.PassiveWeightLoss()
			if pl.IsBot {
				return
			}
			err := eng.notifyPlayer(pl)
			if err != nil {
				log.Println("error when calling eng.notifyPlayer(): ", err)
			}
		}(i)
	}
	wg.Wait()
	eng.removeDeadPlayers()
	eng.removeEatableFood()
}

func (eng *GameEngine) PlayerReverseLookUp(id entity.Id) (*sockethub.Client, error) {
	for k, pId := range eng.PlayersMap {
		if pId == id {
			return k, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no players for id %d found", id))
}

func (eng *GameEngine) HandleMoveEvent(event *MoveEvent, client *sockethub.Client) {
	pl, err := eng.Map.Players.Update(eng.PlayersMap[client], event.NewX, event.NewY)
	if err != nil {
		log.Println(err)
		return
	}
	data, err := MovedSchema.Encode(&MovedEvent{
		constants.Moved,
		pl.X,
		pl.Y,
		pl.VelocityX,
		pl.VelocityY,
		pl.Weight,
		pl.Zoom,
		pl.Shell.Points,
	})
	if err != nil {
		log.Println(err)
		return
	}
	if err := client.Emit(data.Bytes()); err != nil {
		log.Println(err)
	}
}

func (eng *GameEngine) SendPong(data []byte, client *sockethub.Client) {
	pingData := make(map[string]interface{})
	if err := PingPongSchema.Decode(data, &pingData); err != nil {
		log.Println("PingPongSchema.Decode(): ", err)
		return
	}
	pongEvent := map[string]interface{}{"event": constants.Pong, "timestamp": pingData["timestamp"]}
	if data, err := PingPongSchema.Encode(&pongEvent); err != nil {
		log.Println("PingPongSchema.Encode(): ", err)
	} else {
		client.Emit(data.Bytes())
	}
}

func (eng *GameEngine) publishStats() {
	for range time.Tick(time.Duration(2000) * time.Millisecond) {
		stats := eng.Map.GetStats()
		statsEvent := &PlayerStatsEvent{
			Event:      constants.StatsUpdate,
			TopPlayers: stats,
		}
		if data, err := PlayerStatsSchema.Encode(statsEvent); err != nil {
			log.Println(err)
		} else {
			eng.notifyAllPlayers(data.Bytes())
		}
	}
}

func (eng *GameEngine) notifyAllPlayers(data []byte) {
	for client, _ := range eng.PlayersMap {
		client.Emit(data)
	}
}

func (eng *GameEngine) publishAdminStats() {
	for {
		var result = make(map[string]interface{})
		botsCount := eng.Map.Players.BotsCount()
		stats := eng.Map.GetStats()

		result["botsCount"] = botsCount
		result["playersCount"] = eng.Map.Players.Len() - botsCount
		result["topPlayers"] = stats
		// TODO: fix later
		//eng.Hub.Emit("stats", result, "admin")
		time.Sleep(2 * time.Second)
	}
}

func (eng *GameEngine) notifyPlayer(pl *players.Player) error {
	client, err := eng.PlayerReverseLookUp(pl.Id)
	if err != nil {
		log.Println(err)
		return err
	}
	movedEvent := &MovedEvent{
		constants.Moved,
		pl.X,
		pl.Y,
		pl.VelocityX,
		pl.VelocityY,
		pl.Weight,
		pl.Zoom,
		pl.Shell.Points,
	}
	if data, err := MovedSchema.Encode(movedEvent); err != nil {
		log.Println(err)
		return err
	} else {
		if err := client.Emit(data.Bytes()); err != nil {
			pl.IsDead = true
			delete(eng.PlayersMap, client)
		}
	}
	plrs := eng.Map.Players.Closest(pl, constants.NumPlayersResponse)
	plUpdateEvent := &PlayersUpdatedEvent{constants.PlayersUpdate, plrs}
	if data, err := PlayersUpdatedSchema.Encode(plUpdateEvent); err != nil {
		log.Println(err)
		return err
	} else {
		if err := client.Emit(data.Bytes()); err != nil {
			pl.IsDead = true
			delete(eng.PlayersMap, client)
		}
	}
	return nil
}

func (eng *GameEngine) MakeMove(pl *players.Player) {
	food := eng.Map.Food.Closest(pl, 1)
	closestFood := food[0]
	pl.UpdateDirection(closestFood.X, closestFood.Y)
}

func (eng *GameEngine) Run() {

	eng.Hub.OnMessage(func(data []byte, client *sockethub.Client) {
		event := &GenericEvent{}
		if err := GenericSchema.Decode(data, event); err != nil {
			log.Println("GenericSchema.Decode(): ", err)
			return
		}
		switch event.Event {
		case constants.Move:
			moveEvent := &MoveEvent{}
			if err := MoveSchema.Decode(data, moveEvent); err != nil {
				log.Println("MoveSchema.Decode(): ", err)
				return
			}
			eng.HandleMoveEvent(moveEvent, client)
		case constants.Start:
			var startEvent = make(map[string]interface{})
			if err := StartSchema.Decode(data, &startEvent); err != nil {
				log.Println("StartSchema.Decode(): ", err)
				return
			}
			nickname := startEvent["nickname"].(string)
			player := eng.Map.CreatePlayer(nickname, false)
			spikes := eng.Map.Spikes.Spikes
			food := eng.Map.Food.Food
			eng.PlayersMap[client] = player.Id
			startedEvent := &StartedEvent{
				Event: constants.Started,
				Player: &StartedEventPlayer{
					player.X,
					player.Y,
					player.Weight,
					player.Color,
					player.Shell.Points,
				},
				Spikes: spikes,
				Food:   food,
			}
			if data, err := StartedSchema.Encode(startedEvent); err != nil {
				log.Println("StartedSchema.Encode(): ", err)
				return
			} else {
				if err := client.Emit(data.Bytes()); err != nil {
					log.Println(err)
				}
				client.Join(fmt.Sprintf("player/%d", player.Id))
			}
		case constants.Ping:
			eng.SendPong(data, client)
		}
	})
	go eng.publishStats()
	go eng.Hub.Run()
	go eng.Map.PopulateSpikes()
	go eng.publishAdminStats()
	for range time.Tick(eng.runEvery) {
		eng.Loop()
	}
}

func (eng *GameEngine) populateFood() {
	createdFood := eng.Map.PopulateFood()
	if len(createdFood) == 0 {
		return
	}
	event := &FoodCreatedEvent{
		Event: constants.FoodCreated,
		Food:  createdFood,
	}
	if data, err := FoodCreatedSchema.Encode(event); err == nil {
		eng.notifyAllPlayers(data.Bytes())
	} else {
		log.Println("FoodCreatedSchema.Encode()", err)
	}

}

func (eng *GameEngine) removeDeadPlayers() {
	ripEvent := map[string]interface{}{"event": constants.Rip}
	deadPlayers := eng.Map.RemoveDeadPlayers()
	for _, pl := range deadPlayers {
		client, err := eng.PlayerReverseLookUp(pl.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		if data, err := GenericSchema.Encode(&ripEvent); err != nil {
			log.Println(err)
		} else {
			if err := client.Emit(data.Bytes()); err != nil {
				log.Println(err)
			}
		}
	}
}

func (eng *GameEngine) removeEatableFood() {
	eatenFood := eng.Map.RemoveEatableFood()
	for _, f := range eatenFood {
		event := &FoodEatenEvent{
			Event: constants.FoodEaten,
			Id:    f.Id,
		}
		if data, err := FoodEatenSchema.Encode(event); err == nil {
			eng.notifyAllPlayers(data.Bytes())
		} else {
			log.Println("FoodEatenSchema.Decode()", err)
		}
	}
}
