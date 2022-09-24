import p5Types from "p5";
import {EntityData} from "../client"; //Import this for typechecking and intellisense

export default class Entity {
    _x: number;
    _y: number;
    public weight: number;
    public color: number[];

    constructor({x, y, weight, color}: EntityData) {
        this._x = x;
        this._y = y;
        this.weight = weight;
        this.color = color;
    }

    dist(cameraX: number, cameraY: number) {
        return Math.sqrt(Math.pow(this._x - cameraX, 2) + Math.pow(this._y - cameraY, 2))
    }

    draw(p5: p5Types, cameraX: number, cameraY: number) {
        p5.fill(this.color)
        p5.ellipse(this._x - cameraX, this._y - cameraY, this.weight, this.weight)
    }
}
