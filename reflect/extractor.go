package reflect

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ParamsExtractorSource struct {
	// Source is a flat map where keys are parameter names and values are their corresponding string values.
	Source map[string]string
	// Tree represents a hierarchical map, primarily for parameters with multiple values, such as query parameters.
	Tree map[string][]string
}

// ParamsExtract extracts parameters from one or more ParamsExtractorSource instances,
// converts them to their target types, and binds them to the specified target struct.
// The function handles both single-value and multi-value parameters, supports JSON tags for field mapping,
// and ensures type conversions using ConvertValue and ConvertValues methods.
//
// target: The struct to which the extracted and converted parameters will be bound.
// source: A variadic list of ParamsExtractorSource instances, each representing a source of parameters.
//
// Returns an error if parameter extraction or JSON unmarshalling into the target struct fails.
func ParamsExtract(target interface{}, source ...ParamsExtractorSource) error {
	targetType := reflect.TypeOf(target)

	jsonContent, err := ParseToJson(targetType, source...)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonContent, target)
}

func FormBodyExtract(target interface{}, source ...ParamsExtractorSource) {
	targetType := reflect.TypeOf(target).Elem()

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		if field.Name == "Body" {
			err := json.Unmarshal([]byte(`{"body": {}}`), &target)
			if err != nil {
				break
			}

			fieldInstance := reflect.ValueOf(target).Elem().FieldByName("Body").Interface()
			targetBodyType := reflect.TypeOf(fieldInstance)

			jsonContent, err := ParseToJson(targetBodyType, source...)

			if err != nil {
				log.Errorf("FormBodyExtract: %v", err)
				break
			}

			err = json.Unmarshal([]byte(fmt.Sprintf(`{"body": %s}`, jsonContent)), &target)

			if err != nil {
				log.Errorf("FormBodyExtract: %v", err)
			}

			break
		}
	}
}

func ParseToJson(
	targetType reflect.Type,
	source ...ParamsExtractorSource,
) ([]byte, error) {
	params := make(map[string]interface{})

	for _, values := range source {
		if values.Tree != nil {
			for name, value := range values.Tree {
				if field, ok := FindField(targetType, name); ok {
					params[name] = ConvertValues(value, field)
				}
			}
		}

		if values.Source != nil {
			for name, value := range values.Source {
				if field, ok := FindField(targetType, name); ok {
					params[name] = ConvertValue(value, field.Type)
				}
			}
		}
	}
	return json.Marshal(params)
}

// FindField identifies and retrieves a struct field based on its name or its JSON tag.
//
// t: The type of the struct to search within.
// name: The name of the field or its JSON tag.
//
// Returns the struct field and a boolean indicating whether the field was found.
func FindField(t reflect.Type, name string) (reflect.StructField, bool) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if strings.EqualFold(field.Name, name) {
			return field, true
		}
		// Check for json tag
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			jsonName := strings.Split(jsonTag, ",")[0]
			if strings.EqualFold(jsonName, name) {
				return field, true
			}
		}
	}
	return reflect.StructField{}, false
}

// ConvertValues converts a slice of string values to a type matching the specified target type.
//
// values: A slice of string values to be converted.
// targetType: The desired type to which the values should be converted.
//
// Returns the converted value, which may be a slice or a single value, depending on the target type.
func ConvertValues(values []string, field reflect.StructField) interface{} {
	targetType := field.Type
	// Handle slice types
	if targetType.Kind() == reflect.Slice || targetType.Kind() == reflect.Array {
		elemType := targetType.Elem()
		var valuesResult []interface{}
		for _, v := range values {
			converted := ConvertValue(v, elemType)
			valuesResult = append(valuesResult, converted)
		}
		if valuesResult == nil {
			return nil
		}
		return valuesResult
	}

	// For non-slice types, use first value if available
	if len(values) > 0 {
		return ConvertValue(values[0], targetType)
	}
	return nil
}

// ConvertValue converts a single string value to a type matching the specified target type.
//
// value: The string value to be converted.
// targetType: The desired type to which the value should be converted.
//
// Returns the converted value or the original string value if no suitable conversion is possible.
func ConvertValue(value string, targetType reflect.Type) interface{} {
	switch targetType.Kind() {
	case reflect.String:
		return value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			return uintValue
		}
	case reflect.Float32, reflect.Float64:
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	case reflect.Bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	case reflect.Ptr, reflect.Struct:
		// Handle time.Time specifically
		if targetType == reflect.TypeOf(time.Time{}) {
			if t, err := time.Parse("2006-01-02", value); err == nil {
				return &t
			}
			// Add more time formats if needed
			if t, err := time.Parse(time.RFC3339, value); err == nil {
				return &t
			}
		} else {
			fieldValue := reflect.New(targetType).Interface()
			err := json.Unmarshal([]byte(value), &fieldValue)

			if err == nil {
				return fieldValue
			}
		}
	default:
		// No conversion possible
	}

	// Fallback: return the original value as string
	// You might want to handle this differently depending on your needs
	return value
}
