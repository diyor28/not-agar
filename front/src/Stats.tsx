import React from 'react';
import "./components/modal.css";
import './App.css'

export interface Props {
    stats: { nickname: string, weight: number }[]
}

export default class Stats extends React.Component<Props, {}> {
    state = {
        stats: []
    };

    handleChange = (event: any) => {
        this.setState({nickname: event.target.value})
    }

    render() {
        return (
            <div className="leaderboard">
                <h2 className="leaderboard-title"> Leaderboard </h2>
                <ul>
                    {this.props.stats.map((value, index) => {
                        return <li key={index}>
                            <span className="leaderboard-position">{index + 1}.</span> {value.nickname} ({value.weight})
                        </li>
                    })}
                </ul>
            </div>
        );
    }
}

