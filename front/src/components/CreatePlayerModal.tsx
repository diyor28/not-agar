import React, {FormEvent} from 'react';
import Modal from './Modal'
import axios from 'axios'
import '../App.css'

export interface Props {
    onCreated: Function
}

export default class CreatePlayerModal extends React.Component<Props, {}> {
    apiUrl = process.env.REACT_APP_API_URL
    state = {
        show: true,
        nickname: ''
    };
    hideModal = () => {
        this.setState({
            show: false
        });
    };

    createPlayer = (event: FormEvent) => {
        event.preventDefault()
        let data = {nickname: this.state.nickname}
        axios.post(this.apiUrl + "/players", data).then(r => {
            this.props.onCreated(r.data)
        })
        this.hideModal()
    }

    handleChange = (event: any) => {
        this.setState({nickname: event.target.value})
    }

    render() {
        return (
            <div>
                <Modal onClose={this.hideModal}
                       width={300}
                       show={this.state.show} footerClass="align-center"
                       footer={
                           <button className="btn btn-lg btn-success" onClick={this.createPlayer}>
                               Play
                           </button>
                       }>
                    <form onSubmit={this.createPlayer}>
                        <input className={'my-5'} placeholder={'Your nickname'} type="text" value={this.state.nickname} onChange={this.handleChange}/>
                    </form>
                </Modal>
            </div>
        );
    }
}

