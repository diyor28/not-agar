import React from 'react';
import Sketch from 'react-p5';
import Game from "../../engine/GameEngine";
import Stats from "./Stats";
import Ping from "./Ping";
import RIP from "./RIP";
import CreatePlayerModal from "./CreatePlayerModal";
import Tips from "./Tips";
import {Schema} from "../../codec";

let height = window.innerHeight - 10;
let width = window.innerWidth - 10;
// const frameRate = 40

export default class GameCanvas extends React.Component {
    game = new Game(width, height);
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

    onFill = async (data: { nickname: string }) => {
        await this.game.startGame(data);
        this.game.client.onStatsUpdate((data) => {
            this.setState({stats: data});
        });

        this.game.client.socket.onOpen(() => {
            this.setState({socketOpen: true})
        });

        this.game.client.onPong(({ping}) => {
            this.setState({ping});
        });

        this.game.client.onRip(() => {
            // game.socket.close()
            this.setState({show: true})
        });
    }


    render() {
        const schema = new Schema({
            food: 'int16',
            player: {
                x: 'uint8',
                y: 'uint8',
                z: 'uint8'
            },
            spikes: 'int32'
        })
        const buffer = schema.encode({food: 10, player: {x: 5, y: 8}, spikes: 20})
        console.log(buffer.toString('hex'));
        console.log(schema.decode(buffer));
        return (
            <div>
                <CreatePlayerModal onFill={this.onFill}/>
                <Ping ping={this.state.ping}/>
                <RIP show={this.state.show}/>
                <Tips/>
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
