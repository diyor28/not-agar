import Entity, {EntityData} from "./entity";
import p5Types from "p5"; //Import this for typechecking and intellisense
import {lightenDarkenColor, SocketWrapper} from './utils'

const strokeWeight = 8
const textColor = 255

export type PlayerData = EntityData & { nickname: string }
export type MovedEvent = {
    x: number,
    y: number,
    zoom: number,
    weight: number
}


export type SelfPlayerData = {
    uuid: string,
    x: number,
    y: number,
    nickname: string,
    weight: number,
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
    public socket: SocketWrapper;
    public height: number;
    public width: number;
    public zoom: number;

    constructor(socket: SocketWrapper, {uuid, x, y, nickname, weight, color, zoom}: SelfPlayerData, height: number, width: number) {
        super({x, y, nickname, cameraX: x, cameraY: y, weight, color});
        this.uuid = uuid;
        this.socket = socket
        this.zoom = zoom
        this.height = height
        this.width = width
    }

    update(data: MovedEvent) {
        this._x = data.x;
        this._y = data.y;
        this.weight = data.weight;
        this.zoom = data.zoom;
    }

    get x() {
        return 0
    }

    get y() {
        return 0
    }

    move(mouseX: number, mouseY: number) {
        const newX = mouseX - this.width / 2
        const newY = mouseY - this.height / 2
        return {newX, newY}
    }
}
