import {
	FieldType,
	isFixedSize,
	isFixedSizeTypeConf,
	isTypeConf,
	isVarSize,
	isVarSizeTypeConf,
	SchemaType,
	StrictSchemaType,
	StrictTypeConf
} from './types';
import {FieldsMap} from "./field";
import Data from "./data";
import ReadState from "./readState";


const schema: StrictSchemaType = {
	x: {type: 'uint32', optional: false},
	players: {
		type: 'array',
		of: {x: {type: 'uint32', optional: false}},
		optional: false,
		maxLen: 0,
		length: 0
	},
	food: {
		type: 'array',
		of: {type: 'float32', optional: false},
		optional: false,
		maxLen: 0,
		length: 0
	},
	stats: {
		type: 'object',
		of: {x: {type: 'float32', optional: false}},
		optional: false
	}
}

function fieldToConf(field: FieldType): StrictTypeConf<StrictSchemaType> {
	if (isFixedSize(field)) {
		return {type: field, optional: false};
	}
	if (isVarSize(field)) {
		return {type: field, optional: false, length: 0, maxLen: 0};
	}
	if (isFixedSizeTypeConf(field))
		return {type: field.type, optional: field.optional || false}

	if (isVarSizeTypeConf(field)) {
		return {
			type: field.type,
			optional: field.optional || false,
			length: field.length || 0,
			maxLen: field.maxLen || 0
		}
	}
	if (isTypeConf(field)) {
		if (field.type === 'object') {
			const res: StrictSchemaType = {};
			Object.keys(field.of).forEach(key => {
				res[key] = fieldToConf(field.of[key]);
			});
			return {type: field.type, of: res, optional: false};
		}
		return {
			type: field.type,
			of: fieldToConf(field.of),
			optional: false,
			length: field.length || 0,
			maxLen: field.length || 0
		};
	}

	const result: StrictSchemaType = {};
	Object.keys(field).forEach(key => {
		result[key] = fieldToConf(field[key])
	});
	return {
		type: 'object',
		of: result,
		optional: false
	};
}

function strictSchema(schema: SchemaType): StrictSchemaType {
	let result: StrictSchemaType = {};
	Object.keys(schema).forEach(key => {
		result[key] = fieldToConf(schema[key]);
	});
	return result;
}

export default class Schema {
	fields: FieldsMap;
	schema: SchemaType;

	constructor(schema: SchemaType) {
		if (typeof schema !== 'object') {
			throw new TypeError('Invalid type: ' + schema)
		}
		this.fields = new FieldsMap('', strictSchema(schema));
		this.schema = schema;
	}

	encode(value: any): Data {
		const data = new Data();
		this.fields.write(value, data);
		return data;
	}

	decode(buffer: Buffer) {
		return this.fields.read(new ReadState(buffer));
	}

	extends(schema: SchemaType) {
		return new Schema({...this.schema, ...schema})
	}
}

