package gamengine

import (
	"errors"
	"fmt"
	_map "github.com/diyor28/not-agar/cmd/gamengine/map"
	"github.com/diyor28/not-agar/cmd/gamengine/map/entity"
	"github.com/diyor28/not-agar/cmd/gamengine/map/food"
	"github.com/diyor28/not-agar/cmd/gamengine/map/players"
	"github.com/diyor28/not-agar/cmd/gamengine/map/players/shell"
	"github.com/diyor28/not-agar/cmd/gamengine/map/spikes"
	"github.com/diyor28/not-agar/cmd/gamengine/schemas"
	"github.com/diyor28/not-agar/cmd/sockethub"
	"github.com/diyor28/not-agar/pkg/constants"
	"log"
	"sync"
	"time"
)

func castFood(food []*food.Food) []*schemas.Food {
	res := make([]*schemas.Food, len(food))
	for i, f := range food {
		res[i] = &schemas.Food{Id: f.Id, X: f.X, Y: f.Y, Weight: f.Weight, Color: f.Color}
	}
	return res
}

func castSpikes(spikes []*spikes.Spike) []*schemas.Spike {
	res := make([]*schemas.Spike, len(spikes))
	for i, s := range spikes {
		res[i] = &schemas.Spike{X: s.X, Y: s.Y, Weight: s.Weight}
	}
	return res
}

func castPlayers(pls []*players.Player) []*schemas.Player {
	res := make([]*schemas.Player, len(pls))
	for i, p := range pls {
		res[i] = &schemas.Player{X: uint16(p.X), Y: uint16(p.Y), Weight: p.Weight, Nickname: p.Nickname, Color: p.Color}
	}
	return res
}

func castPoints(points []*shell.Point) []*schemas.Point {
	res := make([]*schemas.Point, len(points))
	for i, p := range points {
		res[i] = &schemas.Point{X: int16(p.X * 100), Y: int16(p.Y * 100)}
	}
	return res
}

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

func (eng *GameEngine) HandleMoveEvent(event *schemas.MoveEvent, client *sockethub.Client) {
	pl, err := eng.Map.Players.Update(eng.PlayersMap[client], event.NewX, event.NewY)
	if err != nil {
		log.Println(err)
		return
	}
	data, err := schemas.MovedSchema.Encode(&schemas.MovedEvent{
		Event:     constants.Moved,
		X:         pl.X,
		Y:         pl.Y,
		VelocityX: pl.VelocityX,
		VelocityY: pl.VelocityY,
		Weight:    pl.Weight,
		Zoom:      pl.Zoom,
		Points:    castPoints(pl.Shell.Points),
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
	if err := schemas.PingPongSchema.Decode(data, &pingData); err != nil {
		log.Println("PingPongSchema.Decode(): ", err)
		return
	}
	pongEvent := map[string]interface{}{"event": constants.Pong, "timestamp": pingData["timestamp"]}
	if data, err := schemas.PingPongSchema.Encode(&pongEvent); err != nil {
		log.Println("PingPongSchema.Encode(): ", err)
	} else {
		client.Emit(data.Bytes())
	}
}

func (eng *GameEngine) publishStats() {
	for range time.Tick(time.Duration(2000) * time.Millisecond) {
		stats := eng.Map.GetStats()
		statsEvent := &schemas.PlayerStatsEvent{
			Event:      constants.StatsUpdate,
			TopPlayers: stats,
		}
		if data, err := schemas.PlayerStatsSchema.Encode(statsEvent); err != nil {
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
	movedEvent := &schemas.MovedEvent{
		Event:     constants.Moved,
		X:         pl.X,
		Y:         pl.Y,
		VelocityX: pl.VelocityX,
		VelocityY: pl.VelocityY,
		Weight:    pl.Weight,
		Zoom:      pl.Zoom,
		Points:    castPoints(pl.Shell.Points),
	}
	if data, err := schemas.MovedSchema.Encode(movedEvent); err != nil {
		log.Println(err)
		return err
	} else {
		if err := client.Emit(data.Bytes()); err != nil {
			pl.IsDead = true
			delete(eng.PlayersMap, client)
		}
	}
	plrs := eng.Map.Players.Closest(pl, constants.NumPlayersResponse)
	plUpdateEvent := &schemas.PlayersUpdatedEvent{
		Event:   constants.PlayersUpdate,
		Players: castPlayers(plrs),
	}
	if data, err := schemas.PlayersUpdatedSchema.Encode(plUpdateEvent); err != nil {
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
		event := &schemas.GenericEvent{}
		if err := schemas.GenericSchema.Decode(data, event); err != nil {
			log.Println("GenericSchema.Decode(): ", err)
			return
		}
		switch event.Event {
		case constants.Move:
			moveEvent := &schemas.MoveEvent{}
			if err := schemas.MoveSchema.Decode(data, moveEvent); err != nil {
				log.Println("MoveSchema.Decode(): ", err)
				return
			}
			eng.HandleMoveEvent(moveEvent, client)
		case constants.Start:
			var startEvent = make(map[string]interface{})
			if err := schemas.StartSchema.Decode(data, &startEvent); err != nil {
				log.Println("StartSchema.Decode(): ", err)
				return
			}
			nickname := startEvent["nickname"].(string)
			player := eng.Map.CreatePlayer(nickname, false)
			eng.PlayersMap[client] = player.Id
			startedEvent := &schemas.StartedEvent{
				Event: constants.Started,
				Player: &schemas.StartedEventPlayer{
					X:      player.X,
					Y:      player.Y,
					Weight: player.Weight,
					Color:  player.Color,
					Points: castPoints(player.Shell.Points),
				},
				Spikes: castSpikes(eng.Map.Spikes.Spikes),
				Food:   castFood(eng.Map.Food.Food),
			}
			if data, err := schemas.StartedSchema.Encode(startedEvent); err != nil {
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
	event := &schemas.FoodCreatedEvent{
		Event: constants.FoodCreated,
		Food:  castFood(createdFood),
	}
	if data, err := schemas.FoodCreatedSchema.Encode(event); err == nil {
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
		if data, err := schemas.GenericSchema.Encode(&ripEvent); err != nil {
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
		event := &schemas.FoodEatenEvent{
			Event: constants.FoodEaten,
			Id:    f.Id,
		}
		if data, err := schemas.FoodEatenSchema.Encode(event); err == nil {
			eng.notifyAllPlayers(data.Bytes())
		} else {
			log.Println("FoodEatenSchema.Decode()", err)
		}
	}
}
