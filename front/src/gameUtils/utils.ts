export class SocketWrapper {
    public readonly socket: WebSocket;
    private eventListeners: { event: string, callback: (data: any) => void }[];
    private onOpenListeners: Function[]

    constructor(socket: WebSocket) {
        this.socket = socket
        this.eventListeners = []
        this.onOpenListeners = []
        this.socket.onmessage = this.triggerHandlers.bind(this)
        this.socket.onopen = (event: Event) => {
            this.onOpenListeners.forEach(callback => {
                callback()
            })
        }
    }

    onopen(callback: (data: any) => void) {
        this.onOpenListeners.push(callback)
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

export function calcDistance(x1: number, y1: number, x2: number, y2: number) {
    let dX = x2 - x1
    let dY = y2 - y1
    return Math.sqrt(dX * dX + dY * dY)
}