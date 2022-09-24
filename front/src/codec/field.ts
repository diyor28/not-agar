import {ExtendedPrimitiveType, isVarSizeTypeConf, StrictFieldType, StrictSchemaType, StrictTypeConf} from "./types";
import Data from "./data";
import ReadState, {minBytes} from "./readState";

function isStrictTypeConf(field: any): field is StrictTypeConf<any> {
	return typeof field === 'object' && field.type
}

export class FieldsMap {
	loc: string;
	fields: Field[];

	constructor(loc: string, type: StrictSchemaType) {
		this.loc = loc;
		this.fields = Object.keys(type).map(key => {
			const subType = type[key];
			return new Field(key, subType, loc ? loc + '.' + key : key);
		});
	}

	private calcBitmask(value: any) {
		let bitmask = 0;
		let bitmaskBits = 0;
		for (const field of this.fields) {
			const subValue = value[field.name];
			if (subValue !== undefined && subValue !== null) {
				bitmask |= 1;
			}
			bitmask <<= 1;
			bitmaskBits ++;
		}
		bitmask >>= 1;
		return {bitmask, bitmaskSize: Math.ceil(bitmaskBits / 8)};
	}

	hasOptionalFields(): boolean {
		return this.fields.some(el => el.optional);
	}

	write(value: any, data: Data) {
		if (this.hasOptionalFields()) {
			const {bitmask, bitmaskSize} = this.calcBitmask(value);
			data.writeUInt8(bitmaskSize, 'bitmask size');
			data.writeUint(bitmask, bitmaskSize, 'bitmask');
		}
		for (const field of this.fields) {
			const subValue = value[field.name];
			if (subValue === undefined || subValue === null) {
				if (field.optional) {
					continue;
				}
				throw new TypeError(`Field '${this.loc}.${field.name}' is not optional, got ${subValue}`);
			}
			field.encode(value[field.name], data);
		}
	}

	read(state: ReadState) {
		let bitmask: number = 0;
		const hasOptionalFields = this.hasOptionalFields()
		if (hasOptionalFields) {
			bitmask = this.readBitmask(state);
		}
		const result: Record<string, any> = {};
		this.fields.forEach((field, i) => {
			if (hasOptionalFields && this.isMaskTrue(bitmask, i) || !hasOptionalFields) {
				result[field.name] = field.decode(state);
			} else {
				result[field.name] = undefined;
			}
		});

		return result;
	}

	private readBitmask(state: ReadState): number {
		let bitmask: number = 0;
		const bitmaskBytes = state.readUInt8();
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
				bitmask = state.readUInt64();
				break;
		}
		return bitmask;
	}

	private isMaskTrue(mask: number, idx: number) {
		return (mask >> (this.fields.length - idx - 1)) & 1; // Shift right to put the target at position 0, and AND it with 1
	}
}

export default class Field {
	name: string;
	loc: string;
	optional = false;
	len = 0;
	maxLen = 0;
	type: ExtendedPrimitiveType
	subType: Field | null = null;
	subFields: FieldsMap | null = null;

	constructor(name: string, field: StrictFieldType, loc: string) {
		this.name = name;
		this.loc = loc;
		if (isVarSizeTypeConf(field)) {
			this.len = field.length;
			this.maxLen = field.maxLen;
		}

		if (!isStrictTypeConf(field)) {
			this.type = 'object';
			this.subFields = new FieldsMap(loc, field);
		} else if (field.type === 'array') {
			this.type = 'array';
			this.subType = new Field(name, field.of, `${loc}[]`);
			this.optional = field.optional;
			this.len = field.length;
			this.maxLen = field.maxLen;
		} else if (field.type === 'object') {
			this.type = 'object';
			this.subFields = new FieldsMap(loc, field.of);
			this.optional = field.optional;
		} else {
			this.type = field.type;
			this.optional = field.optional;
		}
	}

