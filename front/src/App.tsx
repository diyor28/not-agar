import React from 'react';
import GameCanvas from "./views/Game/GameCanvas";

class App extends React.Component<any, any> {
    render() {
        return (
            <div className="App">
                <GameCanvas/>
            </div>
        );
    }
}

export default App;
