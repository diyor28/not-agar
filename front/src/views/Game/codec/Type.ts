import types from './types';
import Field from "./Field";
import Data from "./Data";
import ReadState from "./ReadState";

export type BasicType = keyof typeof types;
export type BasicTypeMapping<T> = Record<string, BasicType | BasicType[] | T | T[]>;

export interface RecursiveTypeMapping extends BasicTypeMapping<RecursiveTypeMapping> {

}

export type CombinedType = BasicType | Array<BasicType> | RecursiveTypeMapping;

const TYPE = {
	UINT: 'uint',
	INT: 'int',
	FLOAT: 'float',
	STRING: 'string',
	BUFFER: 'Buffer',
	BOOLEAN: 'boolean',
	JSON: 'json',
	OID: 'oid',
	REGEX: 'regex',
	DATE: 'date',
	ARRAY: '[array]',
	OBJECT: '{object}'
};

export default class Type {
	public TYPE = TYPE;
	public types = types;
	public type: BasicType;
	public subType?: Type;
	public fields: Field[] = [];

	constructor(type: CombinedType) {
		if (typeof type === 'string') {
			if (type in this.TYPE && type !== this.TYPE.ARRAY && type !== this.TYPE.OBJECT) {
				throw new TypeError('Unknown basic type: ' + type)
			}

			this.type = type;
		} else if (Array.isArray(type)) {
			if (type.length !== 1) {
				throw new TypeError('Invalid array type, it must have exactly one element')
			}

			this.type = this.TYPE.ARRAY as BasicType;
			this.subType = new Type(type[0]);
		} else {
			if (!type || typeof type !== 'object') {
				throw new TypeError('Invalid type: ' + type)
			}

			this.type = this.TYPE.OBJECT as BasicType;
			this.fields = Object.keys(type).map((name) => {
				return new Field(name, type[name]);
			})
		}
	}

	encode(value: any) {
		const data = new Data();
		this.write(value, data, '');
		return data.toBuffer();
	}

	decode(buffer: Buffer) {
		return this.read(new ReadState(buffer))
	}

	write(value: any, data: Data, path: string) {
		let field, subpath, subValue, bitmask = 0;

		if (this.type === this.TYPE.ARRAY && this.subType) {
			// Array field
			return this.writeArray(value, data, path, this.subType)
		} else if (this.type !== this.TYPE.OBJECT) {
			// Simple type
			return this.types[this.type].write(value, data, path)
		}

		// Check for object type
		if (!value || typeof value !== 'object') {
			throw new TypeError('Expected an object at ' + path)
		}

		// Write each field
		for (let i = 0, len = this.fields.length; i < len; i ++) {
			field = this.fields[i];
			subpath = path ? path + '.' + field.name : field.name;
			subValue = value[field.name];

			if (subValue !== undefined && subValue !== null) {
				bitmask |= 1;
				bitmask <<= 1;
			} else {
				bitmask <<= 1;
				continue;
			}

			if (!field.array) {
				// Scalar field
				field.type.write(subValue, data, subpath)
				continue
			}

			// Array field
			this.writeArray(subValue, data, subpath, field.type);
		}
		const bitmaskData = new Data();
		const bitmaskBytes = types.uint.write(bitmask, bitmaskData, 'bitmask');
		bitmaskData.writeUInt8(bitmaskBytes, true);
		data.prependBuffer(bitmaskData.toBuffer());
	}

	read(state: ReadState): any {
		if (this.type !== this.TYPE.OBJECT && this.type !== this.TYPE.ARRAY) {
			// Scalar type
			// In this case, there is no need to write custom code
			return types[this.type].read(state);
		} else if (this.type === this.TYPE.ARRAY) {
			// @ts-ignore
			return this.readArray.bind(this, this.subType)(state);
		}

		const bitmaskBytes = state.readUInt8();
		let bitmask: number;
		switch (bitmaskBytes) {
			case 1:
				bitmask = state.readUInt8();
				break;
			case 2:
				bitmask = state.readUInt16();
				break;
			case 4:
				bitmask = state.readUInt32();
				break;
			case 8:
				bitmask = state.readDouble();
				break;
		}

		const result: Record<string, any> = {};
		this.fields.forEach((field: any, i: number) => {
			let name = JSON.stringify(field.name), value = undefined;
			if (field.array) {
				value = this.readArray(this.fields[i].type, state);
			} else {
				if (this.isMaskTrue(bitmask, i)) {
					value = this.fields[i].type.read(state);
				}
			}
			result[name] = value;
		});

		return result;
	}

	getHash() {
		const hashType = (type: Type, array: boolean, optional: boolean) => {
			// Write type (first char + flags)
			// AOxx xxxx
			hash.writeUInt8((type.type.charCodeAt(0) & 0x3f) | (array ? 0x80 : 0) | (optional ? 0x40 : 0))

			if (type.type === this.TYPE.ARRAY) {
				hashType(type.subType as Type, false, false)
			} else if (type.type === this.TYPE.OBJECT) {
				types.uint.write(type.fields.length, hash)
				type.fields.forEach(function (field: any) {
					hashType(field.type, field.array, field.optional)
				})
			}
		}
		let hash = new Data();
		hashType(this, false, false)
		return hash.toBuffer()
	}

	private isMaskTrue(mask: number, idx: number) {
		return (mask >> (this.fields.length - idx)) & 1; // Shift right to put the target at position 0, and AND it with 1
	}

	private writeArray(value: any, data: Data, path: string, type: Type) {
		let i, len;
		if (!Array.isArray(value)) {
			throw new TypeError('Expected an Array at ' + path)
		}
		len = value.length
		types.uint.write(len, data)
		for (i = 0; i < len; i ++) {
			type.write(value[i], data, path + '.' + i)
		}
	}

	private readArray(type: Type, state: ReadState) {
		let arr = new Array(types.uint.read(state));
		for (let j = 0; j < arr.length; j ++) {
			arr[j] = type.read(state)
		}
		return arr
	}
}

