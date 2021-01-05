import React from 'react';
import Sketch from 'react-p5';
import Game from "./engine/GameEngine";
import {SelfPlayerData} from "./engine/player";
import {SpikeData} from "./engine/spike";
import Stats from "./Stats";
import Ping from "./Ping";
import RIP from "./RIP";
import CreatePlayerModal from "./CreatePlayerModal";

let height = window.innerHeight - 10;
let width = window.innerWidth - 10;
// const frameRate = 40

export default class GameCanvas extends React.Component {
    game?: Game
    // game = new Game(width, height, {
    //     uuid: "2423423",
    //     color: [255, 0, 255],
    //     nickname: "",
    //     weight: 40,
    //     x: 3000,
    //     y: 3000,
    //     zoom: 1.0
    // })
    state = {
        show: false,
        stats: [],
        socketOpen: false,
        ping: null
    };

    onCreated = (data: { player: SelfPlayerData, spikes: SpikeData[] }) => {
        let game = new Game(width, height, data.player)
        game.spikesCreated(data.spikes)
        game.socket.on('stats', (data) => {
            this.setState({stats: data})
        })

        game.socket.on('open', () => {
            this.setState({socketOpen: true})
        })

        game.socket.on('pong', () => {
            this.setState({ping: game.ping})
        })

        game.socket.on('rip', () => {
            game.socket.close()
            this.setState({show: true})
        })

        window.addEventListener('keydown', (event: KeyboardEvent) => {
            if (event.key !== 'a') {
                return
            }
            game.accelerate()
        })
        this.game = game
    }

    render() {
        return (
            <div>
                <CreatePlayerModal onCreated={this.onCreated}/>
                <Ping ping={this.state.ping}/>
                <RIP show={this.state.show}/>
                <Stats stats={this.state.stats}/>
                <Sketch
                    setup={(p5, parentRef) => {
                        p5.createCanvas(width, height).parent(parentRef);
                        // p5.frameRate(frameRate);
                    }}
                    draw={p5 => {
                        if (!this.game)
                            return
                        this.game.draw(p5)
                    }}
                    touchStarted={p5 => {
                        if (!this.game)
                            return
                        this.game.touchStarted(p5)
                    }}
                    touchMoved={p5 => {
                        if (!this.game)
                            return
                        this.game.touchMoved(p5)
                    }}
                    touchEnded={p5 => {
                        if (!this.game)
                            return
                        this.game.touchEnded(p5)
                    }}
                    mouseMoved={p5 => {
                        if (!this.game)
                            return
                        this.game.mouseMoved(p5.mouseX, p5.mouseY)
                    }}
                    windowResized={p5 => {
                        height = window.innerHeight - 10;
                        width = window.innerWidth - 10;
                        p5.resizeCanvas(width, height)
                        if (!this.game)
                            return
                        this.game.windowResized(width, height)
                    }}
                />
            </div>
        );
    }
}
