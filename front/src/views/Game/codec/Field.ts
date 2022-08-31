import Type from "./Type";

export default class Field {
	public optional = false;
	public name: string;
	public array: boolean;
	public type: Type

	constructor(name: string, type: any) {
		if (name[name.length - 1] === '?') {
			this.optional = true;
			name = name.substr(0, name.length - 1);
		}

		/** @member {string} */
		this.name = name;

		/** @member {boolean} */
		this.array = Array.isArray(type);

		if (Array.isArray(type)) {
			if (type.length !== 1) {
				throw new TypeError('Invalid array type, it must have exactly one element');
			}
			type = type[0];
		}

		this.type = new Type(type);
	}
}
