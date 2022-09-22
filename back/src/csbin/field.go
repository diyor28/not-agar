package csbin

import (
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/csbin/bitmask"
	"github.com/diyor28/not-agar/src/csbin/bytesIO"
	"reflect"
	"strings"
)

type Field struct {
	Name       string
	Type       reflect.Kind
	optional   bool
	loc        string
	structType *reflect.Type
	subType    *Field
	subFields  Fields
	maxLen     uint64
	len        uint64
}

func NewField(name string, primitiveType reflect.Kind) *Field {
	return &Field{Name: name, loc: name, Type: primitiveType}
}

func (f *Field) Len(exactLen uint64) *Field {
	if f.Type != reflect.Slice && f.Type != reflect.String {
		panic(fmt.Sprintf("type %s does not support Len()", f.Type.String()))
	}
	f.len = exactLen
	return f
}

func (f *Field) Optional() *Field {
	f.optional = true
	return f
}

func (f *Field) MaxLen(maxLen uint64) *Field {
	if f.Type != reflect.Slice && f.Type != reflect.String {
		panic(fmt.Sprintf("type %s does not support MaxLen()", f.Type.String()))
	}
	f.maxLen = maxLen
	return f
}

func (f *Field) UseStruct(s interface{}) *Field {
	if f.Type != reflect.Struct {
		panic(fmt.Sprintf("type %s does not support UseStruct()", f.Type.String()))
	}
	structType := reflect.TypeOf(s)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected: struct, got: %s", structType.Kind().String()))
	}
	f.structType = &structType
	return f
}

func (f *Field) SubType(field *Field) *Field {
	if f.Type != reflect.Slice && f.Type != reflect.Array {
		panic(fmt.Sprintf("type %s does not support SubType()", f.Type.String()))
	}
	field.loc = f.loc + "." + field.Name
	f.subType = field
	return f
}

func (f *Field) SubFields(fields ...*Field) *Field {
	if f.Type != reflect.Struct && f.Type != reflect.Map {
		panic(fmt.Sprintf("type %s does not support SubFields()", f.Type.String()))
	}
	for _, field := range fields {
		field.loc = f.loc + "." + field.Name
		f.subFields = append(f.subFields, field)
	}
	return f
}

func (f *Fields) writeBitmask(reflection *reflect.Value, writer *bytesIO.BytesWriter) error {
	bMask := bitmask.New()
	for _, field := range *f {
		var fieldName string
		var value reflect.Value
		if reflection.Kind() == reflect.Map {
			fieldName = field.Name
			mapEl := reflection.MapIndex(reflect.ValueOf(fieldName))
			if !mapEl.IsValid() {
				return errors.New(fmt.Sprintf("key %s does not exist", fieldName))
			}
			value = mapEl.Elem()
		} else {
			fieldName = strings.Title(field.Name)
			value = reflection.FieldByName(fieldName)
		}
		if value.IsZero() {
			bMask.Set(false)
		} else {
			bMask.Set(true)
		}
	}
	writer.WriteBytes(bMask.ToBytes(), "bitmask")
	return nil
}

func (f *Fields) hasOptionalFields() bool {
	for _, field := range *f {
		if field.optional {
			return true
		}
	}
	return false
}

