import Data from "./data";
import ReadState from "./readState";

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


export interface TypeIO {
	write(u: any, data: Data, path?: string): number

	read(state: ReadState): any
}

/*
 * Formats (big-endian):
 * 7b	0xxx xxxx
 * 14b	10xx xxxx  xxxx xxxx
 * 29b	110x xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx
 * 61b	111x xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx  xxxx xxxx
 */
export const uint8: TypeIO = {
	write: function (u: number, data: Data, path?: string) {
		// Check the input
		if (Math.round(u) !== u || u > MAX_UINT8 || u < 0) {
			throw new TypeError('Expected unsigned integer at ' + path + ', got ' + u)
		}
		data.writeUInt8(u);
		return 1;
	},
	read: (state: ReadState) => {
		return state.readUInt8();
	}
}

export const uint16: TypeIO = {
	read: (state: ReadState) => {
		return state.readUInt16();
	},
	write: (u: number, data: Data, path?: string) => {
		// Check the input
		if (Math.round(u) !== u || u > MAX_UINT16 || u < 0) {
			throw new TypeError('Expected unsigned integer at ' + path + ', got ' + u);
		}
		data.writeUInt16(u + 0x8000);
		return 2;
	}
}

export const uint32: TypeIO = {
	read: (state: ReadState) => {
		return state.readUInt32() - 0xc0000000;
	},
	write: (u: number, data: Data, path?: string) => {
		if (Math.round(u) !== u || u > MAX_UINT32 || u < 0) {
			throw new TypeError('Expected unsigned integer at ' + path + ', got ' + u);
		}
		data.writeUInt32(u + 0xc0000000);
		return 4;
	}
}


export const uint64: TypeIO = {
	read: (state: ReadState) => {
		return (state.readUInt32() - 0xe0000000) * POW[32] + state.readUInt32()
	},
	write: (u: number, data: Data, path?: string) => {
		if (Math.round(u) !== u || u > MAX_DOUBLE_INT || u < 0) {
			throw new TypeError('Expected unsigned integer at ' + path + ', got ' + u);
		}
		// Split in two 32b uints
		data.writeUInt32(Math.floor(u / POW[32]) + 0xe0000000);
		data.writeUInt32(u >>> 0);
		return 8;
	}
}

export const uintT: TypeIO = {
	write: function (u: number, data: Data, path?: string): number {
		// Check the input
		if (Math.round(u) !== u || u > MAX_DOUBLE_INT || u < 0) {
			throw new TypeError('Expected unsigned integer at ' + path + ', got ' + u)
		}

		if (u < MAX_UINT8) {
			data.writeUInt8(u);
			return 1;
		}

		if (u < MAX_UINT16) {
			data.writeUInt16(u + 0x8000);
			return 2;
		}

		if (u < MAX_UINT32) {
			data.writeUInt32(u + 0xc0000000)
			return 4;
		}

		// Split in two 32b uints
		data.writeUInt32(Math.floor(u / POW[32]) + 0xe0000000)
		data.writeUInt32(u >>> 0)
		return 8;
	},
	read: function (state: ReadState) {
		const firstByte = state.peekUInt8();

		if (!(firstByte & 0x80)) {
			state.offset ++
			return firstByte
		}

		if (!(firstByte & 0x40)) {
			return state.readUInt16() - 0x8000
		}

		if (!(firstByte & 0x20)) {
			return state.readUInt32() - 0xc0000000
		}

		return (state.readUInt32() - 0xe0000000) * POW[32] + state.readUInt32()
	}
}

export const int8: TypeIO = {
	read: (state: ReadState) => {
		const firstByte = state.readUInt8();
		return (firstByte & 0x40) ? (firstByte | 0xffffff80) : firstByte;
	},
	write: (i: number, data: Data, path?: string) => {
		if (Math.round(i) !== i || i > MAX_INT8 || i < - MAX_INT8) {
			throw new TypeError('Expected signed integer at ' + path + ', got ' + i);
		}
		data.writeUInt8(i & 0x7f);
		return 1;
	}
}

export const int16: TypeIO = {
	read: (state: ReadState) => {
		let i = state.readUInt16() - 0x8000;
		return (i & 0x2000) ? (i | 0xffffc000) : i;
	},
	write: (i: number, data: Data, path?: string) => {
		if (Math.round(i) !== i || i > MAX_INT16 || i < - MAX_INT16) {
			throw new TypeError('Expected signed integer at ' + path + ', got ' + i);
		}
		data.writeUInt16((i & 0x3fff) + 0x8000);
		return 2;
	}
}

