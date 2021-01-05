import React from 'react';
import Modal from '../../components/Modal'

export interface Props {
    show: boolean,
}

export default class RIP extends React.Component<Props, {}> {
    apiUrl = process.env.REACT_APP_API_URL
    state = {
        nickname: ''
    };

    hideModal = () => {

    }

    refreshPage = () => {
        window.location.reload()
    }

    render() {
        if (!this.props.show)
            return null
        return (
            <div>
                <Modal onClose={this.hideModal}
                       width={300}
                       show={this.props.show} footerClass="align-center"
                       footer={
                           <button onClick={this.refreshPage} className="btn btn-primary">Play again</button>
                       }>
                    <div style={{minHeight: "100px", display: "flex", alignItems: "center", justifyContent: "center"}}>
                        <h3>What a looser!</h3>
                    </div>
                </Modal>
            </div>
        );
    }
}

