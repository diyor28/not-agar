export class SocketWrapper {
    public readonly socket: WebSocket;
    private eventListeners: { event: string, callback: (data: any) => void }[];

    constructor(socket: WebSocket) {
        this.socket = socket
        this.eventListeners = []
        this.socket.onmessage = this.triggerHandlers.bind(this)
    }

    triggerHandlers(event: MessageEvent) {
        const data = JSON.parse(event.data)
        this.eventListeners.forEach(listener => {
            if (listener.event === data.event) {
                listener.callback(data.data)
            }
        })
    }

    emit(event: string, data: any) {
        if (this.socket.readyState === WebSocket.CLOSED)
            return
        if (this.socket.readyState === WebSocket.CONNECTING)
            return
        this.socket.send(JSON.stringify({event, data}))
    }

    on(event: string, callback: (data: any) => void) {
        this.eventListeners.push({event, callback})
    }
}

export function lightenDarkenColor(color: number[], percent: number) {
    const [R, G, B] = color
    return [R + percent, G + percent, B + percent]
}