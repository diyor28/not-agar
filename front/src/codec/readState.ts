'use strict'

/**
 * Wraps a buffer with a read head pointer
 * @class
 * @param {Buffer} buffer
 */

export default class ReadState {
	offset = 0;
	private readonly buffer: Buffer;

	constructor(buffer: Buffer) {
		this.buffer = buffer;
	}

	peekUInt8() {
		return this.buffer.readUInt8(this.offset);
	}

	readUInt8(): number {
		return this.buffer.readUInt8(this.offset ++);
	}

	readUInt16(): number {
		const r = this.buffer.readUInt16BE(this.offset);
		this.offset += 2;
		return r;
	}

	readUInt32(): number {
		const r = this.buffer.readUInt32BE(this.offset);
		this.offset += 4;
		return r;
	}

	readDouble(): number {
		const r = this.buffer.readDoubleBE(this.offset);
		this.offset += 8;
		return r;
	}

	readBuffer(length: number) {
		if (this.offset + length > this.buffer.length) {
			throw new RangeError('Trying to access beyond buffer length');
		}
		const r = this.buffer.slice(this.offset, this.offset + length);
		this.offset += length;
		return r;
	}

	hasEnded() {
		return this.offset === this.buffer.length;
	}
}