func (f *Fields) Encode(reflection *reflect.Value, writer *bytesIO.BytesWriter) error {
	if f.hasOptionalFields() {
		if err := f.writeBitmask(reflection, writer); err != nil {
			return err
		}
	}

	for _, field := range *f {
		var fieldName string
		var value reflect.Value
		if reflection.Kind() == reflect.Map {
			fieldName = field.Name
			value = reflection.MapIndex(reflect.ValueOf(fieldName)).Elem()
		} else {
			fieldName = strings.Title(field.Name)
			value = reflection.FieldByName(fieldName)
		}
		if value.IsZero() {
			continue
		}
		if !value.IsValid() {
			return errors.New(fmt.Sprintf("field %s is not valid", fieldName))
		}
		err := field.Encode(&value, writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Fields) Decode(reflection *reflect.Value, reader *bytesIO.BytesReader) error {
	var bMask *bitmask.Bitmask
	var err error
	if f.hasOptionalFields() {
		bMask, err = reader.ReadBitmask()
		if err != nil {
			return err
		}
	}
	for i, field := range *f {
		if bMask != nil && !bMask.Has(i, len(*f)) {
			continue
		}
		var fieldName string
		var value reflect.Value
		if reflection.Kind() == reflect.Map {
			fieldName = field.Name
			value = reflect.New(field.ConstructType()).Elem()
		} else {
			fieldName = strings.Title(field.Name)
			value = reflection.FieldByName(fieldName)
		}

		if !value.IsValid() {
			return errors.New(fmt.Sprintf("field %s is not valid", fieldName))
		}
		if !value.CanSet() {
			return errors.New(fmt.Sprintf("field %s is not writeable", fieldName))
		}
		err := field.Decode(&value, reader)
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s", field.Name, err.Error()))
		}
		if reflection.Kind() == reflect.Map {
			reflection.SetMapIndex(reflect.ValueOf(fieldName), value)
		}
	}
	return nil
}

func (f *Field) Decode(value *reflect.Value, reader *bytesIO.BytesReader) error {
	if f.Type != value.Kind() {
		return errors.New(fmt.Sprintf("at %s expected: %s, got: %s", f.loc, f.Type, value.Kind()))
	}
	switch value.Kind() {
	case reflect.String:
		s, err := reader.ReadString()
		if err != nil {
			return err
		}
		value.SetString(s)
		return nil
	case reflect.Bool:
		b, err := reader.ReadBool()
		if err != nil {
			return err
		}
		value.SetBool(b)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := reader.ReadUint(f.Size())
		if err != nil {
			return err
		}
		value.SetUint(u)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := reader.ReadInt(f.Size())
		if err != nil {
			return err
		}
		value.SetInt(i)
		return nil
	case reflect.Float32, reflect.Float64:
		i, err := reader.ReadFloat(f.Size())
		if err != nil {
			return err
		}
		value.SetFloat(i)
		return nil
	case reflect.Struct:
		err := f.subFields.Decode(value, reader)
		if err != nil {
			return err
		}
		return nil
	case reflect.Slice:
		err := f.DecodeArray(value, reader)
		if err != nil {
			return err
		}
		return nil
	case reflect.Map:
		err := f.subFields.Decode(value, reader)
		if err != nil {
			return err
		}
	}
	return errors.New(fmt.Sprintf("type %s is not supported", value.Kind().String()))
}

func (f *Field) Encode(value *reflect.Value, writer *bytesIO.BytesWriter) error {
	if value.Kind() == reflect.Ptr {
		elem := value.Elem()
		return f.Encode(&elem, writer)
	}
	if f.Type != value.Kind() {
		return errors.New(fmt.Sprintf("at %s expected: %s, got: %s", f.loc, f.Type, value.Kind()))
	}
	switch value.Kind() {
	case reflect.String:
		writer.Write(value.String(), f.loc)
		return nil
	case reflect.Bool:
		writer.WriteBool(value.Bool(), f.loc)
		return nil
	case reflect.Uint:
		writer.Write(uint(value.Uint()), f.loc)
		return nil
	case reflect.Uint8:
		writer.Write(uint8(value.Uint()), f.loc)
		return nil
	case reflect.Uint16:
		writer.Write(uint16(value.Uint()), f.loc)
		return nil
	case reflect.Uint32:
		writer.Write(uint32(value.Uint()), f.loc)
		return nil
	case reflect.Uint64:
		writer.Write(value.Uint(), f.loc)
		return nil
	case reflect.Int:
		writer.Write(value.Int(), f.loc)
		return nil
	case reflect.Int8:
		writer.Write(int8(value.Int()), f.loc)
		return nil
	case reflect.Int16:
		writer.Write(int16(value.Int()), f.loc)
		return nil
	case reflect.Int32:
		writer.Write(int32(value.Int()), f.loc)
		return nil
	case reflect.Int64:
		writer.Write(value.Int(), f.loc)
		return nil
	case reflect.Float32:
		writer.Write(float32(value.Float()), f.loc)
		return nil
	case reflect.Float64:
		writer.Write(value.Float(), f.loc)
		return nil
	case reflect.Slice:
		err := f.EncodeArray(value, writer)
		if err != nil {
			return err
		}
		return nil
	case reflect.Struct:
		err := f.subFields.Encode(value, writer)
		if err != nil {
			return err
		}
		return nil
	case reflect.Array:
		err := f.EncodeArray(value, writer)
		if err != nil {
			return err
		}
		return nil
	case reflect.Interface:
		copyValue := reflect.ValueOf(value.Interface())
		return f.Encode(&copyValue, writer)
	case reflect.Ptr:
		copyValue := value.Elem()
		return f.Encode(&copyValue, writer)
	}
	return errors.New(fmt.Sprintf("type %s is not supported", value.Kind().String()))
}

func (f *Field) DecodeArray(value *reflect.Value, reader *bytesIO.BytesReader) error {
	arrLength, err := reader.ReadUint16()
	if err != nil {
		return err
	}
	value.Set(reflect.MakeSlice(f.ConstructType(), int(arrLength), int(arrLength)))
	for i := 0; i < int(arrLength); i++ {
		el := reflect.New(f.subType.ConstructType()).Elem()
		err := f.subType.Decode(&el, reader)
		if err != nil {
			return errors.New(fmt.Sprintf("At %d ", i) + err.Error())
		}
		value.Index(i).Set(el)
	}
	return nil
}

func (f *Field) EncodeArray(value *reflect.Value, writer *bytesIO.BytesWriter) error {
	writer.WriteUint16(uint16(value.Len()), "array length")
	for i := 0; i < value.Len(); i++ {
		el := value.Index(i)
		err := f.subType.Encode(&el, writer)
		if err != nil {
			return errors.New(fmt.Sprintf("At %d ", i) + err.Error())
		}
	}
	return nil
}

func (f *Field) ConstructType() reflect.Type {
	switch f.Type {
	case reflect.Uint8:
		return reflect.TypeOf(uint8(0))
	case reflect.Uint16:
		return reflect.TypeOf(uint16(0))
	case reflect.Uint32:
		return reflect.TypeOf(uint32(0))
	case reflect.Uint64:
		return reflect.TypeOf(uint64(0))
	case reflect.Int8:
		return reflect.TypeOf(int8(0))
	case reflect.Int16:
		return reflect.TypeOf(int16(0))
	case reflect.Int32:
		return reflect.TypeOf(int32(0))
	case reflect.Int64:
		return reflect.TypeOf(int64(0))
	case reflect.Float32:
		return reflect.TypeOf(float32(0))
	case reflect.Float64:
		return reflect.TypeOf(float64(0))
	case reflect.String:
		return reflect.TypeOf("")
	case reflect.Bool:
		return reflect.TypeOf(false)
	case reflect.Slice:
		return reflect.SliceOf(f.subType.ConstructType())
	case reflect.Map:
		return reflect.TypeOf(map[string]interface{}{})
	case reflect.Struct:
		return *f.structType
	}
	panic(fmt.Sprintf("could not convert type %s to reflect.Type", f.Type.String()))
}

func (f *Field) Size() int {
	switch f.Type {
	case reflect.Bool, reflect.Uint8, reflect.Int8:
		return 1
	case reflect.Uint16, reflect.Int16:
		return 2
	case reflect.Uint32, reflect.Int32, reflect.Float32:
		return 4
	case reflect.Uint64, reflect.Int64, reflect.Float64:
		return 8
	}
	panic(fmt.Sprintf("type %s has no size", f.Type.String()))
}

type Fields []*Field
