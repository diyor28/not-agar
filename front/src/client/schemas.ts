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
	points: {
		type: 'array',
		of: {x: 'float32', y: 'float32'},
		maxLen: 255
	}
});

export const startSchema = genericSchema.extends({
	nickname: {type: 'string', maxLen: 255}
});

export const foodSchema = genericSchema.extends({
	food: {
		type: 'array',
		of: {
			x: 'float32',
			y: 'float32',
			weight: 'float32',
			color: {type: 'array', of: 'uint8', length: 3}
		},
		maxLen: 255
	}
});

export const playersSchema = genericSchema.extends({
	players: {
		type: 'array',
		of: {
			x: 'float32',
			y: 'float32',
			weight: 'float32',
			nickname: {type: 'string', maxLen: 255},
			color: {type: 'array', of: 'uint8', length: 3}
		},
		maxLen: 255
	}
});

export const startedSchema = genericSchema.extends({
	player: {
		x: 'float32',
		y: 'float32',
		weight: 'float32',
		color: {type: 'array', of: 'uint8', length: 3},
		points: {
			type: 'array',
			of: {x: 'float32', y: 'float32'},
			maxLen: 255
		}
	},
	spikes: {
		type: 'array',
		of: {
			x: 'float32',
			y: 'float32',
			weight: 'float32'
		}
	}
});

export const statsSchema = genericSchema.extends({
	topPlayers: {
		type: 'array',
		of: {
			nickname: {type: 'string', maxLen: 255},
			weight: 'int16'
		},
		maxLen: 255
	}
});

export const ripSchema = genericSchema;
