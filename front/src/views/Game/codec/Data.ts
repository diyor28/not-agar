/**
 * A mutable-length write-only Buffer
 * @class
 * @param {number} [capacity=128] - initial Buffer size
 */

export default class Data {
	private length = 0
	private buffer: Buffer

	constructor(capacity?: number) {
		this.buffer = new Buffer(capacity || 128)
	}

	public prependBuffer(data: Buffer) {
		this.buffer = Buffer.concat([data, this.toBuffer()]);
		this.length += data.length;
	}

	public appendBuffer(data: Buffer) {
		this.buffer = Buffer.concat([this.toBuffer(), data]);
		this.length += data.length;
	}

	public writeUInt8(value: number, prepend: boolean = false) {
		if (prepend) {
			this.prependBuffer(new Buffer(1));
			this.buffer.writeUInt8(value, 0);
		} else {
			this.alloc(1);
			this.buffer.writeUInt8(value, this.length);
		}
		this.length ++;
	}

	public writeUInt16(value: number, prepend: boolean = false) {
		if (prepend) {
			this.prependBuffer(new Buffer(2));
			this.buffer.writeUInt16BE(value, 0);
		} else {
			this.alloc(2);
			this.buffer.writeUInt16BE(value, this.length);
		}
		this.length += 2;
	}

	public writeUInt32(value: number, prepend: boolean = false) {
		if (prepend) {
			this.prependBuffer(new Buffer(4));
			this.buffer.writeUInt16BE(value, 0);
		} else {
			this.alloc(4);
			this.buffer.writeUInt32BE(value, this.length);
		}
		this.length += 4;
	}

	public writeDouble(value: number, prepend: boolean = false) {
		if (prepend) {
			this.prependBuffer(new Buffer(8));
			this.buffer.writeDoubleBE(value, 0);
		} else {
			this.alloc(8);
			this.buffer.writeDoubleBE(value, this.length);
		}
		this.length += 8;
	}

	public toBuffer() {
		return this.buffer.slice(0, this.length);
	}

	private alloc(bytes: number, shift: number = 0) {
		let buffLen = this.buffer.length;


		if (this.length + bytes > buffLen) {
			do {
				buffLen *= 2;
			} while (this.length + bytes > buffLen)

			const newBuffer = new Buffer(buffLen);
			this.buffer.copy(newBuffer, shift, 0, this.length + shift);
			this.buffer = newBuffer;
		}
	}
}
