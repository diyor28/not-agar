import p5Types from "p5"; //Import this for typechecking and intellisense

export type EntityData = {
    x: number,
    y: number,
    cameraX: number,
    cameraY: number,
    weight: number,
    color: number[]
}


export default class Entity {
    _x: number;
    _y: number;
    public cameraX: number;
    public cameraY: number;
    public weight: number;
    public color: number[];

    constructor({x, y, cameraX, cameraY, weight, color}: EntityData) {
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

    draw(p5: p5Types) {
        p5.fill(this.color)
        p5.ellipse(this.x, this.y, this.weight, this.weight)
    }

}
