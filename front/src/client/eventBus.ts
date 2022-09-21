export class EventBus {
	listeners: { event: string, callback: (data: any) => void, once: boolean }[] = []

	constructor() {
	}

	off(event: string, callback: (data: any) => void) {
		const indices: number[] = [];
		this.listeners.forEach((listener, index) => {
			if (listener.event === event && listener.callback === callback) {
				indices.push(index);
			}
		});
		for (let i = indices.length; i >= 0; i --) {
			this.listeners.splice(indices[i], 1);
		}
	}

	on(event: string, callback: (data: any) => void) {
		this.listeners.push({event, callback, once: false});
	}

	once(event: string, callback: (data: any) => void) {
		this.listeners.push({event, callback, once: true});
	}

	emit(event: string, data: any) {
		this.listeners.forEach((listener, index) => {
			if (listener.event === event) {
				listener.callback(data);
				if (listener.once) {
					this.off(listener.event, listener.callback);
				}
			}
		})
	}
}
