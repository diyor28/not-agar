package csbin

import (
	"errors"
	"fmt"
	"github.com/diyor28/not-agar/src/csbin/bytesIO"
	"reflect"
)

type Schema struct {
	Fields   Fields
	compress bool
}

func New(fields ...*Field) *Schema {
	return &Schema{Fields: fields}
}

func (s *Schema) UseCompression() *Schema {
	s.compress = true
	return s
}

func (s *Schema) Extends(fields ...*Field) *Schema {
	combinedFields := s.Fields
	for _, f := range fields {
		combinedFields = append(combinedFields, f)
	}
	schema := New(combinedFields...)
	schema.compress = s.compress
	return schema
}

func (s *Schema) Add(fields ...*Field) *Schema {
	s.Fields = append(s.Fields, fields...)
	return s
}

func (s *Schema) NewField(name string, primitiveType reflect.Kind) *Field {
	field := &Field{Name: name, loc: name, Type: primitiveType}
	s.Fields = append(s.Fields, field)
	return field
}

func (s *Schema) Encode(data interface{}) (*bytesIO.BytesWriter, error) {
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
	if s.compress {
		if err := writer.Compress(); err != nil {
			return nil, err
		}
	}
	return writer, nil
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
	if s.compress {
		if err := reader.Decompress(); err != nil {
			return err
		}
	}
	return s.Fields.Decode(&reflection, reader)
}
