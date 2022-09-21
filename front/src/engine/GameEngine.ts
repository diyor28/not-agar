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
import Player, {SelfPlayer} from "./player";
import {isMobile} from "./utils";
import Food from "./food";
import p5Types from "p5";
import Spike from "./spike"; //Import this for typechecking and intellisense
import JoyStick from "./joystick";
import GameClient, {FoodData, MovedEvent, PlayerData, SpikeData} from "../client";

export type StatsUpdate = {
    weight: number,
    nickname: string
}

export default class Game {
    public players: Player[];
    public food: Food[];
    public spikes: Spike[];
    public stats: StatsUpdate[];
    public client: GameClient;
    public selfPlayer!: SelfPlayer;
    public started = false;
    public _zoom: number;
    private readonly socketUrl: string;
    public width: number;
    public height: number;
    public joystick: JoyStick;

    constructor(width: number, height: number) {
        this.width = width;
        this.height = height;
        this.players = [];
        this.food = [];
        this.stats = [];
        this.spikes = [];
        this._zoom = 1.0;
        this.socketUrl = process.env.REACT_APP_WS_URL as string;
        this.client = new GameClient(this.socketUrl, 1000);
        this.client.onMove(this.onMoved.bind(this));
        this.client.onPlayersUpdate(this.playersUpdated.bind(this));
        this.client.onFoodUpdate(this.foodUpdated.bind(this));
        this.client.onStatsUpdate(this.onStatsUpdate.bind(this));
        this.joystick = new JoyStick(width, height);
    }

    get isMobile() {
        return isMobile();
    }

    get ping() {
        return this.client.ping;
    }

    async startGame({nickname}: { nickname: string }) {
        const {player, spikes} = await this.client.startGame({nickname});
        this.selfPlayer = new SelfPlayer({...player, nickname}, this.height, this.width);
        this.started = true;
        this.initSpikes(spikes);
    }

    windowResized(width: number, height: number) {
        this.width = width;
        this.height = height;
        this.joystick.width = width;
        this.joystick.height = height;
        if (!this.started)
            return
        this.selfPlayer.width = width;
        this.selfPlayer.height = height;
    }

    onStatsUpdate(data: {topPlayers: StatsUpdate[]}) {
        this.stats = data.topPlayers;
    }

    playersUpdated(data: {players: PlayerData[]}) {
        let cameraX = this.selfPlayer._x;
        let cameraY = this.selfPlayer._y;
        this.players = [];
        data.players.forEach((player: PlayerData) => {
            this.players.push(new Player({
                x: player.x,
                y: player.y,
                nickname: player.nickname,
                cameraX,
                cameraY,
                weight: player.weight,
                color: player.color
            }));
        })
    }

    foodUpdated(data: {food: FoodData[]}) {
        let cameraX = this.selfPlayer._x;
        let cameraY = this.selfPlayer._y;
        this.food = [];
        data.food.forEach((food: FoodData) => {
            this.food.push(new Food({
                x: food.x,
                y: food.y,
                cameraX,
                cameraY,
                weight: food.weight,
                color: food.color
            }));
        })
    }

    onMoved(data: MovedEvent) {
        this.selfPlayer.update(data);
        this.spikes.forEach(spike => {
            spike.cameraX = data.x;
            spike.cameraY = data.y;
        })
    }

    emitMove(newX: number, newY: number) {
        const data = {
            newX: this.selfPlayer._x + newX,
            newY: this.selfPlayer._y + newY
        }
        this.client.move(data);
    }

    get zoom() {
        if (this.isMobile)
            return this._zoom * 0.8
        return this._zoom
    }

    updateZoom() {
        this._zoom -= 0.1 * (this._zoom - this.selfPlayer.zoom)
    }

    draw(p5: p5Types) {
        if (!this.started)
            return;
        this.updateZoom()
        p5.background(240);
        p5.translate(this.width / 2, this.height / 2);
        p5.scale(this.zoom);
        this.selfPlayer.draw(p5)
        this.drawAll(p5)
        if (this.isMobile)
            this.joystick.draw(p5)
    }

    touchStarted(p5: p5Types) {
        if (!this.started)
            return;
        this.joystick.touchStarted(p5.mouseX, p5.mouseY)
    }

    touchMoved(p5: p5Types) {
        if (!this.started)
            return;
        const {newX, newY} = this.joystick.move(p5.mouseX, p5.mouseY)
        this.emitMove(newX, newY)
    }

    drawAll(p5: p5Types) {
        this.food.forEach(food => {
            food.draw(p5)
        })
        this.spikes.forEach(spike => {
            spike.draw(p5)
        })
        this.players.forEach((player: Player) => {
            player.draw(p5)
        })
    }

    touchEnded(p5: p5Types) {
        if (!this.started)
            return;
        this.joystick.touchEnded()
    }

    mouseMoved(mouseX: number, mouseY: number) {
        if (!this.started)
            return;
        const {newX, newY} = this.selfPlayer.move(mouseX, mouseY)
        this.emitMove(newX, newY)
    }

    private initSpikes(data: SpikeData[]) {
        let cameraX = this.selfPlayer._x;
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
            }));
        })
    }
}