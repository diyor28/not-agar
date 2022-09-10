import types, {ExtendedPrimitiveType, FieldType} from "./types";
import Data from "./data";

export default class Field {
	public name: string;
	public array: boolean;
	public type: ExtendedPrimitiveType
	public subType: Field | null = null;
	public subFields: Field[] = [];

	constructor(name: string, type: FieldType) {
		this.name = name;
		this.array = Array.isArray(type);

		if (Array.isArray(type)) {
			if (type.length !== 1) {
				throw new TypeError(`Invalid array type for ${name}, it must have exactly one element`);
			}
			this.type = 'array';
			this.subType = new Field(name, type[0]);
		} else if (typeof type === 'object') {
			this.type = 'object';
			this.subFields = Object.keys(type).map(key => {
				return new Field(name + '.' + key, type[key]);
			})
		} else {
			this.type = type;
		}
	}

	encode(value: any, data: Data) {
		switch (this.type) {
			case "array":

			case "boolean":
			case "buffer":
			case "date":
			case "float32":
			case "float64":
			case "int8":
			case "int16":
			case "int32":
			case "int64":
			case "json":
			case "object":
			case "oid":
			case "regex":
			case "string":
			case "uint8":
			case "uint16":
			case "uint32":
			case "uint64":

		}
	}

	private writeArray(value: any, data: Data) {
		let len = value.length;
		if (!Array.isArray(value)) {
			throw new TypeError('Expected an Array at ' + this.name)
		}
		types.uint.write(len, data)
		for (let i = 0; i < len; i ++) {
			type.write(value[i], data, path + '.' + i)
		}
	}
}
