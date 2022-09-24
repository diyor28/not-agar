import {SocketWrapper} from "./socket";
import {StatsUpdate} from "../engine/GameEngine";
import {FoodData, InitialData, MoveCommand, MovedEvent, PlayerData} from "./types";
import {
	foodCreatedSchema,
	foodEatenSchema,
	GameEvent,
	genericSchema,
	movedSchema,
	moveSchema,
	pingSchema,
	playersSchema,
	startedSchema,
	startSchema,
	statsSchema
} from './schemas'
import {EventBus} from "./eventBus";

type MixedGameEvent = 'open' | 'error' | GameEvent;
type GameData =
	Event
	| MovedEvent
	| InitialData
	| { players: PlayerData[] }
	| { id: number }
	| { food: FoodData[] }
	| { topPlayers: StatsUpdate[] }
	| { ping: number };
type GameCallback = (() => void) | ((data: GameData) => void);

export class GameClient {
	ping: number | null = null;
	socket: SocketWrapper;
	pingInterval: number
	private bus: EventBus;

	constructor(url: string, pingInterval: number) {
		this.pingInterval = pingInterval;
		this.socket = new SocketWrapper(url);
		this.bus = new EventBus();
		this.socket.once('open', this.pingPong.bind(this));
		this.socket.on('open', (event) => {
			this.bus.emit('open', event);
		});
		this.socket.on('error', (event) => {
			this.bus.emit('error', event);
		});
		this.socket.on('message', (data) => {
			const {event} = genericSchema.decode(data);
			switch (event) {
				case GameEvent.Moved:
					return this.bus.emit(event, movedSchema.decode(data));
				case GameEvent.Started:
					return this.bus.emit(event, startedSchema.decode(data));
				case GameEvent.PlayersUpdate:
					return this.bus.emit(event, playersSchema.decode(data));
				case GameEvent.FoodEaten:
					return this.bus.emit(event, foodEatenSchema.decode(data));
				case GameEvent.FoodCreated:
					return this.bus.emit(event, foodCreatedSchema.decode(data));
				case GameEvent.StatsUpdate:
					return this.bus.emit(event, statsSchema.decode(data));
				case GameEvent.Pong:
					const {timestamp} = pingSchema.decode(data);
					const ping = new Date().getTime() - timestamp;
					this.ping = ping;
					return this.bus.emit(event, {ping});
				case GameEvent.Rip:
					return this.bus.emit(event, {});
				default:
					console.log(`Received unknown event: ${event}`)
			}
		});
	}

	connect() {
		return this.socket.connect();
	}

	on(event: 'open', callback: (data: Event) => void): void
	on(event: 'error', callback: (data: Event) => void): void
	on(event: GameEvent.Moved, callback: (data: MovedEvent) => void): void
	on(event: GameEvent.Started, callback: (data: InitialData) => void): void
	on(event: GameEvent.PlayersUpdate, callback: (data: { players: PlayerData[] }) => void): void
	on(event: GameEvent.FoodEaten, callback: (data: { id: number }) => void): void
	on(event: GameEvent.FoodCreated, callback: (data: { food: FoodData[] }) => void): void
	on(event: GameEvent.StatsUpdate, callback: (data: { topPlayers: StatsUpdate[] }) => void): void
	on(event: GameEvent.Pong, callback: (data: { ping: number }) => void): void
	on(event: GameEvent.Rip, callback: () => void): void
	on(event: MixedGameEvent, callback: GameCallback) {
		this.bus.on(event, callback);
	}

	once(event: 'open', callback: (data: Event) => void): void
	once(event: 'error', callback: (data: Event) => void): void
	once(event: GameEvent.Moved, callback: (data: MovedEvent) => void): void
	once(event: GameEvent.Started, callback: (data: InitialData) => void): void
	once(event: GameEvent.PlayersUpdate, callback: (data: { players: PlayerData[] }) => void): void
	once(event: GameEvent.FoodEaten, callback: (data: { id: number }) => void): void
	once(event: GameEvent.FoodCreated, callback: (data: { food: FoodData[] }) => void): void
	once(event: GameEvent.StatsUpdate, callback: (data: { topPlayers: StatsUpdate[] }) => void): void
	once(event: GameEvent.Pong, callback: (data: { ping: number }) => void): void
	once(event: GameEvent.Rip, callback: () => void): void
	once(event: MixedGameEvent, callback: GameCallback) {
		this.bus.once(event, callback);
	}

	startGame({nickname}: { nickname: string }): Promise<InitialData> {
		const data = startSchema.encode({event: GameEvent.Start, nickname});
		this.socket.emit(data.toBuffer());
		return new Promise(resolve => this.bus.on(GameEvent.Started, resolve));
	}

	move(data: MoveCommand) {
		this.socket.emit(moveSchema.encode({event: GameEvent.Move, ...data}).toBuffer());
	}

	private pingPong() {
		const data = {timestamp: new Date().getTime()};
		this.socket.emit(pingSchema.encode({event: GameEvent.Ping, ...data}).toBuffer())
		setTimeout(() => this.pingPong(), this.pingInterval);
	}
}
