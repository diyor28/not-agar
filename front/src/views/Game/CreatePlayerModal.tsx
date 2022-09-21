import React, {FormEvent} from 'react';
import {EventBus} from "../../client";
import ButtonLoading from "../../components/ButtonLoading";

export interface Props {
    eventBus: EventBus,
    createPlayer: Function,
}

export default class CreatePlayerModal extends React.Component<Props, {}> {
    constructor(props: Props) {
        super(props);
        this.props.eventBus.on('create', () => {
            this.props.createPlayer({nickname: this.state.nickname});
            this.hideModal();
        });
    }

    state = {
        show: true,
        nickname: ''
    };

    componentDidMount() {
        const input = document.querySelector<HTMLElement>('.nickname-input');
        if (!input)
            return;
        const rect = input.getBoundingClientRect();
        this.props.eventBus.emit('inputmounted', {top: rect.top, left: rect.left});
    }

    hideModal = () => {
        this.setState({
            show: false
        });
    };

    createPlayer = (event: FormEvent) => {
        event.preventDefault();
        this.props.createPlayer({nickname: this.state.nickname});
        this.hideModal();
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
                        <input className="nickname-input input-block" autoFocus={true} placeholder={'Your nickname'}
                               type="text"
                               value={this.state.nickname}
                               onChange={this.handleChange}/>
                    </form>
                    <ButtonLoading btnClass="btn btn-lg btn-success btn-block" onClick={this.createPlayer}>
                        Play
                    </ButtonLoading>
                </div>
            </div>
        );
    }
}