export const int32: TypeIO = {
	read: (state: ReadState) => {
		let i = state.readUInt32() - 0xc0000000;
		return (i & 0x10000000) ? (i | 0xe0000000) : i;
	},
	write: (i: number, data: Data, path?: string) => {
		if (Math.round(i) !== i || i > MAX_INT32 || i < - MAX_INT32) {
			throw new TypeError('Expected signed integer at ' + path + ', got ' + i);
		}
		data.writeUInt32((i & 0x1fffffff) + 0xc0000000);
		return 4;
	}
}

export const int64: TypeIO = {
	read: (state: ReadState) => {
		let i = state.readUInt32() - 0xe0000000;
		i = (i & 0x10000000) ? (i | 0xe0000000) : i;
		return i * POW[32] + state.readUInt32();
	},
	write: (i: number, data: Data, path?: string) => {
		if (Math.round(i) !== i || i > MAX_DOUBLE_INT || i < - MAX_DOUBLE_INT) {
			throw new TypeError('Expected signed integer at ' + path + ', got ' + i);
		}
		// Split in two 32b uints
		data.writeUInt32((Math.floor(i / POW[32]) & 0x1fffffff) + 0xe0000000);
		data.writeUInt32(i >>> 0);
		return 8;
	}
}

export const intT: TypeIO = {
	write: function (i: number, data: Data, path?: string): number {
		// Check the input
		if (Math.round(i) !== i || i > MAX_DOUBLE_INT || i < - MAX_DOUBLE_INT) {
			throw new TypeError('Expected signed integer at ' + path + ', got ' + i)
		}

		if (i >= - MAX_INT8 && i < MAX_INT8) {
			data.writeUInt8(i & 0x7f);
			return 1;
		}

		if (i >= - MAX_INT16 && i < MAX_INT16) {
			data.writeUInt16((i & 0x3fff) + 0x8000);
			return 2;
		}

		if (i >= - MAX_INT32 && i < MAX_INT32) {
			data.writeUInt32((i & 0x1fffffff) + 0xc0000000);
			return 4;
		}

		// Split in two 32b uints
		data.writeUInt32((Math.floor(i / POW[32]) & 0x1fffffff) + 0xe0000000)
		data.writeUInt32(i >>> 0)
		return 8;
	},
	read: function (state: ReadState) {
		let firstByte = state.peekUInt8(), i;

		if (!(firstByte & 0x80)) {
			state.offset ++;
			return (firstByte & 0x40) ? (firstByte | 0xffffff80) : firstByte;
		}

		if (!(firstByte & 0x40)) {
			i = state.readUInt16() - 0x8000;
			return (i & 0x2000) ? (i | 0xffffc000) : i;
		}

		if (!(firstByte & 0x20)) {
			i = state.readUInt32() - 0xc0000000;
			return (i & 0x10000000) ? (i | 0xe0000000) : i;
		}

		i = state.readUInt32() - 0xe0000000;
		i = (i & 0x10000000) ? (i | 0xe0000000) : i;
		return i * POW[32] + state.readUInt32();
	}
}

export const float16: TypeIO = {
	write: function (f: any, data: Data, path: string) {
		if (typeof f !== 'number') {
			throw new TypeError('Expected a number at ' + path + ', got ' + f);
		}
		data.writeDouble(f);
		return 8;
	},
	read: function (state: ReadState) {
		return state.readDouble();
	}
}

export const float32: TypeIO = {
	write: function (f: any, data: Data, path: string) {
		if (typeof f !== 'number') {
			throw new TypeError('Expected a number at ' + path + ', got ' + f);
		}
		data.writeDouble(f);
		return 8;
	},
	read: function (state: ReadState) {
		return state.readDouble();
	}
}

export const float64: TypeIO = {
	write: function (f: any, data: Data, path: string) {
		if (typeof f !== 'number') {
			throw new TypeError('Expected a number at ' + path + ', got ' + f);
		}
		data.writeDouble(f);
		return 8;
	},
	read: function (state: ReadState) {
		return state.readDouble()
	}
}

/*
 * <uint_length> <buffer_data>
 */
export const stringT: TypeIO = {
	write: function (s: any, data: Data, path: string) {
		if (typeof s !== 'string') {
			throw new TypeError('Expected a string at ' + path + ', got ' + s);
		}
		const buffer = new Buffer(s);
		bufferT.write(buffer, data, path);
		return buffer.length;
	},
	read: function (state: ReadState) {
		return bufferT.read(state).toString();
	}
}

/*
 * <uint_length> <buffer_data>
 */
