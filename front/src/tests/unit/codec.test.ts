import {Schema} from '../../codec'
import * as assert from 'assert'

describe('Schema.encode', () => {
	test('encode string', () => {
		const schema = new Schema({
			event: 'string',
		});
		const data = schema.encode({event: 'update'});
		const decoded = schema.decode(data.toBuffer());
		assert.strictEqual(data.toBuffer().toString('hex'), '0006757064617465');
		assert.strictEqual(decoded.event, 'update');
	});

	test('encode string.maxLen=255', () => {
		const schema = new Schema({
			event: {type: 'string', maxLen: 255},
		});
		const data = schema.encode({event: 'update'});
		assert.strictEqual(data.toBuffer().toString('hex'), '06757064617465');
	});

	test('encode {string, [{string, int16}]}', () => {
		const schema = new Schema({
			event: 'string',
			topPlayers: {
				type: 'array',
				of: {
					nickname: 'string',
					weight: 'int16'
				}
			}
		});
		const data = schema.encode({
			event: 'stats',
			topPlayers: [
				{
					nickname: 'demo',
					weight: 300
				}
			]
		});
		assert.strictEqual(data.toBuffer().toString('hex'), '000573746174730001000464656d6f012c');
	});

	test('encode {string, {uint8, uint8, uint8}, [{string, int32}]}', () => {
		const schema = new Schema({
			event: 'string',
			player: {x: 'uint8', y: 'uint8', z: {type: 'uint8', optional: true}},
			stats: {
				type: 'array',
				of: {
					playerId: 'string',
					score: 'int32'
				}
			}
		});
		const data = schema.encode({
			event: 'update',
			player: {x: 10, y: 12},
			stats: [{playerId: '432142c', score: 32}]
		});
		assert.strictEqual(data.toBuffer().toString('hex'), '000675706461746501060a0c000100073433323134326300000020');
	});
});


describe('Schema.decode', () => {
	test('decode string.maxLen=255', () => {
		const schema = new Schema({
			event: {type: 'string', maxLen: 255},
		});
		const decoded = schema.decode(Buffer.from('06757064617465', 'hex'));
		assert.strictEqual(decoded.event, 'update');
	})

	test('decode {string, int16}', () => {
		const schema = new Schema({
			event: 'string',
			topPlayers: {
				type: 'array',
				of: {
					nickname: 'string',
					weight: 'int16'
				}
			}
		});
		const decoded = schema.decode(Buffer.from('000573746174730001000464656d6f012c', 'hex'));
		assert.strictEqual(decoded.event, 'stats');
		assert.strictEqual(decoded.topPlayers[0].nickname, 'demo');
		assert.strictEqual(decoded.topPlayers[0].weight, 300);
	});

	test('decode {string, {uint8, uint8, uint8}, [{string, int32}]}', () => {
		const schema = new Schema({
			event: 'string',
			player: {x: 'uint8', y: 'uint8', z: {type: 'uint8', optional: true}},
			stats: {
				type: 'array',
				of: {
					playerId: 'string',
					score: 'int32'
				}
			}
		});
		const decoded = schema.decode(Buffer.from('000675706461746501060a0c000100073433323134326300000020', 'hex'));
		assert.strictEqual(decoded.event, 'update');
		assert.strictEqual(decoded.player.x, 10);
		assert.strictEqual(decoded.player.y, 12);
		assert.strictEqual(decoded.player.z, undefined);
		assert.strictEqual(decoded.stats[0].playerId, '432142c');
		assert.strictEqual(decoded.stats[0].score, 32);
	});

	test('decode {string, {[uint8]}}', () => {
		const schema = new Schema({
			event: 'string',
			player: {color: {type: 'array', of: 'uint8', length: 3}},
		});
		const decoded = schema.decode(Buffer.from('0006757064617465ff7864', 'hex'));
		assert.strictEqual(decoded.event, 'update');
		assert.strictEqual(decoded.player.color[0], 255);
		assert.strictEqual(decoded.player.color[1], 120);
		assert.strictEqual(decoded.player.color[2], 100);
	});

	test('decode {[{float32, float32}].maxLen=255, string}', () => {
		const schema = new Schema({
			points: {
				type: 'array',
				of: {x: 'float32', y: 'float32'},
				maxLen: 255
			},
			nickname: 'string'
		});
		const decoded = schema.decode(Buffer.from('0241f0000042340000c1c80000c1f00000000464656d6f', 'hex'));
		assert.strictEqual(decoded.points[0].x, 30);
		assert.strictEqual(decoded.points[0].y, 45);
		assert.strictEqual(decoded.points[1].x, -25);
		assert.strictEqual(decoded.points[1].y, -30);
		assert.strictEqual(decoded.nickname, 'demo');
	});
});