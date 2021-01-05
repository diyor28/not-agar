import React, {FormEvent} from 'react';
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
        if (!this.state.show) {
            return null;
        }
        return (
            <div className="modal modal-transparent" id="modal">
                <div className="modal-content" style={{width: 300 + 'px'}}>
                    <form className={'my-5'} onSubmit={this.createPlayer}>
                        <input className="input-block" autoFocus={true} placeholder={'Your nickname'} type="text"
                               value={this.state.nickname}
                               onChange={this.handleChange}/>
                    </form>
                    <ButtonLoading btnClass="btn btn-lg btn-success btn-block"
                                   onClick={this.createPlayer}>
                        Play
                    </ButtonLoading>
                </div>
            </div>
        );
    }
}

