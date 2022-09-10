import {Schema} from "../codec";

export class BaseSocket {
	public readonly socket: WebSocket;
	protected openListeners: Function[] = [];
	protected onMessageListeners: Function[] = [];

	constructor(url: string) {
		this.socket = new WebSocket(url);
		this.socket.onmessage = this.handleMessage.bind(this);
		this.socket.onopen = (event: Event) => {
			this.openListeners.forEach(callback => {
				callback(event)
			})
		}
	}

	onOpen(callback: () => void) {
		this.openListeners.push(callback)
	}

	onMessage(callback: (data: ArrayBuffer | Record<string, any>) => void) {
		this.onMessageListeners.push(callback)
	}

	close(code?: number, reason?: string) {
		this.socket.close(code, reason)
	}

	protected handleMessage(event: MessageEvent) {
		this.onMessageListeners.forEach(callback => {
			if (event.data instanceof ArrayBuffer) {
				callback(event.data);
			} else {
				callback(JSON.parse(event.data));
			}
		})
	}
}

export class BinarySocket extends BaseSocket {
	schema: Schema;
	protected eventListeners: { event: string, callback: (data: any) => void, schema: Schema }[] = [];

	constructor(url: string, schema: Schema) {
		super(url);
		this.schema = schema;
		this.onMessage((data) => {
			if (!(data instanceof ArrayBuffer)) {
				throw new Error(`Received ${typeof data} message in BinarySocket`);
			}
			const dataBuffer = Buffer.from(data);
			const {event} = this.schema.decode(dataBuffer);
			this.eventListeners.forEach(listener => {
				if (listener.event === event) {
					listener.callback(listener.schema.decode(dataBuffer));
				}
			})
		})
	}

	emit(data: ArrayBuffer | Buffer) {
		if (this.socket.readyState === WebSocket.CLOSED)
			throw new Error('Socket closed');
		if (this.socket.readyState === WebSocket.CONNECTING)
			throw new Error('Socket connecting');
		this.socket.send(data);
	}

	on(event: string, callback: (data: any) => void, schema: Schema) {
		this.eventListeners.push({event, callback, schema});
	}
}

export class JsonSocket extends BaseSocket {
	protected eventListeners: { event: string, callback: (data: any) => void }[] = [];

	constructor(url: string) {
		super(url);
		this.onMessage((data) => {
			if (data instanceof ArrayBuffer) {
				throw new Error('Received ArrayBuffer in JsonSocket');
			}
			this.eventListeners.forEach(listener => {
				if (listener.event === data.event) {
					listener.callback(data.data);
				}
			})

		})
	}

	emit(event: string, data: any) {
		if (this.socket.readyState === WebSocket.CLOSED)
			throw new Error('Socket closed');
		if (this.socket.readyState === WebSocket.CONNECTING)
			throw new Error('Socket connecting');
		this.socket.send(JSON.stringify({event, data}))
	}

	on(event: string, callback: (data: any) => void) {
		this.eventListeners.push({event, callback})
	}
}