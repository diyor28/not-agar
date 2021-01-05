import React from "react";
import {isMobile} from "./engine/utils";

export default class Tips extends React.Component {
    render() {
        if (isMobile())
            return null
        return (
            <div className="tips">
                Tip: Press <kbd>a</kbd> to accelerate
            </div>
        );
    }
}

