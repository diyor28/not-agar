import React from 'react';
import '../App.css'

export interface Props {
    onClick: (event: any) => void,
    btnClass: string,
    loading: boolean
}


export default class ButtonLoading extends React.Component<Props, {}> {
    render() {
        const buttonContent = () => {
            if (this.props.loading)
                return (
                    <div>
                        <span className="spinner-grow spinner-grow-sm" role="status" aria-hidden={true}/>
                        <span>Loading...</span>
                    </div>
                )
            return this.props.children
        }
        return (
            // <span className="spinner"/>
            <button className={this.props.btnClass}
                    disabled={this.props.loading}
                    onClick={this.props.onClick}>
                {buttonContent()}
            </button>
        );
    }
}

