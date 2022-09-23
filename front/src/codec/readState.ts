import {MAX_UINT16, MAX_UINT32, MAX_UINT8, POW} from "./data";

/**
 * Wraps a buffer with a read head pointer
 * @class
 * @param {Buffer} buffer
 */

export function minBytes(n: number) {
	if (n < MAX_UINT8) {
		return 1;
	}

	if (n < MAX_UINT16) {
		return 2;
	}

	if (n < MAX_UINT32) {
		return 4;
	}
	return 8;
}

export default class ReadState {
	offset = 0;
	readonly buffer: Buffer;

	constructor(buffer: Buffer) {
		this.buffer = buffer;
	}

	peekUInt8() {
		return this.buffer.readUInt8(this.offset);
	}

	readUint(size: number): number {
		switch (size) {
			case 1:
				return this.readUInt8()
			case 2:
				return this.readUInt16()
			case 4:
				return this.readUInt32()
			case 8:
				return this.readUInt64()
		}
		throw new TypeError(`Expected size in [1, 2, 4, 8], got ${size}`);
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

	readString(len?: number, maxLen?: number): string {
		return this.readBuffer(len, maxLen).toString();
	}

	readBoolean(): boolean {
		const b = this.readUInt8();
		if (b > 1) {
			throw new Error('Invalid boolean value')
		}
		return Boolean(b);
	}

	readFloat32(): number {
		const r = this.buffer.readFloatBE(this.offset);
		this.offset += 4;
		return r;
	}

	readFloat64(): number {
		const r = this.buffer.readDoubleBE(this.offset);
		this.offset += 8;
		return r;
	}

	readBuffer(len?: number, maxLen?: number): Buffer {
		let length: number;
		if (len) {
			length = len;
		} else if (maxLen) {
			length = this.readUint(minBytes(maxLen));
		} else {
			length = this.readUInt16();
		}
		if (this.offset + length > this.buffer.length) {
			console.log(length, len, maxLen, this.buffer.toString('hex'));
			throw new RangeError('Trying to access beyond buffer length');
		}
		const r = this.buffer.slice(this.offset, this.offset + length);
		this.offset += length;
		return r;
	}
}
