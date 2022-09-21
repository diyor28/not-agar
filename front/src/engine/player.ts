import Entity from "./entity";
import p5Types from "p5"; //Import this for typechecking and intellisense
import {lightenDarkenColor} from './utils'
import {MovedEvent, PlayerData, SelfPlayerData} from "../client";

const strokeWeight = 4
const textColor = 255

export default class Player extends Entity {
    public nickname;

    constructor({x, y, cameraX, cameraY, nickname, weight, color}: PlayerData) {
        super({x, y, cameraX, cameraY, weight, color})
        this.nickname = nickname
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
    }
}

export class SelfPlayer extends Player {
    uuid: string;
    height: number;
    width: number;
    zoom: number;
    velocityX = 0;
    velocityY = 0;
    points: { x: number, y: number }[]

    constructor(data: SelfPlayerData, height: number, width: number) {
        super({
            x: data.x,
            y: data.y,
            nickname: data.nickname,
            cameraX: data.x,
            cameraY: data.y,
            weight: data.weight,
            color: data.color
        });
        this.uuid = data.uuid;
        this.zoom = data.zoom
        this.height = height
        this.width = width
        this.points = data.points;
    }

    update(data: MovedEvent) {
        this._x = data.x;
        this._y = data.y;
        this.weight = data.weight;
        this.zoom = data.zoom;
        this.velocityX = data.velocityX;
        this.velocityY = data.velocityY;
        this.points = data.points;
    }

    get x() {
        return 0
    }

    get y() {
        return 0
    }

    drawPerimeter(p5: p5Types) {
        p5.beginShape()
        for (let i = 0; i < this.points.length; i ++) {
            p5.vertex(this.points[i].x, this.points[i].y)
        }
        p5.vertex(this.points[0].x, this.points[0].y)
        p5.endShape()
    }

    draw(p5: p5Types) {
        let ringColor = lightenDarkenColor(this.color, - 20)
        p5.fill(this.color)
        p5.strokeWeight(strokeWeight)
        p5.stroke(ringColor)
        this.drawPerimeter(p5)
        p5.noStroke()
        p5.fill(textColor)
        p5.textAlign(p5.CENTER, p5.CENTER)
        p5.textSize(this.weight / 5)
        p5.text(this.nickname, this.x, this.y)
    }

    move(mouseX: number, mouseY: number) {
        const newX = mouseX - this.width / 2
        const newY = mouseY - this.height / 2
        return {newX, newY}
    }
}
