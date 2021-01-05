import React, {FormEvent} from 'react';
import Modal from '../../components/Modal'
import axios from 'axios'
import ButtonLoading from '../../components/ButtonLoading'

export interface Props {
    onCreated: Function,
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
                           <ButtonLoading btnClass="btn btn-lg btn-success"
                                          onClick={this.createPlayer}>
                               Play
                           </ButtonLoading>
                       }>
                    <form className={'my-5'} onSubmit={this.createPlayer}>
                        <input autoFocus={true} placeholder={'Your nickname'} type="text" value={this.state.nickname}
                               onChange={this.handleChange}/>
                    </form>
                    <div style={{marginBottom: "10px"}}>
                        Tip: Press <kbd>a</kbd> to accelerate
                    </div>
                </Modal>
            </div>
        );
    }
}

