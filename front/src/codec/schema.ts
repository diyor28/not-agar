import {RecursiveTypeMapping} from './types';
import {FieldsMap} from "./field";
import Data from "./data";
import ReadState from "./readState";

export default class Schema {
	fields: FieldsMap;
	schema: RecursiveTypeMapping;

	constructor(schema: RecursiveTypeMapping) {
		if (typeof schema !== 'object') {
			throw new TypeError('Invalid type: ' + schema)
		}
		this.fields = new FieldsMap('', schema);
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

	extends(schema: RecursiveTypeMapping) {
		return new Schema({...this.schema, ...schema})
	}
}

