class SocketWrapper {
    constructor(socket, pingInterval) {
        this.socket = socket
        this.eventListeners = []
        this.openListeners = []
        this.pingInterval = pingInterval
        this.ping = null
        this.socket.onmessage = this.triggerHandlers.bind(this)
        this.socket.onopen = event => {
            this.pingPong()
            this.openListeners.forEach(callback => {
                callback(event)
            })
        }

        this.on('pong', (data) => {
            const currentTime = new Date().getTime()
            this.ping = currentTime - data.timestamp
        })

    }

    pingPong() {
        this.emit('ping', {timestamp: new Date().getTime()})
        setTimeout(() => this.pingPong(), this.pingInterval)
    }

    triggerHandlers(event) {
        const data = JSON.parse(event.data)
        this.eventListeners.forEach(listener => {
            if (listener.event === data.event) {
                listener.callback(data.data)
            }
        })
    }

    close(code, reason) {
        this.socket.close(code, reason)
    }

    emit(event, data) {
        if (this.socket.readyState === WebSocket.CLOSED)
            return
        if (this.socket.readyState === WebSocket.CONNECTING)
            return
        this.socket.send(JSON.stringify({event, data}))
    }

    on(event, callback) {
        if (event === 'open')
            this.openListeners.push(callback)
        else
            this.eventListeners.push({event, callback})
    }
}

const Main = {
    template: `
        <div>
            <div>
            Bots: {{ stats.botsCount }}
            Players: {{ stats.playersCount }}
            </div>
        </div>
    `,
    data() {
        return {
            socket: null,
            stats: {}
        }
    },
    created() {
        this.socket = new SocketWrapper(new WebSocket('ws://localhost:3100/admin'), 1000)
        this.socket.on('stats', data => {
            this.stats = data
        })
    }
}
