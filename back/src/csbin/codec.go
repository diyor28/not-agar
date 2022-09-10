package csbin

import (
	"errors"
	"reflect"
)

type Codec struct {
	Fields Fields
}

func New(fields Fields) *Codec {
	return &Codec{Fields: fields}
}

func (c *Codec) Encode(data interface{}) ([]byte, error) {
	reflection := reflect.ValueOf(data).Elem()
	writer := NewWriter()
	if reflection.Kind() != reflect.Struct {
		return nil, errors.New("provided value is not a struct")
	}
	err := c.Fields.Encode(&reflection, writer)
	if err != nil {
		return nil, err
	}
	return writer.Bytes, nil
}

func (c *Codec) Decode(data []byte, result interface{}) error {
	reflection := reflect.ValueOf(result).Elem()
	if reflection.Kind() != reflect.Struct {
		return errors.New("provided value is not a struct")
	}
	reader := NewReader(data)
	return c.Fields.Decode(&reflection, reader)
}
