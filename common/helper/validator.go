package helper

import (
	"fmt"
	"reflect"
)

// TODO CREATE A GO MODULE LIB and change to use as ...string
// Validate a struct fields by reflection
func ValidateEmpty(x interface{}, skip string) []string {

	v := reflect.ValueOf(x).Elem()

	values := make([]string, 0)

	for i := 0; i < v.NumField(); i++ {

		if v.Type().Field(i).Name != skip {
			if v.Field(i).Interface() == nil || v.Field(i).Interface() == "" {
				values = append(values, fmt.Sprintf("The %s must be not empty", v.Type().Field(i).Name))
			}
		}
	}
	return values
}

func contains(arr []string, value string) bool {

	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}
