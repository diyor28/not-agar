import React from 'react';
import Sketch from 'react-p5';
import Game from "../../engine/GameEngine";
import Stats from "./Stats";
import Ping from "./Ping";
import RIP from "./RIP";
import CreatePlayerModal from "./CreatePlayerModal";
import Tips from "./Tips";
import {EventBus} from "../../client";
import {GameEvent} from "../../client/schemas";

let height = window.innerHeight - 10;
let width = window.innerWidth - 10;

export default class GameCanvas extends React.Component<any, any> {
    game = new Game(width, height);
    eventBus = new EventBus();
    state = {
        show: false,
        stats: [],
        socketOpen: false,
        ping: null
    };

    createPlayer = async (data: { nickname: string }) => {
        await this.game.client.connect();
        await this.game.startGame(data);
        this.game.client.on(GameEvent.StatsUpdate, (data) => {
            this.setState({stats: data.topPlayers});
        });

        this.game.client.on('open', () => {
            this.setState({socketOpen: true});
        });

        this.game.client.on(GameEvent.Pong, ({ping}) => {
            this.setState({ping});
        });

        this.game.client.once(GameEvent.Rip, () => {
            this.setState({show: true});
        });
    }


    render() {
        return (
            <div>
                <CreatePlayerModal eventBus={this.eventBus} createPlayer={this.createPlayer}/>
                {/*<PlayButton eventBus={this.eventBus}/>*/}
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
