export class SocketWrapper {
    public readonly socket: WebSocket;
    private eventListeners: { event: string, callback: (data: any) => void }[];
    private openListeners: Function[];
    public ping: number | null;
    public pingInterval: number

    constructor(socket: WebSocket, pingInterval: number) {
        this.socket = socket
        this.eventListeners = []
        this.openListeners = []
        this.pingInterval = pingInterval
        this.ping = null
        this.socket.onmessage = this.triggerHandlers.bind(this)
        this.socket.onopen = (event: Event) => {
            this.pingPong()
            this.openListeners.forEach(callback => {
                callback(event)
            })
        }

        this.on('pong', (data: { timestamp: number }) => {
            const currentTime = new Date().getTime()
            this.ping = currentTime - data.timestamp
        })

    }

    private pingPong() {
        this.emit('ping', {timestamp: new Date().getTime()})
        setTimeout(() => this.pingPong(), this.pingInterval)
    }

    private triggerHandlers(event: MessageEvent) {
        const data = JSON.parse(event.data)
        this.eventListeners.forEach(listener => {
            if (listener.event === data.event) {
                listener.callback(data.data)
            }
        })
    }

    close(code?: number, reason?: string) {
        this.socket.close(code, reason)
    }

    emit(event: string, data: any) {
        if (this.socket.readyState === WebSocket.CLOSED)
            return
        if (this.socket.readyState === WebSocket.CONNECTING)
            return
        this.socket.send(JSON.stringify({event, data}))
    }

    on(event: string, callback: (data: any) => void) {
        if (event === 'open')
            this.openListeners.push(callback)
        else
            this.eventListeners.push({event, callback})
    }
}

export function clipValue(value: number, min: number, max: number) {
    return Math.min(Math.max(value, min), max)
}

export function lightenDarkenColor(color: number[], percent: number) {
    const [R, G, B] = color
    return [
        clipValue(R + percent, 0, 255),
        clipValue(G + percent, 0, 255),
        clipValue(B + percent, 0, 255)
    ]
}

export function calcDistance(x1: number, y1: number, x2: number, y2: number) {
    let dX = x2 - x1
    let dY = y2 - y1
    return Math.sqrt(dX * dX + dY * dY)
}

export function isMobile() {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
}