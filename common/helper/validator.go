package helper

import (
	"fmt"
	"reflect"
)

// Validate a struct fields by reflection
func Validate(x interface{}) []string {
	v := reflect.ValueOf(x).Elem()
	values := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == nil || v.Field(i).Interface() == "" {
			values = append(values, fmt.Sprintf("The %s must be not empty", v.Type().Field(i).Name))
		}
	}
	return values
}

// TODO: EXTRACT IT TO A MODULE
