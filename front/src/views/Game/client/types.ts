export interface EntityData {
	x: number,
	y: number,
	cameraX: number,
	cameraY: number,
	weight: number,
	color: number[]
}

export interface SpikeData extends EntityData {

}

export interface FoodData extends EntityData {

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
}

export interface MovedEvent {
	x: number
	y: number
	zoom: number
	weight: number
}

export interface MoveCommand {
	newX: number
	newY: number
}

export interface InitialData {
	player: SelfPlayerData
	spikes: SpikeData[]
}