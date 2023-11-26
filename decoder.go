package schematic

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
)

// Decoder struct which holds schema for the expected payload.
// SchemaDecoder compares paylod before decoding with the schema.
type SchemaDecoder[T any] struct {
	Schema  map[string]any
	Marsh   func(v any) ([]byte, error)
	UnMarsh func(data []byte, v any) error
	Decoder func(r io.Reader) *json.Decoder
}

// Creates schema of the given model
func (self *SchemaDecoder[T]) createSchema(m any) (map[string]any, error) {
	schema := make(map[string]any)

	if b, err := self.Marsh(m); err != nil {
		return schema, err

	} else {
		return schema, self.UnMarsh(b, &schema)
	}
}

// Creates new schema decoder.
func NewSchemaDecoder[T any](
	marsh func(v any) ([]byte, error),
	unmarsh func(data []byte, v any) error,
	decoder func(r io.Reader) *json.Decoder,

) *SchemaDecoder[T] {

	self := &SchemaDecoder[T]{
		Marsh:   marsh,
		UnMarsh: unmarsh,
		Decoder: decoder,
	}

	b, err := self.createSchema(new(T))
	if err != nil {
		return nil
	}

	self.Schema = b
	return self
}

// Decodes payload with schema comparasing
func (self *SchemaDecoder[T]) Decode(raw []byte) *T {
	r := make(map[string]any)

	if err := self.Decoder(bytes.NewReader(raw)).Decode(&r); err != nil {
		return nil
	}

	if isMapSchemaValid(self.Schema, r) {
		s := new(T)
		if err := json.Unmarshal(raw, s); err != nil {
			return nil
		}

		return s

	} else {
		return nil
	}
}

// Makes a deep validation of the data with schema.
// Returns true, if all fields in data are all presented and
// have proper types.
func isMapSchemaValid(schema, data map[string]interface{}) bool {
	for key, value := range schema {
		if dataValue, exists := data[key]; exists {
			if reflect.TypeOf(value) != reflect.TypeOf(dataValue) {
				return false
			}

			switch value.(type) {
			case map[string]interface{}:
				if dataMapValue, ok := dataValue.(map[string]interface{}); ok {
					if !isMapSchemaValid(value.(map[string]interface{}), dataMapValue) {
						return false
					}
				} else {
					return false
				}
			}

			// Ignore unknown keys in the data map . . .
		} else {
			return false
		}
	}

	for key := range data {
		if _, exists := schema[key]; !exists {
			return false
		}
	}

	return true
}
