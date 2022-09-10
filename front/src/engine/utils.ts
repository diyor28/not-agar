export function clipValue(value: number, min: number, max: number) {
    return Math.min(Math.max(value, min), max)
}

export function lightenDarkenColor(color: number[], percent: number) {
    const [R, G, B] = color
    return [
        clipValue(R + percent, 0, 255),
        clipValue(G + percent, 0, 255),
        clipValue(B + percent, 0, 255)
    ]
}

export function calcDistance(x1: number, y1: number, x2: number, y2: number) {
    let dX = x2 - x1
    let dY = y2 - y1
    return Math.sqrt(dX * dX + dY * dY)
}

export function isMobile() {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
}