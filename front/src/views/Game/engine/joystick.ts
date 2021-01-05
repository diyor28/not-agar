import p5Types from "p5";

export default class JoyStick {
    public width: number;
    public height: number;
    public _touched: boolean = false;
    public middleX: number = 0;
    public middleY: number = 0;
    public _dX: number = 0;
    public _dY: number = 0;
    public r: number
    public R: number

    constructor(width: number, height: number) {
        this.width = width
        this.height = height
        this.r = (this.width + this.height) / 15
        this.R = this.r * 2.0
    }

    touchStarted(x: number, y: number) {
        this.middleX = x
        this.middleY = y
        this._touched = true
    }

    touchEnded() {
        this._dX = 0
        this._dY = 0
        this._touched = false
    }

    get x() {
        return this.translateX(this.middleX + this._dX)
    }

    get y() {
        return this.translateY(this.middleY + this._dY)
    }

    getDistFromCenter(x: number, y: number) {
        let midX = this.middleX
        let midY = this.middleY
        let vX = x - midX
        let vY = y - midY
        return Math.sqrt(vX * vX + vY * vY)
    }

    setXY(x: number, y: number) {
        const maxDist = this.R / 2
        this._dX = x - this.middleX
        this._dY = y - this.middleY
        let dist = this.getDistFromCenter(x, y)
        if (dist <= maxDist)
            return
        // let aX = midX + (vX / dist * maxDist);
        // let aY = midY + (vY / dist * maxDist);
        this._dX = (this._dX / dist * maxDist)
        this._dY = (this._dY / dist * maxDist)
    }

    move(mouseX: number, mouseY: number) {
        this.setXY(mouseX, mouseY)
        return {newX: this.x - this.translateX(this.middleX), newY: this.y - this.translateY(this.middleY)}
    }

    translateX(x: number, reverse?: boolean) {
        if (reverse)
            return x + this.width / 2
        return x - this.width / 2
    }

    translateY(y: number, reverse?: boolean) {
        if (reverse)
            return y + this.height / 2
        return y - this.height / 2
    }

    draw(p5: p5Types) {
        if (!this._touched)
            return
        p5.strokeWeight(0)
        p5.fill(230, 200)
        p5.ellipse(this.translateX(this.middleX), this.translateY(this.middleY), this.R, this.R)
        p5.fill(200, 200)
        p5.ellipse(this.x, this.y, this.r, this.r)
    }
}