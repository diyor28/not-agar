import {Schema} from "../codec";
import {BinarySocket} from "./socket";
import {StatsUpdate} from "../engine/GameEngine";
import {FoodData, InitialData, MoveCommand, MovedEvent, PlayerData, SelfPlayerData, SpikeData} from "./types";

export type {MovedEvent, SpikeData, SelfPlayerData, EntityData, MoveCommand, PlayerData, FoodData, InitialData} from './types';

export default class GameClient {
	ping: number | null = null;
	socket: BinarySocket;
	pingInterval: number
	private schemas = {
		ping: new Schema({
			event: 'string',
			timestamp: 'uint32'
		}),
		move: new Schema({
			event: 'string',
			newX: 'uint16',
			newY: 'uint16'
		}),
		moved: new Schema({
			event: 'string',
			x: 'float32',
			y: 'float32',
			weight: 'float32',
			zoom: 'float32'
		}),
		start: new Schema({
			event: 'string',
			nickname: 'string'
		}),
		started: new Schema({
			event: 'string',
			player: {
				uuid: 'string',
				nickname: 'string',
				color: ['uint8'],
				x: 'float32',
				y: 'float32',
				weight: 'float32',
				speed: 'float32',
				zoom: 'float32'
			},
			spikes: [
				{
					uuid: 'string',
					x: 'float32',
					y: 'float32',
					weight: 'float32'
				}
			]
		})
	};

	constructor(url: string, pingInterval: number) {
		this.pingInterval = pingInterval;
		const schema = new Schema({event: 'string'});
		this.socket = new BinarySocket(url, schema);
		this.registerPongListener();
		this.socket.onOpen(() => {
			this.pingPong();
		});
	}

	public startGame({nickname}: { nickname: string }): Promise<InitialData> {
		const schema = this.schemas.start;
		this.socket.emit(schema.encode({event: 'start', nickname}));
		return new Promise(resolve => this.socket.on('started', resolve, this.schemas.started));
	}

	public move(data: MoveCommand) {
		const schema = this.schemas.move;
		this.socket.emit(schema.encode({event: 'move', ...data}));
	}

	public onMove(callback: (data: MovedEvent) => void) {
		this.socket.on('moved', callback, this.schemas.moved);
	}

	public onPlayersUpdate(callback: (data: PlayerData[]) => void) {
		// this.socket.on('pUpdated', callback, this.schemas)
	}

	public onFoodUpdate(callback: (data: FoodData[]) => void) {

	}

	public onStatsUpdate(callback: (data: StatsUpdate[]) => void) {

	}

	public onPong(callback: (data: { ping: number }) => void) {

	}

	public onRip(callback: () => void) {

	}

	private registerPongListener() {
		const listener = (data: { timestamp: number }) => {
			const currentTime = new Date().getTime();
			const ping = currentTime - data.timestamp;
			// this.onPong({ping});
			this.ping = ping;
		};
		this.socket.on('pong', listener, this.schemas.ping);
	}

	private pingPong() {
		const data = {timestamp: new Date().getTime()};
		this.socket.emit(this.schemas.ping.encode({...data, event: 'ping'}))
		setTimeout(() => this.pingPong(), this.pingInterval);
	}
}
