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
	Name      string
	Type      reflect.Kind
	SubType   *Field
	SubFields Fields
}

func (f *Fields) Encode(reflection *reflect.Value, writer *bytesIO.BytesWriter) error {
	bMask := bitmask.New()
	for i, field := range *f {
		var fieldName string
		var value reflect.Value
		if reflection.Kind() == reflect.Map {
			fieldName = field.Name
			value = reflection.MapIndex(reflect.ValueOf(fieldName))
		} else {
			fieldName = strings.Title(field.Name)
			value = reflection.FieldByName(fieldName)
		}
		if value.Kind() != 0 {
			bMask.Set(i, true)
		}
	}
	writer.WriteBytes(bMask.ToBytes())

	for _, field := range *f {
		var fieldName string
		var value reflect.Value
		if reflection.Kind() == reflect.Map {
			fieldName = field.Name
			value = reflection.MapIndex(reflect.ValueOf(fieldName))
		} else {
			fieldName = strings.Title(field.Name)
			value = reflection.FieldByName(fieldName)
		}
		if !value.IsValid() {
			fmt.Println(value.Kind())
			return errors.New(fmt.Sprintf("field %s is not valid", fieldName))
		}
		err := field.Encode(&value, writer)
		if err != nil {
			return errors.New(fmt.Sprintf("for field %s ", fieldName) + err.Error())
		}
	}
	return nil
}

func (f *Fields) Decode(reflection *reflect.Value, reader *bytesIO.BytesReader) error {
	bMask, err := reader.ReadBitmask()
	if err != nil {
		return err
	}
	for i, field := range *f {
		if !bMask.Has(i) {
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
			return err
		}
		if reflection.Kind() == reflect.Map {
			reflection.SetMapIndex(reflect.ValueOf(fieldName), value)
		}
	}
	return nil
}

func (f *Field) Decode(value *reflect.Value, reader *bytesIO.BytesReader) error {
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
		return nil
	case reflect.Uint:
		if !f.IsUint() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: uint", f.Type))
		}
		u, err := reader.ReadUint(f.Size())
		if err != nil {
			return err
		}
		value.SetUint(u)
		return nil
	case reflect.Int:
		if !f.IsInt() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: int", f.Type))
		}
		i, err := reader.ReadInt(f.Size())
		if err != nil {
			return err
		}
		value.SetInt(i)
		return nil
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
		return nil
	case reflect.Struct:
		if !f.IsStruct() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: struct", f.Type))
		}
		err := f.SubFields.Decode(value, reader)
		if err != nil {
			return err
		}
		return nil
	case reflect.Slice:
		if !f.IsSlice() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: slice", f.Type))
		}
		err := f.DecodeArray(value, reader)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("unexpected type: " + value.Kind().String())
}

func (f *Field) Encode(value *reflect.Value, writer *bytesIO.BytesWriter) error {
	switch value.Kind() {
	case reflect.String:
		if !f.IsString() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: string", f.Type))
		}
		writer.Write(value.String())
		return nil
	case reflect.Uint:
		if !f.IsUint() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: uint", f.Type))
		}
		writer.Write(value.Uint())
		return nil
	case reflect.Int:
		if !f.IsInt() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: int", f.Type))
		}
		writer.Write(value.Int())
		return nil
	case reflect.Float32:
	case reflect.Float64:
		if !f.IsFloat() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: float", f.Type))
		}
		writer.Write(value.Float())
		return nil
	case reflect.Slice:
		if !f.IsSlice() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: slice", f.Type))
		}
		err := f.EncodeArray(value, writer)
		if err != nil {
			return err
		}
		return nil
	case reflect.Struct:
		if !f.IsStruct() {
			return errors.New(fmt.Sprintf("expected type: %s. Got: struct", f.Type))
		}
		err := f.SubFields.Encode(value, writer)
		if err != nil {
			return err
		}
		return nil
	case reflect.Interface:
		copyValue := reflect.ValueOf(value.Interface())
		return f.Encode(&copyValue, writer)
	}
	return errors.New("unexpected type: " + value.Kind().String())
}

func (f *Field) DecodeArray(value *reflect.Value, reader *bytesIO.BytesReader) error {
	arrLength, err := reader.ReadUint16()
	if err != nil {
		return err
	}
	for i := 0; i < int(arrLength); i++ {
		el := value.Index(i)
		err := f.SubType.Decode(&el, reader)
		if err != nil {
			return errors.New(fmt.Sprintf("At %d ", i) + err.Error())
		}
	}
	return nil
}

func (f *Field) EncodeArray(value *reflect.Value, writer *bytesIO.BytesWriter) error {
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
	}
	panic(fmt.Sprintf("could not convert type %s to reflect.Type", f.Type.String()))
}

func (f *Field) IsStruct() bool {
	return f.Type == reflect.Struct
}

func (f *Field) IsSlice() bool {
	return f.Type == reflect.Slice
}

func (f Field) IsString() bool {
	return f.Type == reflect.String
}

func (f *Field) IsUint() bool {
	return f.Type == reflect.Uint8 || f.Type == reflect.Uint16 || f.Type == reflect.Uint32 || f.Type == reflect.Uint64
}

func (f *Field) IsInt() bool {
	return f.Type == reflect.Int8 || f.Type == reflect.Int16 || f.Type == reflect.Int32 || f.Type == reflect.Int64
}

func (f *Field) IsFloat() bool {
	return f.Type == reflect.Float32 || f.Type == reflect.Float64
}

func (f *Field) Size() int {
	switch f.Type {
	case reflect.Bool:
	case reflect.Uint8:
	case reflect.Int8:
		return 1
	case reflect.Uint16:
	case reflect.Int16:
		return 2
	case reflect.Uint32:
	case reflect.Int32:
	case reflect.Float32:
		return 4
	case reflect.Uint64:
	case reflect.Int64:
	case reflect.Float64:
		return 8
	}
	panic("type String has no size")
}

type Fields []*Field
