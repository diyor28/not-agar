import React from 'react';
import "./components/modal.css";
import CSS from 'csstype';

export interface Props {
    stats: { nickname: string, weight: number }[]
}

const modalStyle: CSS.Properties = {
    position: "absolute",
    top: 0,
    right: 0,
    background: "rgba(0, 0, 0, 0.2)",
    minHeight: '200px',
    width: '200px'
}

const ulStyle: CSS.Properties = {}

const liStyle: CSS.Properties = {
    color: "white"
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
            <div style={modalStyle}>
                <h3 style={{color: "white"}}> Leaderboard </h3>
                {this.props.stats.map((value, index) => {
                    return <li style={liStyle} key={index}>{index+1}. {value.nickname} ({value.weight})</li>
                })}
            </div>
        );
    }
}

