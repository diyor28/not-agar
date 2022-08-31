import Entity from "./entity";

import p5Types from "p5";
import {lightenDarkenColor} from "./utils"; //Import this for typechecking and intellisense

export default class Spike extends Entity {
    draw(p5: p5Types) {
        p5.fill(this.color)
        p5.strokeWeight(8)
        p5.stroke(lightenDarkenColor(this.color, -20))
        // p5.ellipse(this.x, this.y, this.weight, this.weight)
        p5.beginShape()
        const points = 50
        for (let i = 0; i <= points; i++) {
            let angle = p5.PI * 2 * i / points
            let rX = this.weight * p5.sin(angle) / 2
            let rY = this.weight * p5.cos(angle) / 2
            if (i % 2) {
                rX *= 1.05
                rY *= 1.05
            }
            let x = this._x - this.cameraX + rX
            let y = this._y - this.cameraY + rY
            // p5.point(x, y)
            p5.vertex(x, y)
        }
        p5.endShape()
    }
}
