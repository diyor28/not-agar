export type UIntT = 'uint8' | 'uint16' | 'uint32' | 'uint64';
export type IntT = 'int8' | 'int16' | 'int32' | 'int64';
export type FloatT = 'float32' | 'float64';
export type StringT = 'string';
export type BooleanT = 'boolean';
export type BufferT = 'buffer';
export type PrimitiveType = UIntT | IntT | FloatT | BooleanT | StringT | BufferT;
export type ExtendedPrimitiveType = PrimitiveType | 'array' | 'object';


export type BasicTypeMapping<T> = Record<string, PrimitiveType | PrimitiveType[] | T | T[]>;

export interface RecursiveTypeMapping extends BasicTypeMapping<RecursiveTypeMapping> {
}
export type FieldType = PrimitiveType | PrimitiveType[] | RecursiveTypeMapping | Array<RecursiveTypeMapping>
