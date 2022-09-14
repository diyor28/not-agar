package csbin

import (
	"errors"
	"github.com/diyor28/not-agar/src/csbin/bytesIO"
	"reflect"
)

type Schema struct {
	Fields Fields
}

func New(fields Fields) *Schema {
	return &Schema{Fields: fields}
}

func NewField(name string, primitiveType reflect.Kind, subType *Field, subFields Fields) *Field {
	return &Field{Name: name, Type: primitiveType, SubType: subType, SubFields: subFields}
}

func (s *Schema) Extends(fields Fields) *Schema {
	combinedFields := s.Fields
	for _, f := range fields {
		combinedFields = append(combinedFields, f)
	}
	return New(combinedFields)
}

func (s *Schema) Encode(data interface{}) ([]byte, error) {
	value := reflect.ValueOf(data).Elem()
	writer := bytesIO.NewWriter()
	if value.Kind() != reflect.Struct && value.Kind() != reflect.Map {
		return nil, errors.New("provided value is not a struct or a map")
	}
	err := s.Fields.Encode(&value, writer)
	if err != nil {
		return nil, err
	}
	return writer.Bytes, nil
}

func (s *Schema) Decode(data []byte, result interface{}) error {
	reflection := reflect.ValueOf(result).Elem()
	if reflection.Kind() != reflect.Struct && reflection.Kind() != reflect.Map {
		return errors.New("provided value is not a struct or a map")
	}
	reader := bytesIO.NewReader(data)
	return s.Fields.Decode(&reflection, reader)
}
