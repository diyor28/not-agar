import React from 'react';
import Sketch from 'react-p5';
import Game from "./gameUtils/GameEngine";
import CreatePlayerModal from "./components/CreatePlayerModal";
import {SelfPlayerData} from "./gameUtils/player";
import Stats from "./Stats";
import Ping from "./Ping";

const height = window.innerHeight - 10;
const width = window.innerWidth - 10;

class App extends React.Component {
    game = new Game(width, height)
    state = {
        show: false,
        stats: [],
        socketOpen: false,
        ping: null
    };

    componentDidMount() {
        this.game.socket.on('stats', data => {
            if (!data)
                return
            this.setState({stats: data})
        })

        this.game.socket.on('open', data => {
            this.setState({socketOpen: true})
        })

        this.game.socket.on('pong', () => {
            this.setState({ping: this.game.ping})
        })

        window.addEventListener('keydown', (event: KeyboardEvent) => {
            if (event.key !== 'a') {
                return
            }
            this.game.accelerate()
        })
    }

    onCreated = (data: SelfPlayerData) => {
        this.game.playerCreated(data)
    }

    render() {
        return (
            <div className="App">
                <Ping ping={this.state.ping}/>
                <CreatePlayerModal socketOpen={this.state.socketOpen} onCreated={this.onCreated}/>
                <Stats stats={this.state.stats}/>
                <Sketch
                    setup={(p5, parentRef) => {
                        p5.createCanvas(width, height).parent(parentRef);
                    }}
                    draw={p5 => {
                        this.game.draw(p5)
                    }}
                />
            </div>
        );
    }
}

export default App;
