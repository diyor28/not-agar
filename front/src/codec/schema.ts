import {RecursiveTypeMapping} from './types';
import {FieldsMap} from "./field";
import Data from "./data";
import ReadState from "./readState";

export default class Schema {
	public fields: FieldsMap;

	constructor(type: RecursiveTypeMapping) {
		if (typeof type !== 'object') {
			throw new TypeError('Invalid type: ' + type)
		}
		this.fields = new FieldsMap('', type);
	}

	encode(value: any): Data {
		const data = new Data();
		this.fields.write(value, data);
		return data;
	}

	decode(buffer: Buffer) {
		return this.fields.read(new ReadState(buffer));
	}
}