	decode(state: ReadState) {
		switch (this.type) {
			case "array":
				return this.readArray(state);
			case "object":
				if (!this.subFields)
					throw new Error('this.subFields is not defined');
				return this.subFields.read(state);
			default:
				try {
					return this.readPrimitive(state);
				} catch (e){
					e.message = `When trying to decode ${this.loc}: ` + e.message;
					throw e;
				}

		}
	}

	private readPrimitive(state: ReadState) {
		switch (this.type) {
			case "string":
				return state.readString(this.len, this.maxLen);
			case "buffer":
				return state.readBuffer(this.len, this.maxLen);
			case "boolean":
				return state.readBoolean();
			case "float32":
				return state.readFloat32();
			case "float64":
				return state.readFloat64();
			case "int8":
				return state.readInt8();
			case "int16":
				return state.readInt16();
			case "int32":
				return state.readInt32();
			case "int64":
				return state.readInt64();
			case "uint8":
				return state.readUInt8();
			case "uint16":
				return state.readUInt16();
			case "uint32":
				return state.readUInt32();
			case "uint64":
				return state.readUInt64();
		}
	}

	encode(value: any, data: Data) {
		switch (this.type) {
			case "array":
				return this.writeArray(value, data);
			case "object":
				if (!this.subFields)
					throw new Error('this.subFields is not defined');
				return this.subFields.write(value, data);
			default:
				try {
					this.writePrimitive(value, data);
				} catch (e: any){
					if (e instanceof TypeError) {
						throw new TypeError(`Error while encoding ${this.loc}: ${e.message}`);
					} else {
						throw e;
					}
				}
		}
	}

	private writePrimitive(value: any, data: Data) {
		switch (this.type) {
			case 'string':
				return data.writeString(value, this.loc, this.len, this.maxLen);
			case 'buffer':
				return data.writeBuffer(value, this.loc, this.len, this.maxLen);
			case 'boolean':
				return data.writeBoolean(value, this.loc);
			case 'float32':
				return data.writeFloat32(value, this.loc);
			case 'float64':
				return data.writeFloat64(value, this.loc);
			case 'int8':
				return data.writeInt8(value, this.loc);
			case 'int16':
				return data.writeInt16(value, this.loc);
			case 'int32':
				return data.writeInt32(value, this.loc);
			case 'int64':
				return data.writeInt64(value, this.loc);
			case 'uint8':
				return data.writeUInt8(value, this.loc);
			case 'uint16':
				return data.writeUInt16(value, this.loc);
			case 'uint32':
				return data.writeUInt32(value, this.loc);
			case 'uint64':
				return data.writeUInt64(value, this.loc);
		}
	}

	private writeArray(value: any, data: Data) {
		if (!Array.isArray(value) || !this.subType) {
			throw new TypeError('Expected an Array at ' + this.loc);
		}
		let arrLen = value.length;
		if (this.len) {
			if (arrLen !== this.len) {
				throw new TypeError(`Expected an Array of length ${this.len}, got ${arrLen}`);
			}
		} else if (this.maxLen) {
			if (arrLen > this.maxLen) {
				throw new TypeError(`Expected a string of length <= ${this.maxLen}, got ${arrLen}`);
			}
			data.writeUint(arrLen, minBytes(this.maxLen), 'array length');
		} else {
			data.writeUInt16(arrLen, 'array length');
		}
		for (let i = 0; i < arrLen; i ++) {
			this.subType.encode(value[i], data);
		}
	}

	private readArray(state: ReadState) {
		if (!this.subType)
			throw new Error('this.subType is not set');

		let length: number;
		if (this.len) {
			length = this.len;
		} else if (this.maxLen) {
			length = state.readUint(minBytes(this.maxLen));
		} else {
			length = state.readUInt16();
		}
		let arr = new Array(length);
		for (let j = 0; j < arr.length; j ++) {
			arr[j] = this.subType.decode(state);
		}
		return arr
	}
}
