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
    public cannonAngle;

    constructor({x, y, cameraX, cameraY, nickname, weight, color}: PlayerData) {
        super({x, y, cameraX, cameraY, weight, color})
        this.nickname = nickname
        this.cannonAngle = 90
    }

    drawCannon(p5: p5Types, ringColor: number[]) {
        let x = this.x
        let y = this.y
        let r = this.weight / 2
        let a = Math.PI * this.cannonAngle / 180
        let width = 20
        let height = 30
        let cornerRadius = 5
        // p5.translate(this.x, this.y)
        y += Math.sin(a) * r
        x += Math.cos(a) * r
        console.log(x, y)
        p5.rotate(a * - 1)
        p5.fill(ringColor)
        p5.rect(x, y, width, height, cornerRadius)
        p5.rotate(a)
    }

    draw(p5: p5Types) {
        let ringColor = lightenDarkenColor(this.color, - 20)
        p5.fill(this.color)
        p5.strokeWeight(strokeWeight)
        p5.stroke(ringColor)
        p5.ellipse(this.x, this.y, this.weight, this.weight)
        p5.noStroke()
        p5.fill(textColor)
        p5.textAlign(p5.CENTER, p5.CENTER)
        p5.textSize(this.weight / 5)
        p5.text(this.nickname, this.x, this.y)
        this.drawCannon(p5, ringColor)
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
        window.addEventListener('keydown', event => {
            if (event.code === 'ArrowUp') {
                this.cannonAngle += 10
            }
            if (event.code === 'ArrowDown') {
                this.cannonAngle -= 10
            }
        })
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
