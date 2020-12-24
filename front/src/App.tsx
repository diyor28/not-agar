import React from 'react';
import Sketch from 'react-p5';
import Game from "./gameUtils/GameEngine";
import CreatePlayerModal from "./components/CreatePlayerModal";
import {SelfPlayerData} from "./gameUtils/player";
import Stats from "./Stats";

const height = window.innerHeight - 10;
const width = window.innerWidth - 10;

class App extends React.Component {
    game = new Game(width, height)
    state = {
        show: false,
        stats: [],
        socketOpen: false
    };

    componentDidMount() {
        this.game.socket.on('stats', data => {
            if (!data)
                return
            this.setState({stats: data})
        })

        this.game.socket.onopen(data => {
            this.setState({socketOpen: true})
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
