import Entity, {EntityData} from "./entity";
import p5Types from "p5"; //Import this for typechecking and intellisense
import {lightenDarkenColor, SocketWrapper} from './utils'

const strokeWeight = 8
const textColor = 255

export type PlayerData = EntityData & { nickname: string }


export type SelfPlayerData = {
    uuid: string,
    x: number,
    y: number,
    nickname: string,
    weight: number,
    speed: number,
    color: number[],
    zoom: number
}

export default class Player extends Entity {
    public nickname;

    constructor({x, y, cameraX, cameraY, nickname, weight, color}: PlayerData) {
        super({x, y, cameraX, cameraY, weight, color})
        this.nickname = nickname
    }

    draw(p5: p5Types) {
        p5.fill(this.color)
        p5.strokeWeight(strokeWeight)
        p5.stroke(lightenDarkenColor(this.color, -20))
        p5.ellipse(this.x, this.y, this.weight, this.weight)
        p5.noStroke()
        p5.fill(textColor)
        p5.textAlign(p5.CENTER, p5.CENTER)
        p5.textSize(this.weight / 5)
        p5.text(this.nickname, this.x, this.y)
    }
}

export class SelfPlayer extends Player {
    public uuid: string;
    public speed: number;
    public socket: SocketWrapper;
    public height: number | undefined;
    public width: number | undefined;
    public zoom: number;

    constructor(socket: SocketWrapper, {uuid, x, y, nickname, weight, speed, color, zoom}: SelfPlayerData) {
        super({x, y, nickname, cameraX: x, cameraY: y, weight, color});
        this.uuid = uuid;
        this.speed = speed
        this.socket = socket
        this.zoom = zoom
    }

    update(data: SelfPlayerData) {
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

    move(mouseX: number, mouseY: number) {
        if (!this.width || !this.height)
            return;
        const newX = this._x + mouseX - this.width / 2
        const newY = this._y + mouseY - this.height / 2

        const data = {
            uuid: this.uuid,
            newX: newX,
            newY: newY
        }
        this.socket.emit('move', data)
    }
}
