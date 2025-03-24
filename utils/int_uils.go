package utils

import "strconv"

func ToIntOrElseNil(value *string) *int {
	if value == nil || *value == "" {
		return nil
	}

	if intValue, err := strconv.Atoi(*value); err == nil {
		return &intValue
	}

	return nil
}

func ToIntOrElse(value *string, defaultValue int) *int {
	if intValue := ToIntOrElseNil(value); intValue != nil {
		return intValue
	}
	return &defaultValue
}
