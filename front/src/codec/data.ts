export const POW = (function () {
	const r = [];
	let i, n = 1
	for (i = 0; i <= 56; i ++) {
		r.push(n)
		n *= 2
	}
	return r
})()

// Pre-calculated constants
export const MAX_DOUBLE_INT = POW[63],
	MAX_UINT8 = POW[8],
	MAX_UINT16 = POW[16],
	MAX_UINT32 = POW[32],
	MAX_INT8 = POW[7],
	MAX_INT16 = POW[15],
	MAX_INT32 = POW[31]


/**
 * A mutable-length write-only Buffer
 * @class
 * @param {number} [capacity=128] - initial Buffer size
 */

export default class Data {
	length = 0
	private buffer: Buffer
	private explanations: { bytes: number, explanation?: string }[] = []

	constructor(capacity?: number) {
		this.buffer = new Buffer(capacity || 128)
	}

	explain(): string {
		let hexString = this.divideFormatString();
		let explanation = '\n';
		const maxExpLength = this.maxExpLength();
		for (let i = 0; i < maxExpLength; i ++) {
			this.explanations.forEach(exp => {
				const bytes = exp.bytes * 2;
				if (exp.explanation && exp.explanation[i]) {
					explanation += exp.explanation[i] + ' '.repeat(bytes - 1) + '|';
				} else {
					explanation += ' '.repeat(bytes) + '|';
				}
			});
			explanation += '\n';
		}
		return hexString + explanation;
	}

	writeUint(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_DOUBLE_INT || value < 0) {
			throw new TypeError('Expected uint, got ' + value);
		}

		if (value < MAX_UINT8) {
			return this.writeUInt8(value, explanation);
		}

		if (value < MAX_UINT16) {
			return this.writeUInt16(value + 0x8000, explanation);
		}

		if (value < MAX_UINT32) {
			return this.writeUInt32(value + 0xc0000000, explanation);
		}

		this.writeUInt64(value, explanation);
	}

	writeUInt8(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_UINT8 || value < 0) {
			throw new TypeError('Expected uint8, got ' + value);
		}
		this.alloc(1);
		this.buffer.writeUInt8(value, this.length);
		this.length ++;
		this.explanations.push({bytes: 1, explanation});
	}

	writeUInt16(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_UINT16 || value < 0) {
			throw new TypeError('Expected uint16, got ' + value);
		}
		this.alloc(2);
		this.buffer.writeUInt16BE(value, this.length);
		this.length += 2;
		this.explanations.push({bytes: 2, explanation});
	}

	writeUInt32(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_UINT32 || value < 0) {
			throw new TypeError('Expected uint32, got ' + value);
		}
		this.alloc(4);
		this.buffer.writeUInt32BE(value, this.length);
		this.length += 4;
		this.explanations.push({bytes: 4, explanation});
	}

	writeUInt64(value: number, explanation?: string) {
		if (value > MAX_DOUBLE_INT || value < 0) {
			throw new TypeError('Expected uint64, got ' + value);
		}
		this.writeUInt32(Math.floor(value / POW[32]) + 0xe0000000, explanation);
		this.writeUInt32(value >>> 0, explanation);
	}

	writeInt(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_DOUBLE_INT || value < - MAX_DOUBLE_INT) {
			throw new TypeError('Expected signed integer at got ' + value);
		}

		if (value >= - MAX_INT8 && value < MAX_INT8) {
			return this.writeInt8(value, explanation);
		}

		if (value >= - MAX_INT16 && value < MAX_INT16) {
			return this.writeInt16(value, explanation);
		}

		if (value >= - MAX_INT32 && value < MAX_INT32) {
			return this.writeInt32(value, explanation);
		}

		this.writeInt64(value, explanation);
	}

	writeInt8(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_INT8 || value < - MAX_INT8) {
			throw new TypeError('Expected int8, got ' + value);
		}
		this.alloc(1);
		this.buffer.writeInt8(value, this.length);
		this.length += 1;
		this.explanations.push({bytes: 1, explanation});
	}

	writeInt16(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_INT16 || value < - MAX_INT16) {
			throw new TypeError('Expected int16, got ' + value);
		}
		this.alloc(2);
		this.buffer.writeInt16BE(value, this.length);
		this.length += 2;
		this.explanations.push({bytes: 2, explanation});
	}

	writeInt32(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_INT32 || value < - MAX_INT32) {
			throw new TypeError('Expected int32, got ' + value);
		}
		this.alloc(4);
		this.buffer.writeInt32BE(value, this.length);
		this.length += 4;
		this.explanations.push({bytes: 4, explanation});
	}

	writeInt64(value: number, explanation?: string) {
		if (Math.round(value) !== value || value > MAX_DOUBLE_INT || value < - MAX_DOUBLE_INT) {
			throw new TypeError('Expected int64, got ' + value);
		}
		this.writeUInt32((Math.floor(value / POW[32]) & 0x1fffffff) + 0xe0000000);
		this.writeUInt32(value >>> 0);
		this.explanations.push({bytes: 8, explanation});
	}

	writeBuffer(b: Buffer, explanation?: string) {
		if (!Buffer.isBuffer(b)) {
			throw new TypeError('Expected a Buffer got ' + b);
		}
		this.writeUInt16(b.length, 'buffer length');
		this.appendBuffer(b, explanation);
	}

	writeString(s: string, explanation?: string) {
		const b = new Buffer(s);
		this.writeUInt16(b.length, 'string length');
		this.appendBuffer(b, explanation);
	}

	writeBoolean(b: boolean, explanation?: string) {
		this.writeUInt8(b ? 1 : 0, explanation);
	}

	writeFloat32(value: number, explanation?: string) {
		this.alloc(4);
		this.buffer.writeFloatBE(value, this.length);
		this.length += 4;
		this.explanations.push({bytes: 4, explanation});
	}

	writeFloat64(value: number, explanation?: string) {
		this.alloc(8);
		this.buffer.writeDoubleBE(value, this.length);
		this.length += 8;
		this.explanations.push({bytes: 8, explanation});
	}

	prependBuffer(data: Buffer, explanation?: string) {
		this.buffer = Buffer.concat([data, this.toBuffer()]);
		this.length += data.length;
		this.explanations.unshift({bytes: data.length, explanation});
	}

	appendBuffer(data: Buffer, explanation?: string) {
		this.buffer = Buffer.concat([this.toBuffer(), data]);
		this.length += data.length;
		this.explanations.push({bytes: data.length, explanation});
	}

	toBuffer() {
		return this.buffer.slice(0, this.length);
	}

	private divideFormatString() {
		let index = 0;
		let hexString = this.toBuffer().toString('hex');
		this.explanations.forEach(el => {
			index += el.bytes * 2;
			hexString = hexString.substr(0, index) + '|' + hexString.substr(index);
			index ++;
		});
		return hexString;
	}

	private maxExpLength(): number {
		return this.explanations.map(el => {
			return (el.explanation || '').length
		}).reduce((el, max) => {
			return el > max ? el : max;
		});
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
