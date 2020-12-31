import React from 'react';
import Sketch from 'react-p5';
import Game from "./gameUtils/GameEngine";
import CreatePlayerModal from "./components/CreatePlayerModal";
import {SelfPlayerData} from "./gameUtils/player";
import Stats from "./Stats";
import Ping from "./Ping";
import RIP from "./RIP";

let height = window.innerHeight - 10;
let width = window.innerWidth - 10;

class App extends React.Component {
    game?: Game
    state = {
        show: false,
        stats: [],
        socketOpen: false,
        ping: null
    };

    onCreated = (data: SelfPlayerData) => {
        let game = new Game(width, height, data)
        game.socket.on('stats', data => {
            if (!data)
                return
            this.setState({stats: data})
        })

        game.socket.on('open', data => {
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
            <div className="App">
                <Ping ping={this.state.ping}/>
                <CreatePlayerModal onCreated={this.onCreated}/>
                <RIP show={this.state.show}/>
                <Stats stats={this.state.stats}/>
                <Sketch
                    setup={(p5, parentRef) => {
                        p5.createCanvas(width, height).parent(parentRef);
                    }}
                    draw={p5 => {
                        if (!this.game)
                            return
                        this.game.draw(p5)
                    }}
                    windowResized={p5 => {
                        height = window.innerHeight - 10;
                        width = window.innerWidth - 10;
                        p5.resizeCanvas(width, height)
                    }}
                />
            </div>
        );
    }
}

export default App;
