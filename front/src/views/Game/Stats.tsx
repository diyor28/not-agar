import React from 'react';
import "../../components/modal.css";
import {StatsUpdate} from './engine/GameEngine'
import {isMobile} from "./engine/utils";

export interface Props {
    stats: StatsUpdate[]
}

function sliceStats(stats: StatsUpdate[], isMobile: boolean) {
    if (isMobile)
        return stats.slice(0, 5)
    return stats
}

export default class Stats extends React.Component<Props, {}> {
    state = {
        stats: []
    };

    handleChange = (event: any) => {
        this.setState({nickname: event.target.value})
    }

    render() {
        if (!this.props.stats.length)
            return null
        return (
            <div className={'leaderboard ' + (isMobile() ? 'leaderboard-mobile' : '')}>
                <h2 className="leaderboard-title"> Leaderboard </h2>
                <ul>
                    {sliceStats(this.props.stats, isMobile()).map((value, index) => {
                        return <li key={index}>
                            <div className="nickname">{index + 1}. {value.nickname}</div>
                            <div className="weight">({value.weight})</div>
                        </li>
                    })}
                </ul>
            </div>
        );
    }
}

