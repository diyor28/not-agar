import {ExtendedPrimitiveType, FieldType, RecursiveTypeMapping} from "./types";
import Data from "./data";
import ReadState from "./readState";

export class FieldsMap {
	loc: string;
	fields: Field[];

	constructor(loc: string, type: RecursiveTypeMapping) {
		this.loc = loc;
		this.fields = Object.keys(type).map(key => {
			return new Field(key, type[key], loc ? loc + '.' + key : key);
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

	write(value: any, data: Data) {
		const {bitmask, bitmaskSize} = this.calcBitmask(value);
		data.writeUInt8(bitmaskSize, 'bitmask size');
		data.writeUint(bitmask, 'bitmask');
		for (const field of this.fields) {
			const subValue = value[field.name];
			if (subValue === undefined || subValue === null)
				continue;
			field.encode(value[field.name], data);
		}
	}

	read(state: ReadState) {
		const bitmaskBytes = state.readUInt8();
		let bitmask: number = 0;
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
		const result: Record<string, any> = {};
		this.fields.forEach((field, i) => {
			let name = field.name;
			if (this.isMaskTrue(bitmask, i)) {
				result[name] = field.decode(state);
			} else {
				result[name] = undefined;
			}
		});

		return result;
	}

	private isMaskTrue(mask: number, idx: number) {
		return (mask >> (this.fields.length - idx - 1)) & 1; // Shift right to put the target at position 0, and AND it with 1
	}
}

export default class Field {
	name: string;
	loc: string;
	array: boolean;
	type: ExtendedPrimitiveType
	subType: Field | null = null;
	subFields: FieldsMap | null = null;

	constructor(name: string, type: FieldType, loc: string) {
		this.name = name;
		this.loc = loc;
		this.array = Array.isArray(type);

		if (Array.isArray(type)) {
			if (type.length !== 1) {
				throw new TypeError(`Invalid array type for ${loc}, it must have exactly one element`);
			}
			this.type = 'array';
			this.subType = new Field(name, type[0], `${loc}[]`);
		} else if (typeof type === 'object') {
			this.type = 'object';
			this.subFields = new FieldsMap(loc, type);
		} else {
			this.type = type;
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
			case "boolean":
				return state.readBoolean();
			case "buffer":
				return state.readBuffer();
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
			case "string":
				return state.readString();
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
			case "boolean":
				return data.writeBoolean(value, this.loc);
			case "buffer":
				return data.writeBuffer(value, this.loc);
			case "float32":
				return data.writeFloat32(value, this.loc);
			case "float64":
				return data.writeFloat64(value, this.loc);
			case "int8":
				return data.writeInt8(value, this.loc);
			case "int16":
				return data.writeInt16(value, this.loc);
			case "int32":
				return data.writeInt32(value, this.loc);
			case "int64":
				return data.writeInt64(value, this.loc);
			case "string":
				return data.writeString(value, this.loc);
			case "uint8":
				return data.writeUInt8(value, this.loc);
			case "uint16":
				return data.writeUInt16(value, this.loc);
			case "uint32":
				return data.writeUInt32(value, this.loc);
			case "uint64":
				return data.writeUInt64(value, this.loc);
		}
	}

	private writeArray(value: any, data: Data) {
		if (!Array.isArray(value) || !this.subType) {
			throw new TypeError('Expected an Array at ' + this.name)
		}
		let len = value.length;
		data.writeUInt16(len, 'array length');
		for (let i = 0; i < len; i ++) {
			this.subType.encode(value[i], data);
		}
	}

	private readArray(state: ReadState) {
		if (!this.subType)
			throw new Error('this.subType is not set');
		let arr = new Array(state.readUInt16());
		for (let j = 0; j < arr.length; j ++) {
			arr[j] = this.subType.decode(state);
		}
		return arr
	}
}
