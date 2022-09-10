const POW = (function () {
	const r = [];
	let i, n = 1
	for (i = 0; i <= 56; i ++) {
		r.push(n)
		n *= 2
	}
	return r
})()

// Pre-calculated constants
const MAX_DOUBLE_INT = POW[53],
	MAX_UINT8 = POW[7],
	MAX_UINT16 = POW[14],
	MAX_UINT32 = POW[29],
	MAX_INT8 = POW[6],
	MAX_INT16 = POW[13],
	MAX_INT32 = POW[28]


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

	public writeUint(value: number, prepend: boolean = false) {
		if (Math.round(value) !== value || value > MAX_DOUBLE_INT || value < 0) {
			throw new TypeError('Expected unsigned integer got ' + value);
		}

		if (value < MAX_UINT8) {
			return this.writeUInt8(value, prepend);
		}

		if (value < MAX_UINT16) {
			return this.writeUInt16(value + 0x8000, prepend);
		}

		if (value < MAX_UINT32) {
			return this.writeUInt32(value + 0xc0000000, prepend);
		}

		this.writeUInt64(value, prepend);
	}


	public writeUInt8(value: number, prepend: boolean = false) {
		if (Math.round(value) !== value || value > MAX_UINT8 || value < 0) {
			throw new TypeError('Expected unsigned integer got ' + value);
		}
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
		if (Math.round(value) !== value || value > MAX_UINT16 || value < 0) {
			throw new TypeError('Expected unsigned integer got ' + value);
		}
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
		if (Math.round(value) !== value || value > MAX_UINT32 || value < 0) {
			throw new TypeError('Expected unsigned integer got ' + value);
		}
		if (prepend) {
			this.prependBuffer(new Buffer(4));
			this.buffer.writeUInt32BE(value, 0);
		} else {
			this.alloc(4);
			this.buffer.writeUInt32BE(value, this.length);
		}
		this.length += 4;
	}

	public writeUInt64(value: number, prepend: boolean = false) {
		if (Math.round(value) !== value || value > MAX_DOUBLE_INT || value < 0) {
			throw new TypeError('Expected unsigned integer got ' + value);
		}
		if (prepend) {
			this.writeUInt32(value >>> 0, prepend);
			this.writeUInt32(Math.floor(value / POW[32]) + 0xe0000000, prepend);
		} else {
			this.writeUInt32(Math.floor(value / POW[32]) + 0xe0000000);
			this.writeUInt32(value >>> 0);
		}
		this.length += 8;
	}

	public writeInt(value: number) {
		if (Math.round(value) !== value || value > MAX_DOUBLE_INT || value < - MAX_DOUBLE_INT) {
			throw new TypeError('Expected signed integer at got ' + value);
		}

		if (value >= - MAX_INT8 && value < MAX_INT8) {
			return this.writeInt8(value);
		}

		if (value >= - MAX_INT16 && value < MAX_INT16) {
			return this.writeInt16(value);
		}

		if (value >= - MAX_INT32 && value < MAX_INT32) {
			return this.writeInt32(value);
		}

		this.writeInt64(value);
	}

	public writeInt8(value: number) {
		if (Math.round(value) !== value || value > MAX_INT8 || value < - MAX_INT8) {
			throw new TypeError('Expected signed integer got ' + value);
		}
		this.writeUInt8(value & 0x7f);
	}

	public writeInt16(value: number) {
		if (Math.round(value) !== value || value > MAX_INT16 || value < - MAX_INT16) {
			throw new TypeError('Expected signed integer got ' + value);
		}
		this.writeUInt16((value & 0x3fff) + 0x8000);
	}

	public writeInt32(value: number) {
		if (Math.round(value) !== value || value > MAX_INT32 || value < - MAX_INT32) {
			throw new TypeError('Expected signed integer got ' + value);
		}
		this.writeUInt32((value & 0x1fffffff) + 0xc0000000);
	}

	public writeInt64(value: number) {
		if (Math.round(value) !== value || value > MAX_DOUBLE_INT || value < - MAX_DOUBLE_INT) {
			throw new TypeError('Expected signed integer got ' + value);
		}
		this.writeUInt32((Math.floor(value / POW[32]) & 0x1fffffff) + 0xe0000000);
		this.writeUInt32(value >>> 0);
	}

	public writeBuffer(b: Buffer) {
		if (!Buffer.isBuffer(b)) {
			throw new TypeError('Expected a Buffer got ' + b);
		}
		this.writeUInt32(b.length);
		this.appendBuffer(b);
	}

	public writeString(s: string) {
		this.writeBuffer(new Buffer(s));
	}

	public writeBoolean(b: boolean) {
		this.writeUInt8(b ? 1 : 0);
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

	public prependBuffer(data: Buffer) {
		this.buffer = Buffer.concat([data, this.toBuffer()]);
		this.length += data.length;
	}

	public appendBuffer(data: Buffer) {
		this.buffer = Buffer.concat([this.toBuffer(), data]);
		this.length += data.length;
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
