package context

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func extractAllParams(r *http.Request, target interface{}) error {
	params := make(map[string]interface{})
	targetType := reflect.TypeOf(target)

	// Process headers
	for name, values := range r.Header {
		if field, ok := findField(targetType, name); ok {
			params[name] = convertValues(values, field.Type)
		}
	}

	// Process URL query parameters
	//for name, values := range r.URL.Query() {
	//	field, ok := findField(targetType, name)
	//	if ok {
	//		params[name] = convertValues(values, field.Type)
	//	}
	//}

	// Process path variables
	if vars := mux.Vars(r); vars != nil {
		for name, value := range vars {
			if field, ok := findField(targetType, name); ok {
				params[name] = convertValue(value, field.Type)
			}
		}
	}

	// Process request body
	if r.Body != nil {
		// Check if target has a Body field
		if field, ok := findField(targetType, "Body"); ok {
			bodyValue := reflect.New(field.Type.Elem()).Interface()
			if err := json.NewDecoder(r.Body).Decode(bodyValue); err == nil {
				params["body"] = bodyValue
			}
		}
	}

	jsonContent, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return json.NewDecoder(bytes.NewBuffer(jsonContent)).Decode(&target)
}

// Helper function to find a field in the target struct (case insensitive)
func findField(t reflect.Type, name string) (reflect.StructField, bool) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == name {
			return field, true
		}
		// Check for json tag
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			jsonName := strings.Split(jsonTag, ",")[0]
			if jsonName == name {
				return field, true
			}
		}
	}
	return reflect.StructField{}, false
}

// Convert string values based on target type
func convertValues(values []string, targetType reflect.Type) interface{} {
	// Handle slice types
	if targetType.Kind() == reflect.Slice {
		elemType := targetType.Elem()
		slice := reflect.MakeSlice(targetType, len(values), len(values))
		for i, v := range values {
			converted := convertValue(v, elemType)
			slice.Index(i).Set(reflect.ValueOf(converted))
		}
		return slice.Interface()
	}

	// For non-slice types, use first value if available
	if len(values) > 0 {
		return convertValue(values[0], targetType)
	}
	return nil
}

// Convert single string value to target type
func convertValue(value string, targetType reflect.Type) interface{} {
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
	}
	// Fallback to string if conversion fails or type not handled
	return value
}
