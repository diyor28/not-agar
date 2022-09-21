import {Schema} from "../codec";

export class BaseSocket {
	socket?: WebSocket;
	url: string
	protected openListeners: Function[] = [];
	protected errorListeners: Function[] = [];
	protected onMessageListeners: Function[] = [];

	constructor(url: string) {
		this.url = url;
	}

	async connect(): Promise<Event> {
		this.socket = new WebSocket(this.url);
		const connectPromise = new Promise<Event>((resolve, reject) => {
			this.onOpen(resolve)
			this.onError(reject)
		})
		this.socket.onopen = (event: Event) => {
			this.openListeners.forEach(callback => {
				callback(event);
			});
		}
		this.socket.onerror = (event: Event) => {
			this.errorListeners.forEach(callback => {
				callback(event);
			});
		}
		this.socket.onmessage = this.handleMessage.bind(this);
		return connectPromise
	}

	onOpen(callback: (event: Event) => void) {
		this.openListeners.push(callback);
	}

	onError(callback: (event: Event) => void) {
		this.errorListeners.push(callback);
	}

	onMessage(callback: (data: ArrayBuffer | Record<string, any>) => void) {
		this.onMessageListeners.push(callback)
	}

	close(code?: number, reason?: string) {
		if (!this.socket)
			throw new Error('Call socket.connect first')
		this.socket.close(code, reason)
	}

	protected handleMessage(event: MessageEvent) {
		this.onMessageListeners.forEach(async callback => {
			if (event.data instanceof ArrayBuffer) {
				return callback(event.data);
			}

			if (event.data instanceof Blob) {
				const arrayBuffer = await this.blobToArrayBuffer(event.data);
				return callback(arrayBuffer);
			}

			return callback(JSON.parse(event.data));
		})
	}

	private blobToArrayBuffer(blob: Blob) {
		const fileReader = new FileReader();
		const result = new Promise<ArrayBuffer>((resolve, reject) => {
			fileReader.onload = function (event) {
				if (!event.target)
					return reject(new Error('event.target is set null'));
				const res = event.target.result;
				if (!(res instanceof ArrayBuffer))
					return reject(new Error('could not convert to ArrayBuffer'))
				return resolve(res);
			};
		})
		fileReader.readAsArrayBuffer(blob);
		return result
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
		if (!this.socket)
			throw new Error('Call socket.connect first')
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
		if (!this.socket)
			throw new Error('Call socket.connect first')
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