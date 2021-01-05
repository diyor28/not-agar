import React from 'react';
import "../../components/modal.css";
import '../../App.css'

export interface Props {
    ping: number | null
}

function getColor(ping: number) {
    if (ping < 100)
        return 'dot-success'
    if (ping > 100 && ping < 200)
        return 'dot-warning'
    return 'dot-danger'
}

export default class Ping extends React.Component<Props, {}> {
    state = {};

    render() {
        if (!this.props.ping)
            return null
        return (
            <div className="ping-div">
                <span className={"dot dot-sm " + getColor(this.props.ping)}/>
                <span className="ping-text">{this.props.ping} ms</span>
            </div>
        );
    }
}

