package structutil

import (
	"errors"
	"reflect"
)

// ParseStruct parses values from `from` to `to` based on matching JSON tags.
// Supports nested structs and arrays/slices of structs.
func ParseStruct(from, to interface{}) error {
	fromVal := reflect.ValueOf(from)
	toVal := reflect.ValueOf(to)

	// Ensure `to` is a pointer to a struct
	if toVal.Kind() != reflect.Ptr || toVal.Elem().Kind() != reflect.Struct {
		return errors.New("to must be a pointer to a struct")
	}

	fromType := reflect.TypeOf(from)
	toElem := toVal.Elem()

	// Create a map of JSON tags to field values from the `from` struct
	tagToField := buildTagMap(fromType, fromVal)

	// Set values to the `to` struct recursively
	return setStructValues(tagToField, toElem)
}

// buildTagMap creates a map of JSON tags to field values from a struct.
func buildTagMap(fromType reflect.Type, fromVal reflect.Value) map[string]reflect.Value {
	tagToField := make(map[string]reflect.Value)

	for i := 0; i < fromType.NumField(); i++ {
		field := fromType.Field(i)
		value := fromVal.Field(i)

		// Handle nested structs
		if field.Type.Kind() == reflect.Struct {
			nestedMap := buildTagMap(field.Type, value)
			for k, v := range nestedMap {
				tagToField[k] = v
			}
			continue
		}

		// Handle slices of structs
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
			tagToField[field.Tag.Get("json")] = value
			continue
		}

		// Map JSON tag to the field value
		tag := field.Tag.Get("json")
		if tag != "" {
			tagToField[tag] = value
		}
	}

	return tagToField
}

// setStructValues sets values to the `to` struct based on the tag map.
func setStructValues(tagToField map[string]reflect.Value, toElem reflect.Value) error {
	toType := toElem.Type()

	for i := 0; i < toType.NumField(); i++ {
		toField := toType.Field(i)
		toFieldValue := toElem.Field(i)
		toTag := toField.Tag.Get("json")

		// Handle nested structs in `to`
		if toFieldValue.Kind() == reflect.Struct {
			if nestedMap, ok := tagToField[toTag]; ok {
				err := setStructValues(buildTagMap(nestedMap.Type(), nestedMap), toFieldValue)
				if err != nil {
					return err
				}
			}
			continue
		}

		// Handle slices of structs in `to`
		if toFieldValue.Kind() == reflect.Slice && toFieldValue.Type().Elem().Kind() == reflect.Struct {
			fromSlice, exists := tagToField[toTag]
			if exists && fromSlice.Kind() == reflect.Slice {
				newSlice := reflect.MakeSlice(toFieldValue.Type(), fromSlice.Len(), fromSlice.Cap())

				for j := 0; j < fromSlice.Len(); j++ {
					fromElem := fromSlice.Index(j)
					toElem := reflect.New(toFieldValue.Type().Elem()).Elem()
					err := setStructValues(buildTagMap(fromElem.Type(), fromElem), toElem)
					if err != nil {
						return err
					}
					newSlice.Index(j).Set(toElem)
				}

				toFieldValue.Set(newSlice)
			}
			continue
		}

		// Set value if tags match and types are compatible
		if fromFieldValue, exists := tagToField[toTag]; exists {
			if toFieldValue.CanSet() && fromFieldValue.Type().AssignableTo(toFieldValue.Type()) {
				toFieldValue.Set(fromFieldValue)
			}
		}
	}

	return nil
}
