const height = window.innerHeight - 10;
const width = window.innerWidth - 10;
let game;
let grid;
let zoom = 1.0;

class Grid {
    constructor() {
        this.gridPoints = []
        this.cameraX = null
        this.cameraY = null
        this.gridPoints = []
        this.step = width / 20
    }

    setCamera(cameraX, cameraY) {
        const step = width / 20
        this.cameraX = cameraX
        this.cameraY = cameraY
        let newPoints = []
        this.gridPoints.forEach((p, index) => {
            const distance = Math.sqrt(Math.pow(p.x - this.cameraX, 2) + Math.pow(p.y - this.cameraY, 2))
            // if (distance < Math.min(height, width) / 2) {
            //     newPoints.push(p)
            // }
        })
        for (let range = 0; range < Math.max(width, height); range += this.step) {
            // this.gridPoints.push({
            //     x: this.cameraX + range,
            //     y: this.cameraY + range
            // })
        }
    }

    draw() {
        if (!this.cameraX)
            return
        if (!this.cameraY)
            return
        this.gridPoints.forEach(point => {
            const x = point.x - this.cameraX
            const y = point.y - this.cameraY
            stroke(200);
            strokeWeight(1);
            line(x, 0, x, height);
            line(0, y, width, y);
        })
    }
}

class Entity {
    constructor({x, y, cameraX, cameraY, weight, color}) {
        this._x = x;
        this._y = y;
        this.cameraX = cameraX;
        this.cameraY = cameraY;
        this.weight = weight;
        this.color = color;
    }

    get dist() {
        return Math.sqrt(Math.pow(this._x - this.cameraX, 2) + Math.pow(this._y - this.cameraY, 2))
    }

    get x() {
        return this._x - this.cameraX
    }

    get y() {
        return this._y - this.cameraY
    }

    draw() {
        fill(this.color)
        ellipse(this.x, this.y, this.weight, this.weight)
    }

}

class Food extends Entity {
}

class Player extends Entity {
    constructor({x, y, cameraX, cameraY, nickname, weight, color}) {
        super({x, y, cameraX, cameraY, weight, color})
        this.nickname = nickname
    }

    draw() {
        if (!this.color || !this.weight)
            return
        fill(this.color)
        ellipse(this.x, this.y, this.weight, this.weight)
        console.log('###nickanme', this.nickname)
        text(this.nickname, this.x, this.y)
    }
}

class SelfPlayer extends Player {
    constructor(socket, {uuid, x, y, nickname, weight, speed, color, zoom}) {
        super({x, y, nickname, cameraX: x, cameraY: y, weight, color});
        this.uuid = uuid;
        this.speed = speed
        this.socket = socket
        this.zoom = zoom
    }

    update(data) {
        this._x = data.x;
        this._y = data.y;
        this.weight = data.weight;
        this.speed = data.speed;
        this.zoom = data.zoom;
        this.nickname = data.nickname;
        this.color = data.color;
    }

    get x() {
        return 0
    }

    get y() {
        return 0
    }

    move(mouseX, mouseY) {
        const diffX = mouseX - (width / 2)
        const diffY = mouseY - (height / 2)
        const directionX = diffX ? diffX / Math.abs(diffX) : 0
        const directionY = diffY ? diffY / Math.abs(diffY) : 0
        if (this.socket.readyState === WebSocket.CLOSED)
            return
        if (this.socket.readyState === WebSocket.CONNECTING)
            return
        const data = {
            uuid: this.uuid,
            directionX,
            directionY
        }
        this.socket.send(JSON.stringify(data))
    }
}


class Game {
    constructor() {
        this.players = []
        this.grid = new Grid()
        this.foods = []
        this.socket = new WebSocket('ws://localhost:3000/ws')
        this.selfPlayer = null
        this.playerCreated = false
        // this.selfPlayer = new SelfPlayer(this.socket, {})
        this.zoom = 1.0
        this.socket.onmessage = (event) => {
            this.onmessage(event)
        }
    }

    onmessage(event) {
        const data = JSON.parse(event.data)
        const selfPlayer = data.selfPlayer
        this.selfPlayer.uuid = data.selfPlayer.uuid
        this.selfPlayer.update(selfPlayer)
        this.zoom = selfPlayer.zoom
        let cameraX = selfPlayer.x
        let cameraY = selfPlayer.y;
        let players = data.players || [];
        let foods = data.foods || [];
        this.players = [];
        this.foods = [];
        foods.forEach(food => {
            this.foods.push(new Food({
                x: food.x,
                y: food.y,
                cameraX,
                cameraY,
                weight: food.weight,
                color: food.color
            }))
        })
        players.forEach(player => {
            this.players.push(new Player({
                x: player.x,
                y: player.y,
                nickname: player.nickname,
                cameraX,
                cameraY,
                weight: player.weight,
                color: player.color
            }))
        })
    }

    animate(mouseX, mouseY) {
        if (!this.playerCreated)
            return
        this.selfPlayer.move(mouseX, mouseY)
        this.selfPlayer.draw()
        this.grid.setCamera(this.selfPlayer.cameraX, this.selfPlayer.cameraY)
        this.grid.draw()
        this.drawAll()
    }

    drawAll() {
        this.players.forEach(player => {
            if (player.uuid === this.selfPlayer.uuid)
                return
            player.draw()
        })
        this.foods.forEach(food => {
            food.draw()
        })
    }
}


function setup() {
    let canvas = createCanvas(width, height);
    canvas.parent('#canvas')
    game = new Game()
}

function draw() {
    translate(width / 2, height / 2);
    zoom -= 0.1 * (zoom - game.zoom)
    scale(zoom);
    background(240);
    game.animate(mouseX, mouseY);
}

const Main = {
    template: `
    <div>
        <div id="canvas">
        </div>
    </div>
    `,
    mounted() {
    }
}
