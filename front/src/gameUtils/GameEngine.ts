// class Grid {
//     constructor() {
//         this.gridPoints = []
//         this.cameraX = null
//         this.cameraY = null
//         this.gridPoints = []
//         this.step = width / 20
//     }
//
//     setCamera(cameraX, cameraY) {
//         const step = width / 20
//         this.cameraX = cameraX
//         this.cameraY = cameraY
//         let newPoints = []
//         this.gridPoints.forEach((p, index) => {
//             const distance = Math.sqrt(Math.pow(p.x - this.cameraX, 2) + Math.pow(p.y - this.cameraY, 2))
//             // if (distance < Math.min(height, width) / 2) {
//             //     newPoints.push(p)
//             // }
//         })
//         for (let range = 0; range < Math.max(width, height); range += this.step) {
//             // this.gridPoints.push({
//             //     x: this.cameraX + range,
//             //     y: this.cameraY + range
//             // })
//         }
//     }
//
//     draw() {
//         if (!this.cameraX)
//             return
//         if (!this.cameraY)
//             return
//         this.gridPoints.forEach(point => {
//             const x = point.x - this.cameraX
//             const y = point.y - this.cameraY
//             stroke(200);
//             strokeWeight(1);
//             line(x, 0, x, height);
//             line(0, y, width, y);
//         })
//     }
// }
import Player, {PlayerData, SelfPlayer, SelfPlayerData} from "./player";
import {SocketWrapper} from "./utils";
import Food, {FoodData} from "./food";
import p5Types from "p5"; //Import this for typechecking and intellisense

type StatsUpdate = {
    weight: number,
    nickname: string
}

export default class Game {
    public players: Player[];
    public foods: Food[];
    public stats: StatsUpdate[]
    public socket: SocketWrapper;
    public selfPlayer?: SelfPlayer
    public zoom: number;
    private readonly socketUrl: string;
    public readonly width: number;
    public readonly height: number;

    constructor(width: number, height: number) {
        this.width = width
        this.height = height
        this.players = []
        this.foods = []
        this.stats = []
        // @ts-ignore
        this.socketUrl = process.env.REACT_APP_WS_URL
        console.log('###created socket')
        this.socket = new SocketWrapper(new WebSocket(this.socketUrl), 1000)
        this.zoom = 1.0
        this.socket.on('moved', data => this.onMoved(data))
        this.socket.on('stats', data => this.onStatsUpdate(data))
    }

    playerCreated({uuid, x, y, nickname, weight, speed, color, zoom}: SelfPlayerData) {
        this.selfPlayer = new SelfPlayer(this.socket, {uuid, x, y, nickname, weight, speed, color, zoom})
        this.selfPlayer.width = this.width
        this.selfPlayer.height = this.height
    }

    onStatsUpdate(data: StatsUpdate[]) {
        this.stats = data
    }

    get ping() {
        return this.socket.ping
    }

    accelerate() {
        if (!this.selfPlayer)
            return
        this.socket.emit('accelerate', {uuid: this.selfPlayer.uuid})
    }

    onMoved(data: { selfPlayer: SelfPlayerData, players: PlayerData[], foods: FoodData[] }) {
        const selfPlayer = data.selfPlayer
        if (!this.selfPlayer)
            return
        this.selfPlayer.update(selfPlayer)
        let cameraX = selfPlayer.x
        let cameraY = selfPlayer.y;
        let players = data.players || [];
        let foods = data.foods || [];
        this.players = [];
        this.foods = [];
        foods.forEach((food: FoodData) => {
            this.foods.push(new Food({
                x: food.x,
                y: food.y,
                cameraX,
                cameraY,
                weight: food.weight,
                color: food.color
            }))
        })
        players.forEach((player: PlayerData) => {
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

    draw(p5: p5Types) {
        if (!this.selfPlayer)
            return
        p5.background(240);
        p5.translate(this.width / 2, this.height / 2);
        this.zoom -= 0.1 * (this.zoom - this.selfPlayer.zoom)
        p5.scale(this.zoom);
        this.selfPlayer.move(p5.mouseX, p5.mouseY)
        this.selfPlayer.draw(p5)
        this.drawAll(p5)
    }

    drawAll(p5: p5Types) {
        this.players.forEach((player: Player) => {
            // @ts-ignore
            if (player.uuid === this.selfPlayer.uuid)
                return
            player.draw(p5)
        })
        this.foods.forEach(food => {
            food.draw(p5)
        })
    }
}

