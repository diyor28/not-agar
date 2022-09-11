'use strict'

import {POW} from "./data";

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

	readUInt64(): number {
		return (this.readUInt32() - 0xe0000000) * POW[32] + this.readUInt32()
	}

	readInt8(): number {
		return this.buffer.readInt8(this.offset ++);
	}

	readInt16(): number {
		const r = this.buffer.readInt16BE(this.offset);
		this.offset += 2;
		return r;
	}

	readInt32(): number {
		const r = this.buffer.readInt32BE(this.offset);
		this.offset += 4;
		return r;
	}

	readInt64(): bigint {
		const r = this.buffer.readBigInt64BE(this.offset);
		this.offset += 8;
		return r;
	}

	readString(): string {
		return this.readBuffer().toString();
	}

	readBoolean(): boolean {
		const b = this.readUInt8();
		if (b > 1) {
			throw new Error('Invalid boolean value')
		}
		return Boolean(b);
	}

	readDouble(): number {
		const r = this.buffer.readDoubleBE(this.offset);
		this.offset += 8;
		return r;
	}

	readBuffer(): Buffer {
		const length = this.readUInt16();
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
