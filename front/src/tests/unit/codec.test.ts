import {Schema} from '../../codec'
import * as assert from 'assert'

describe('Schema.encode', () => {
	test('string', () => {
		const schema = new Schema({
			event: 'string',
		});
		const data = schema.encode({event: 'update'});
		const decoded = schema.decode(data.toBuffer());
		assert.strictEqual(data.toBuffer().toString('hex'), '01010006757064617465');
		assert.strictEqual(decoded.event, 'update');
	});

	test('string, {uint8, uint8, uint8}, [{string, int32}]', () => {
		const schema = new Schema({
			event: 'string',
			player: {x: 'uint8', y: 'uint8', z: 'uint8'},
			stats: [{
				playerId: 'string',
				score: 'int32'
			}]
		});
		const data = schema.encode({
			event: 'update',
			player: {x: 10, y: 12},
			stats: [{playerId: '432142c', score: 32}]
		});
		const decoded = schema.decode(data.toBuffer());
		assert.strictEqual(data.toBuffer().toString('hex'), '0107000675706461746501060a0c0001010300073433323134326300000020');
		assert.strictEqual(decoded.event, 'update');
		assert.strictEqual(decoded.player.x, 10);
		assert.strictEqual(decoded.player.y, 12);
		assert.strictEqual(decoded.player.z, undefined);
		assert.strictEqual(decoded.stats[0].playerId, '432142c');
		assert.strictEqual(decoded.stats[0].score, 32);
	});
})