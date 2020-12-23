import React from 'react';
import Sketch from 'react-p5';
import Game from "./gameUtils/GameEngine";
import CreatePlayerModal from "./components/CreatePlayerModal";
import p5Types from "p5"; //Import this for typechecking and intellisense
import {SelfPlayerData} from "./gameUtils/player";
import Stats from "./Stats";

const height = window.innerHeight - 10;
const width = window.innerWidth - 10;

class App extends React.Component {
    game = new Game(width, height)
    state = {
        show: false,
        stats: []
    };

    componentDidMount() {
        this.game.socket.on('stats', data => {
            if (!data)
                return
            this.setState({stats: data})
        })
    }

    onCreated = (data: SelfPlayerData) => {
        this.game.playerCreated(data)
    }

    render() {
        return (
            <div className="App">
                <CreatePlayerModal onCreated={this.onCreated}/>
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
