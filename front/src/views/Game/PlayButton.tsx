import React from "react";
import ButtonLoading from "../../components/ButtonLoading";
import {EventBus} from "../../client";
import {calcDistance} from "../../engine/utils";


export interface Props {
	eventBus: EventBus,
}

const getRandomNumber = (num: number) => Math.floor(Math.random() * (num + 1));

export default class PlayButton extends React.Component<Props, {}> {
	state = {
		top: 0,
		left: 0
	}
	mountPromise: Promise<void>
	resolveMountPromise = (value: any) => {
	}

	constructor(props: Props) {
		super(props);
		this.mountPromise = new Promise(resolve => {
			this.resolveMountPromise = resolve;
		});
		this.props.eventBus.on('inputmounted', ({top, left}) => {
			this.mountPromise.then(() => {
				this.setState({top: top + 50, left});
			});
		})

		window.addEventListener('mousemove', (event: MouseEvent) => {
			const btn = document.querySelector<HTMLElement>('.play-btn');
			if (!btn)
				return;
			if (btn.getAnimations().length) {
				return;
			}
			const x = Math.max(btn.offsetLeft, Math.min(event.clientX, btn.offsetLeft + btn.offsetWidth));
			const y = Math.max(btn.offsetTop, Math.min(event.clientY, btn.offsetTop + btn.offsetHeight));
			const dist = calcDistance(event.clientX, event.clientY, x, y);
			if (dist > 10) {
				return;
			}
			const top = getRandomNumber(window.innerHeight - btn.offsetHeight);
			const left = getRandomNumber(window.innerWidth - btn.offsetWidth);
			const animation = btn.animate({
				top: `${top}px`,
				left: `${left}px`,
				easing: 'ease-out'
			}, 1000);
			animation.finished.then(() => {
				btn.style.top = `${top}px`;
				btn.style.left = `${left}px`;
			})
		});
	}

	componentDidMount() {
		this.resolveMountPromise(null);
	}


	render() {
		return <div className="play-btn"
					style={{
						position: 'absolute',
						width: '300px',
						top: `${this.state.top}px`,
						left: `${this.state.left}px`
					}}>
			<ButtonLoading btnClass="btn btn-lg btn-success btn-block"
						   onClick={() => this.props.eventBus.emit('create', {})}>
				Play
			</ButtonLoading>
		</div>
	}
}