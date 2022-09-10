package csbin

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type PrimitiveType string

const (
	String  PrimitiveType = "string"
	Uint8                 = "uint8"
	Uint16                = "uint16"
	Uint32                = "uint32"
	Uint64                = "uint64"
	Int8                  = "int8"
	Int16                 = "int16"
	Int32                 = "int32"
	Int64                 = "int64"
	Float32               = "float32"
	Float64               = "float64"
	Array                 = "array"
	Struct                = "struct"
)

type Field struct {
	Type      PrimitiveType
	SubType   *Field
	SubFields Fields
}

func (f *Field) Decode(value *reflect.Value, reader *BytesReader) error {
	switch value.Kind() {
	case reflect.String:
		if !f.IsString() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: string", f.Type))
		}
		s, err := reader.ReadString()
		if err != nil {
			return err
		}
		value.SetString(s)

	case reflect.Uint:
		if !f.IsUint() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: uint", f.Type))
		}
		u, err := reader.ReadUint(f.Size())
		if err != nil {
			return err
		}
		value.SetUint(u)
	case reflect.Int:
		if !f.IsInt() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: int", f.Type))
		}
		i, err := reader.ReadInt(f.Size())
		if err != nil {
			return err
		}
		value.SetInt(i)
	case reflect.Float32:
	case reflect.Float64:
		if !f.IsFloat() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: float", f.Type))
		}
		i, err := reader.ReadFloat(f.Size())
		if err != nil {
			return err
		}
		value.SetFloat(i)
	case reflect.Struct:
		if !f.IsStruct() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: struct", f.Type))
		}
		f.SubFields
	case reflect.Slice:
		if !f.IsArray() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: slice", f.Type))
		}
	}
	return errors.New("unexpected type")
}

func (f *Field) Encode(value *reflect.Value, writer *BytesWriter) error {
	switch value.Kind() {
	case reflect.String:
		if !f.IsString() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: string", f.Type))
		}
		writer.Write(value.String())
	case reflect.Uint:
		if !f.IsUint() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: uint", f.Type))
		}
		writer.Write(value.Uint())
	case reflect.Int:
		if !f.IsInt() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: int", f.Type))
		}
		writer.Write(value.Int())
	case reflect.Float32:
	case reflect.Float64:
		if !f.IsFloat() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: float", f.Type))
		}
		writer.Write(value.Float())
	case reflect.Slice:
		if !f.IsArray() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: slice", f.Type))
		}
		err := f.EncodeArray(value, writer)
		if err != nil {
			return err
		}
	case reflect.Struct:
		if !f.IsStruct() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: struct", f.Type))
		}
		err := f.SubFields.Encode(value, writer)
		if err != nil {
			return err
		}
	}
	return errors.New("unexpected type")
}

func (f *Field) EncodeArray(value *reflect.Value, writer *BytesWriter) error {
	writer.WriteUint16(uint16(value.Len()))
	for i := 0; i < value.Len(); i++ {
		el := value.Index(i)
		err := f.SubType.Encode(&el, writer)
		if err != nil {
			return errors.New(fmt.Sprintf("At %d ", i) + err.Error())
		}
	}
	return nil
}

func (f *Field) IsStruct() bool {
	return f.Type == Struct
}

func (f *Field) IsArray() bool {
	return f.Type == Array
}

func (f Field) IsString() bool {
	return f.Type == String
}

func (f *Field) IsUint() bool {
	return f.Type == Uint8 || f.Type == Uint16 || f.Type == Uint32 || f.Type == Uint64
}

func (f *Field) IsInt() bool {
	return f.Type == Int8 || f.Type == Int16 || f.Type == Int32 || f.Type == Int64
}

func (f *Field) IsFloat() bool {
	return f.Type == Float32 || f.Type == Float64
}

func (f *Field) Size() int {
	if f.Type == Uint8 || f.Type == Int8 {
		return 1
	}

	if f.Type == Uint16 || f.Type == Int16 {
		return 2
	}

	if f.Type == Uint32 || f.Type == Int32 || f.Type == Float32 {
		return 4
	}

	if f.Type == Uint64 || f.Type == Int64 || f.Type == Float64 {
		return 8
	}

	panic("type String has no size")
}

type Fields map[string]*Field

func (f *Fields) Encode(reflection *reflect.Value, writer *BytesWriter) error {
	for key, field := range *f {
		fieldName := strings.Title(key)
		value := reflection.FieldByName(fieldName)
		if !value.IsValid() {
			return errors.New(fmt.Sprintf("field %s is not valid", fieldName))
		}
		err := field.Encode(&value, writer)
		if err != nil {
			return errors.New(fmt.Sprintf("for field %s ", fieldName) + err.Error())
		}
	}
	return nil
}

func (f *Fields) Decode(reflection *reflect.Value, reader *BytesReader) error {
	bitmap, err := reader.ReadBitmap()
	if err != nil {
		return err
	}
	idx := 0
	for key, field := range *f {
		if bitmap[idx] == 1 {
			continue
		}
		fieldName := strings.Title(key)
		value := reflection.FieldByName(fieldName)

		if !value.IsValid() {
			return errors.New(fmt.Sprintf("field %s is not valid", fieldName))
		}
		if !value.CanSet() {
			return errors.New(fmt.Sprintf("field %s is not writeable", fieldName))
		}
		err := field.Decode(&value, reader)
		if err != nil {
			return err
		}
		idx++
	}
}
