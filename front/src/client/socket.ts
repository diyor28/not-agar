import {EventBus} from "./eventBus";

export class SocketWrapper {
	socket?: WebSocket;
	url: string
	protected bus: EventBus;

	constructor(url: string) {
		this.url = url;
		this.bus = new EventBus();
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

	async connect(): Promise<Event> {
		this.socket = new WebSocket(this.url);
		const connectPromise = new Promise<Event>((resolve, reject) => {
			this.once('open', resolve);
			this.once('error', reject);
		});
		this.socket.onopen = (event: Event) => {
			this.bus.emit('open', event);
		}
		this.socket.onerror = (event: Event) => {
			this.bus.emit('error', event);
		}
		this.socket.onmessage = this.handleMessage.bind(this);
		return connectPromise
	}

	once(event: 'open', callback: (data: Event) => void): void
	once(event: 'error', callback: (data: Event) => void): void
	once(event: 'open' | 'error', callback: (data: Event) => void): void {
		this.bus.once(event, callback);
	}

	on(event: 'open', callback: (data: Event) => void): void
	on(event: 'error', callback: (data: Event) => void): void
	on(event: 'message', callback: (data: Buffer) => void): void
	on(event: 'open' | 'error' | 'message', callback: ((data: Event) => void) | ((data: Buffer) => void)): void {
		this.bus.on(event, callback);
	}

	close(code?: number, reason?: string) {
		if (!this.socket)
			throw new Error('Call socket.connect first')
		this.socket.close(code, reason)
	}

	protected async handleMessage(event: MessageEvent<Buffer | ArrayBuffer | Blob>) {
		if (event.data instanceof ArrayBuffer) {
			return this.bus.emit('message', Buffer.from(event.data));
		}

		if (event.data instanceof Blob) {
			const arrayBuffer = await this.blobToArrayBuffer(event.data);
			return this.bus.emit('message', Buffer.from(arrayBuffer));
		}

		return this.bus.emit('message', event.data);
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
