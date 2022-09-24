import Entity from "./entity";
import p5Types from "p5"; //Import this for typechecking and intellisense
import {lightenDarkenColor} from './utils'
import {MovedEvent, PlayerData, SelfPlayerData} from "../client";

const strokeWeight = 4
const textColor = 255

export default class Player extends Entity {
    public nickname;

    constructor({x, y, nickname, weight, color}: PlayerData) {
        super({x, y, weight, color});
        this.nickname = nickname;
    }

    draw(p5: p5Types, cameraX: number, cameraY: number) {
        const x = this._x - cameraX;
        const y = this._y = cameraY;
        let ringColor = lightenDarkenColor(this.color, - 20);
        p5.fill(this.color);
        p5.strokeWeight(strokeWeight);
        p5.stroke(ringColor);
        p5.ellipse(x, y, this.weight, this.weight);
        p5.noStroke();
        p5.fill(textColor);
        p5.textAlign(p5.CENTER, p5.CENTER);
        p5.textSize(this.weight / 5);
        p5.text(this.nickname, x, y);
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
        p5.text(this.nickname, 0, 0)
    }

    move(mouseX: number, mouseY: number) {
        const newX = mouseX - this.width / 2
        const newY = mouseY - this.height / 2
        return {newX, newY}
    }
}