export const bufferT: TypeIO = {
	write: function (B: any, data: Data, path: string) {
		if (!Buffer.isBuffer(B)) {
			throw new TypeError('Expected a Buffer at ' + path + ', got ' + B);
		}
		uintT.write(B.length, data, path);
		data.appendBuffer(B);
		return B.length;
	},
	read: function (state: ReadState) {
		const length = uintT.read(state)
		return state.readBuffer(length)
	}
}

/*
 * either 0x00 or 0x01
 */
export const booleanT: TypeIO = {
	write: function (b: boolean, data: Data) {
		data.writeUInt8(b ? 1 : 0);
		return 1;
	},
	read: function (state: ReadState) {
		const b = state.readUInt8();
		if (b > 1) {
			throw new Error('Invalid boolean value')
		}
		return Boolean(b)
	}
}

/*
 * <uint_length> <buffer_data>
 */
export const jsonT: TypeIO = {
	write: function (j: any, data: Data, path: string) {
		return stringT.write(JSON.stringify(j), data, path);
	},
	read: function (state: ReadState) {
		return JSON.parse(stringT.read(state))
	}
}

/*
 * <12B_buffer_data>
 */
export const oidT = {
	write: function (o: any, data: Data, path: string) {
		const buffer = new Buffer(String(o), 'hex');
		if (buffer.length !== 12) {
			throw new TypeError('Expected an object id (12 bytes) at ' + path + ', got ' + o);
		}
		data.appendBuffer(buffer);
		return buffer.length;
	},
	read: function (state: ReadState) {
		return state.readBuffer(12).toString('hex');
	}
}

/*
 * <uint_source_length> <buffer_source_data> <flags>
 * flags is a bit-mask: g=1, i=2, m=4
 */
export const regexT = {
	write: function (r: any, data: Data, path: string) {
		let g, i, m;
		if (!(r instanceof RegExp)) {
			throw new TypeError('Expected an instance of RegExp at ' + path + ', got ' + r)
		}
		const sLen = stringT.write(r.source, data, path);
		g = r.global ? 1 : 0;
		i = r.ignoreCase ? 2 : 0;
		m = r.multiline ? 4 : 0;
		data.writeUInt8(g + i + m);
		return 1 + sLen;
	},
	read: function (state: ReadState) {
		const source = stringT.read(state),
			flags = state.readUInt8(),
			g = flags & 0x1 ? 'g' : '',
			i = flags & 0x2 ? 'i' : '',
			m = flags & 0x4 ? 'm' : '';
		return new RegExp(source, g + i + m);
	}
}

/*
 * <uint_time_ms>
 */
export const dateT: TypeIO = {
	write: function (d: any, data: Data, path: string) {
		if (!(d instanceof Date)) {
			throw new TypeError('Expected an instance of Date at ' + path + ', got ' + d)
		}
		if (isNaN(d.getTime())) {
			throw new TypeError('Expected a valid Date at ' + path + ', got ' + d)
		}
		return uintT.write(d.getTime(), data, path);
	},
	read: function (state: ReadState) {
		return new Date(uintT.read(state))
	}
}


export type UIntT = 'uint8' | 'uint16' | 'uint32' | 'uint64';
export type IntT = 'int8' | 'int16' | 'int32' | 'int64';
export type FloatT = 'float32' | 'float64';
export type StringT = 'string';
export type BooleanT = 'boolean';
export type DateT = 'date';
export type BufferT = 'buffer';
export type OidT = 'oid';
export type RegexT = 'regex';
export type JsonT = 'json';
export type PrimitiveType = UIntT | IntT | DateT | FloatT | BooleanT | StringT | BufferT | OidT | RegexT | JsonT;
export type ExtendedPrimitiveType = PrimitiveType | 'array' | 'object';


export type BasicTypeMapping<T> = Record<string, PrimitiveType | PrimitiveType[] | T | T[]>;

export interface RecursiveTypeMapping extends BasicTypeMapping<RecursiveTypeMapping> {
}
export type FieldType = PrimitiveType | PrimitiveType[] | RecursiveTypeMapping | Array<RecursiveTypeMapping>


const types: Record<PrimitiveType | 'int' | 'uint' | 'float', TypeIO> = {
	uint8,
	uint16,
	uint32,
	uint64,
	int8,
	int16,
	int32,
	int64,
	float32,
	float64,
	uint: uintT,
	float: uintT,
	int: intT,
	date: dateT,
	string: stringT,
	json: jsonT,
	boolean: booleanT,
	buffer: bufferT,
	regex: regexT,
	oid: oidT
};

export default types;