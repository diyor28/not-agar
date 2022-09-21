import {Schema} from "../codec";
import {BinarySocket} from "./socket";
import {StatsUpdate} from "../engine/GameEngine";
import {FoodData, InitialData, MoveCommand, MovedEvent, PlayerData, SelfPlayerData, SpikeData} from "./types";

export type {
	MovedEvent, SpikeData, SelfPlayerData, EntityData, MoveCommand, PlayerData, FoodData, InitialData
} from './types';

const genericSchema = new Schema({
	event: 'string'
})

export default class GameClient {
	ping: number | null = null;
	socket: BinarySocket;
	pingInterval: number
	private onPongListeners: ((data: {ping: number}) => void)[] = [];

	private schemas = {
		ping: genericSchema.extends({
			timestamp: 'uint64'
		}),
		move: genericSchema.extends({
			newX: 'float32',
			newY: 'float32'
		}),
		moved: genericSchema.extends({
			x: 'float32',
			y: 'float32',
			weight: 'float32',
			velocityX: 'float32',
			velocityY: 'float32',
			zoom: 'float32',
			points: [{x: 'float32', y: 'float32'}]
		}),
		start: genericSchema.extends({
			nickname: 'string'
		}),
		food: genericSchema.extends({
			food: [
				{
					x: 'float32',
					y: 'float32',
					weight: 'float32',
					color: ['uint8']
				}
			]
		}),
		players: genericSchema.extends({
			players: [
				{
					x: 'float32',
					y: 'float32',
					weight: 'float32',
					nickname: 'string',
					color: ['uint8']
				},
			]
		}),
		started: genericSchema.extends({
			player: {
				x: 'float32',
				y: 'float32',
				weight: 'float32',
				color: ['uint8'],
				points: [{x: 'float32', y: 'float32'}]
			},
			spikes: [
				{
					x: 'float32',
					y: 'float32',
					weight: 'float32'
				}
			]
		}),
		stats: genericSchema.extends({
			topPlayers: [
				{
					nickname: 'string',
					weight: 'int16'
				}
			]
		}),
		rip: genericSchema
	};

	constructor(url: string, pingInterval: number) {
		const schema = new Schema({event: 'string'});
		this.pingInterval = pingInterval;
		this.socket = new BinarySocket(url, schema);
		this.registerPongListener();
		this.socket.onOpen((event) => {
			this.pingPong();
		});
	}

	connect() {
		return this.socket.connect();
	}

	onOpen(callback: (event: Event) => void) {
		this.socket.onOpen(callback)
	}

	startGame({nickname}: { nickname: string }): Promise<InitialData> {
		const data = this.schemas.start.encode({event: 'start', nickname});
		this.socket.emit(data.toBuffer());
		return new Promise(resolve => this.socket.on('started', resolve, this.schemas.started));
	}

	move(data: MoveCommand) {
		const schema = this.schemas.move;
		this.socket.emit(schema.encode({event: 'move', ...data}).toBuffer());
	}

	onMove(callback: (data: MovedEvent) => void) {
		this.socket.on('moved', callback, this.schemas.moved);
	}

	onPlayersUpdate(callback: (data: { players: PlayerData[] }) => void) {
		this.socket.on('players', callback, this.schemas.players);
	}

	onFoodUpdate(callback: (data: { food: FoodData[] }) => void) {
		this.socket.on('food', callback, this.schemas.food);
	}

	onStatsUpdate(callback: (data: { topPlayers: StatsUpdate[] }) => void) {
		this.socket.on('stats', callback, this.schemas.stats);
	}

	onPong(callback: (data: { ping: number }) => void) {
		this.onPongListeners.push(callback);
	}

	onRip(callback: () => void) {
		this.socket.on('rip', callback, this.schemas.rip);
	}

	private registerPongListener() {
		const listener = (data: { timestamp: number }) => {
			const currentTime = new Date().getTime();
			const ping = currentTime - data.timestamp;
			this.ping = ping;
			this.onPongListeners.forEach(callback => {
				callback({ping});
			});
		};
		this.socket.on('pong', listener, this.schemas.ping);
	}

	private pingPong() {
		const data = {timestamp: new Date().getTime()};
		this.socket.emit(this.schemas.ping.encode({event: 'ping', ...data}).toBuffer())
		setTimeout(() => this.pingPong(), this.pingInterval);
	}
}
