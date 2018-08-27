package utils

import (
	"reflect"
)

// ConvertStructToMap convert struct value into map[string]interface{}
func ConvertStructToMap(in interface{}) (out map[string]interface{}) {
	// out = make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("StructToMap only support struct")
	}

	// to do if this is empty
	if v == reflect.Zero(reflect.TypeOf(in)) {
		panic("error: StructToMap empty param in")
	}
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		jsonTagName := typ.Field(i).Tag.Get("json")

		if jsonTagName == "" {
			fieldName := typ.Field(i).Name
			out[fieldName] = v.Field(i).Interface()
		} else {
			out[jsonTagName] = v.Field(i).Interface() // set key of map to value in struct field
		}
	}
	return
}
