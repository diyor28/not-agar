import {Schema} from '../../codec'
import * as assert from 'assert'

describe('Type.decode', () => {
	const schema = new Schema({
		event: 'string',
		player: {x: 'uint8', y: 'uint8', z: 'uint8'},
		stats: [{
			playerId: 'string',
			score: 'int32'
		}]
	});
	const buffer = schema.encode({event: 'update', player: {x: 10, y: 12}, stats: [{playerId: '432142c', score: 32}]});
	assert.strictEqual(buffer.toString('hex'), '');
})