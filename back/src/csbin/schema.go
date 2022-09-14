package csbin

import (
	"errors"
	"fmt"
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
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Ptr {
		return nil, errors.New(fmt.Sprintf("expected pointer to struct or map, got %s", value.Kind().String()))
	}
	value = value.Elem()
	writer := bytesIO.NewWriter()
	if value.Kind() != reflect.Struct && value.Kind() != reflect.Map {
		return nil, errors.New(fmt.Sprintf("expected struct or map, got %s", value.Kind().String()))
	}
	err := s.Fields.Encode(&value, writer)
	if err != nil {
		return nil, err
	}
	return writer.Bytes, nil
}

func (s *Schema) Decode(data []byte, result interface{}) error {
	reflection := reflect.ValueOf(result)
	if reflection.Kind() != reflect.Ptr {
		return errors.New(fmt.Sprintf("expected pointer to struct or map, got %s", reflection.Kind().String()))
	}
	reflection = reflection.Elem()
	if reflection.Kind() != reflect.Struct && reflection.Kind() != reflect.Map {
		return errors.New(fmt.Sprintf("expected struct or map, got %s", reflection.Kind().String()))
	}
	reader := bytesIO.NewReader(data)
	return s.Fields.Decode(&reflection, reader)
}
