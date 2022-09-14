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
			newX: 'uint16',
			newY: 'uint16'
		}),
		moved: genericSchema.extends({
			x: 'float32',
			y: 'float32',
			weight: 'float32',
			zoom: 'float32'
		}),
		start: genericSchema.extends({
			nickname: 'string'
		}),
		fUpdated: genericSchema.extends({
			foods: [
				{
					x: 'float32',
					y: 'float32',
					weight: 'float32',
					color: ['uint8']
				}
			]
		}),
		pUpdated: genericSchema.extends({
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
				color: ['uint8']
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
		console.log('socket initialized');
		this.registerPongListener();
		this.socket.onOpen((event) => {
			console.log('onOpen called', event)
			this.pingPong();
		});
	}

	onOpen(callback: (event: Event) => void) {
		this.socket.onOpen(callback)
	}

	startGame({nickname}: { nickname: string }): Promise<InitialData> {
		const schema = this.schemas.start;
		this.socket.emit(schema.encode({event: 'start', nickname}).toBuffer());
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
		this.socket.on('pUpdated', callback, this.schemas.pUpdated);
	}

	onFoodUpdate(callback: (data: { foods: FoodData[] }) => void) {
		this.socket.on('fUpdated', callback, this.schemas.fUpdated);
	}

	onStatsUpdate(callback: (data: { topPlayers: StatsUpdate[] }) => void) {
		this.socket.on('stats', callback, this.schemas.fUpdated);
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
			console.log('ping', ping);
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
