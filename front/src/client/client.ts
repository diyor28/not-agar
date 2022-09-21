import {SocketWrapper} from "./socket";
import {StatsUpdate} from "../engine/GameEngine";
import {FoodData, InitialData, MoveCommand, MovedEvent, PlayerData} from "./types";
import {
	foodSchema,
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

type GameEvent = 'open' | 'error' | 'moved' | 'started' | 'players' | 'food' | 'stats' | 'pong' | 'rip';
type GameData =
	Event
	| MovedEvent
	| InitialData
	| { players: PlayerData[] }
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
				case 'moved':
					return this.bus.emit(event, movedSchema.decode(data));
				case 'started':
					return this.bus.emit(event, startedSchema.decode(data));
				case 'players':
					return this.bus.emit(event, playersSchema.decode(data));
				case 'food':
					return this.bus.emit(event, foodSchema.decode(data));
				case 'stats':
					return this.bus.emit(event, statsSchema.decode(data));
				case 'pong':
					const {timestamp} = pingSchema.decode(data);
					const ping = new Date().getTime() - timestamp;
					this.ping = ping;
					return this.bus.emit(event, {ping});
				case 'rip':
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
	on(event: 'moved', callback: (data: MovedEvent) => void): void
	on(event: 'started', callback: (data: InitialData) => void): void
	on(event: 'players', callback: (data: { players: PlayerData[] }) => void): void
	on(event: 'food', callback: (data: { food: FoodData[] }) => void): void
	on(event: 'stats', callback: (data: { topPlayers: StatsUpdate[] }) => void): void
	on(event: 'pong', callback: (data: { ping: number }) => void): void
	on(event: 'rip', callback: () => void): void
	on(event: GameEvent, callback: GameCallback) {
		this.bus.on(event, callback);
	}

	once(event: 'open', callback: (data: Event) => void): void
	once(event: 'error', callback: (data: Event) => void): void
	once(event: 'moved', callback: (data: MovedEvent) => void): void
	once(event: 'started', callback: (data: InitialData) => void): void
	once(event: 'players', callback: (data: { players: PlayerData[] }) => void): void
	once(event: 'food', callback: (data: { food: FoodData[] }) => void): void
	once(event: 'stats', callback: (data: { topPlayers: StatsUpdate[] }) => void): void
	once(event: 'pong', callback: (data: { ping: number }) => void): void
	once(event: 'rip', callback: () => void): void
	once(event: GameEvent, callback: GameCallback) {
		this.bus.once(event, callback);
	}

	startGame({nickname}: { nickname: string }): Promise<InitialData> {
		const data = startSchema.encode({event: 'start', nickname});
		this.socket.emit(data.toBuffer());
		return new Promise(resolve => this.bus.on('started', resolve));
	}

	move(data: MoveCommand) {
		this.socket.emit(moveSchema.encode({event: 'move', ...data}).toBuffer());
	}

	private pingPong() {
		const data = {timestamp: new Date().getTime()};
		this.socket.emit(pingSchema.encode({event: 'ping', ...data}).toBuffer())
		setTimeout(() => this.pingPong(), this.pingInterval);
	}
}
