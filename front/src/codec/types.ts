export type UIntT = 'uint8' | 'uint16' | 'uint32' | 'uint64';
export type IntT = 'int8' | 'int16' | 'int32' | 'int64';
export type FloatT = 'float32' | 'float64';
export type StringT = 'string';
export type BooleanT = 'boolean';
export type BufferT = 'buffer';
export type FixedSizePrimitive = UIntT | IntT | FloatT | BooleanT
export type VarSizePrimitive = UIntT | IntT | FloatT | BooleanT | StringT | BufferT
export type PrimitiveType = FixedSizePrimitive | VarSizePrimitive;
export type ExtendedPrimitiveType = PrimitiveType | 'array' | 'object';

interface ObjectTypeConf<T> {
	type: 'object'
	of: T
	optional?: boolean
}

interface ArrayTypeConf<T> {
	type: 'array'
	of: PrimitiveType | T | TypeConf<T>
	optional?: boolean
	length?: number
	maxLen?: number
}

interface StrictObjectTypeConf<T> {
	type: 'object'
	of: T
	optional: boolean
}

interface StrictArrayTypeConf<T> {
	type: 'array'
	of: T | StrictTypeConf<T>
	optional: boolean
	length: number
	maxLen: number
}

interface FixedSizeTypeConf {
	type: UIntT | IntT | FloatT | BooleanT
	optional?: boolean
}

interface VarSizeTypeConf {
	type: StringT | BufferT
	optional?: boolean
	length?: number
	maxLen?: number
}

export type TypeConf<T> = ObjectTypeConf<T> | ArrayTypeConf<T> | FixedSizeTypeConf | VarSizeTypeConf;
export type StrictTypeConf<T> = StrictObjectTypeConf<T> | StrictArrayTypeConf<T> | Required<FixedSizeTypeConf> | Required<VarSizeTypeConf>;
export type TypeMapping<T> = Record<string, PrimitiveType | TypeConf<T> | T>;
export type StrictTypeMapping<T> = Record<string, StrictTypeConf<T>>;

export interface SchemaType extends TypeMapping<SchemaType> {
}

export interface StrictSchemaType extends StrictTypeMapping<StrictSchemaType> {

}

export type FieldType = PrimitiveType | TypeConf<SchemaType> | SchemaType;
export type StrictFieldType = StrictTypeConf<StrictSchemaType> | StrictSchemaType;

export function isFixedSize(field: any): field is FixedSizePrimitive {
	if (typeof field !== 'string')
		return false;
	return ['uint8', 'uint16', 'uint32', 'uint64', 'int8', 'int16', 'int32', 'int64', 'float32', 'float64'].includes(field);
}

export function isVarSize(field: any): field is VarSizePrimitive {
	if (typeof field !== 'string')
		return false;
	return ['string', 'buffer'].includes(field);
}

export function isTypeConf(field: any): field is TypeConf<any> {
	return typeof field === 'object' && field.type
}

export function isFixedSizeTypeConf(field: any): field is FixedSizeTypeConf {
	return typeof field === 'object' && isFixedSize(field.type)
}

export function isVarSizeTypeConf(field: any): field is VarSizeTypeConf {
	return typeof field === 'object' && isVarSize(field.type)
}
