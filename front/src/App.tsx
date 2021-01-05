import React from 'react';
import GameCanvas from "./views/Game/GameCanvas";

class App extends React.Component {
    render() {
        return (
            <div className="App">
                <GameCanvas/>
            </div>
        );
    }
}

export default App;
