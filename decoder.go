package schematic

import (
	"bytes"

	"reflect"

	"github.com/goccy/go-json"
)

// Decoder struct which holds schema for the expected payload.
// SchemaDecoder compares paylod before decoding with the schema.
type SchemaDecoder[T any] struct {
	Schema map[string]any
}

// Creates schema of the given model
func createSchema(m any) (map[string]any, error) {
	schema := make(map[string]any)

	if b, err := json.Marshal(m); err != nil {
		return schema, err

	} else {
		return schema, json.Unmarshal(b, &schema)
	}
}

// Creates new schema decoder.
func NewSchemaDecoder[T any]() *SchemaDecoder[T] {

	b, err := createSchema(new(T))
	if err != nil {
		return nil
	}

	return &SchemaDecoder[T]{
		Schema: b,
	}
}

// Decodes payload with schema comparasing
func (self *SchemaDecoder[T]) Decode(raw []byte) *T {
	r := make(map[string]any)

	if err := json.NewDecoder(bytes.NewReader(raw)).Decode(&r); err != nil {
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
