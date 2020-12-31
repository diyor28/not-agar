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
import Player, {MovedEvent, PlayerData, SelfPlayer, SelfPlayerData} from "./player";
import {SocketWrapper} from "./utils";
import Food, {FoodData} from "./food";
import p5Types from "p5";
import Spike, {SpikeData} from "./spike"; //Import this for typechecking and intellisense

type StatsUpdate = {
    weight: number,
    nickname: string
}

export default class Game {
    public players: Player[];
    public foods: Food[];
    public spikes: Spike[]
    public stats: StatsUpdate[]
    public socket: SocketWrapper;
    public selfPlayer: SelfPlayer
    public zoom: number;
    private readonly socketUrl: string;
    public width: number;
    public height: number;

    constructor(width: number, height: number, data: SelfPlayerData) {
        this.width = width
        this.height = height
        this.players = []
        this.foods = []
        this.stats = []
        this.spikes = []
        this.zoom = 1.0
        // @ts-ignore
        this.socketUrl = process.env.REACT_APP_WS_URL + `/${data.uuid}/`
        this.socket = new SocketWrapper(new WebSocket(this.socketUrl), 1000)
        this.socket.on('moved', data => this.onMoved(data))
        this.socket.on('playersUpdated', data => this.playersUpdated(data))
        this.socket.on('foodUpdated', data => this.foodUpdated(data))
        this.socket.on('spikesUpdated', data => this.spikesUpdated(data))
        this.socket.on('stats', data => this.onStatsUpdate(data))
        this.selfPlayer = new SelfPlayer(this.socket, data, this.height, this.width)
    }

    windowResized(width: number, height: number) {
        this.width = width
        this.height = height
        this.selfPlayer.width = width
        this.selfPlayer.height = height
    }

    onStatsUpdate(data: StatsUpdate[]) {
        this.stats = data
    }

    get ping() {
        return this.socket.ping
    }

    accelerate() {
        this.socket.emit('accelerate', {uuid: this.selfPlayer.uuid})
    }

    playersUpdated(data: PlayerData[]) {
        if (!this.selfPlayer)
            return
        let cameraX = this.selfPlayer._x
        let cameraY = this.selfPlayer._y;
        let players = data || [];
        this.players = [];
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
        // console.log(this.players.length, this.foods)
    }

    spikesUpdated(data: PlayerData[]) {
        if (!this.selfPlayer)
            return
        let cameraX = this.selfPlayer._x
        let cameraY = this.selfPlayer._y;
        let spikes = data || [];
        this.spikes = [];
        spikes.forEach((spike: SpikeData) => {
            this.spikes.push(new Spike({
                x: spike.x,
                y: spike.y,
                cameraX,
                cameraY,
                weight: spike.weight,
                color: [0, 255, 0]
            }))
        })
        // console.log(this.players.length, this.foods)
    }

    foodUpdated(data: FoodData[]) {
        if (!this.selfPlayer)
            return
        let cameraX = this.selfPlayer._x
        let cameraY = this.selfPlayer._y;
        let foods = data || [];
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
    }

    onMoved(data: MovedEvent) {
        this.selfPlayer.update(data)
    }

    draw(p5: p5Types) {
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
            player.draw(p5)
        })
        this.foods.forEach(food => {
            food.draw(p5)
        })
        this.spikes.forEach(spike => {
            spike.draw(p5)
        })
    }
}