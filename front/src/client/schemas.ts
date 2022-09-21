import {Schema} from "../codec";

export const genericSchema = new Schema({
	event: 'string'
});

export const pingSchema = genericSchema.extends({
	timestamp: 'uint64'
});

export const moveSchema = genericSchema.extends({
	newX: 'float32',
	newY: 'float32'
});

export const movedSchema = genericSchema.extends({
	x: 'float32',
	y: 'float32',
	weight: 'float32',
	velocityX: 'float32',
	velocityY: 'float32',
	zoom: 'float32',
	points: [{x: 'float32', y: 'float32'}]
});

export const startSchema = genericSchema.extends({
	nickname: 'string'
});

export const foodSchema = genericSchema.extends({
	food: [
		{
			x: 'float32',
			y: 'float32',
			weight: 'float32',
			color: ['uint8']
		}
	]
});

export const playersSchema = genericSchema.extends({
	players: [
		{
			x: 'float32',
			y: 'float32',
			weight: 'float32',
			nickname: 'string',
			color: ['uint8']
		},
	]
});

export const startedSchema = genericSchema.extends({
	player: {
		x: 'float32',
		y: 'float32',
		weight: 'float32',
		color: ['uint8'],
		points: [{x: 'float32', y: 'float32'}]
	},
	spikes: [
		{
			x: 'float32',
			y: 'float32',
			weight: 'float32'
		}
	]
});

export const statsSchema = genericSchema.extends({
	topPlayers: [
		{
			nickname: 'string',
			weight: 'int16'
		}
	]
});

export const ripSchema = genericSchema;
