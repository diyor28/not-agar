export interface EntityData {
	x: number,
	y: number,
	weight: number,
	color: number[]
}

export interface SpikeData extends EntityData {

}

export interface FoodData extends EntityData {
	id: number
}

export interface PlayerData extends EntityData {
	nickname: string
}


export interface SelfPlayerData {
	uuid: string
	x: number
	y: number
	nickname: string
	weight: number
	color: number[]
	zoom: number
	points: {x: number, y: number}[]
}

export interface MovedEvent {
	x: number
	y: number
	zoom: number
	weight: number
	velocityX: number
	velocityY: number
	points: {x: number, y: number}[]
}

export interface MoveCommand {
	newX: number
	newY: number
}

export interface InitialData {
	player: SelfPlayerData
	spikes: SpikeData[]
	food: FoodData[]
}