import React from "react";
import PropTypes from "prop-types";
import "./modal.css";

export interface Props {
    onClose: Function,
    show: boolean,
    footer?: any,
    width?: number,
    footerClass?: string
}

export default class Modal extends React.Component<Props, {}> {
    onClose = (e: any) => {
        this.props.onClose && this.props.onClose(e);
    };

    render() {
        if (!this.props.show) {
            return null;
        }
        return <div className="modal" id="modal">
            <div className="modal-content" style={{width: this.props.width || 500 + 'px'}}>
                <div className="modal-body">
                    {this.props.children}
                </div>
                <div className={'modal-footer ' + this.props.footerClass || ''}>
                    {this.props.footer}
                </div>
            </div>
        </div>;
    }
}
