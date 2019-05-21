package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// ConvertStructToMap convert struct value into map[string]interface{}
func ConvertStructToMap(in interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		println("StructToMap only support struct")
		return out
	}

	// to do if this is empty
	if v == reflect.Zero(reflect.TypeOf(in)) {
		println("error: StructToMap empty param in")
		return out
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
	return out
}

// ConvertUint8Slice2String ...
func ConvertUint8Slice2String(u8s []uint8) string {
	bs := make([]byte, 0)
	for _, u := range u8s {
		bs = append(bs, byte(u))
	}
	return string(bs)
}

// ToLower convert struct(only) fields to lower string
// if the field value is string type ~
func ToLower(v interface{}) {
	if !mustbePtr(v) {
		fmt.Printf("ToLower param(%s) not type of pointer\n",
			reflect.ValueOf(v).Kind().String())
		return
	}
	// must be pointer
	value := reflect.ValueOf(v).Elem()
	if !typeEqual(value, reflect.Struct) {
		fmt.Println("ToLower param not type of struct")
		return
	}

	// range and toLower
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Type().Kind() {
		case reflect.String:
			field.SetString(
				strings.ToLower(field.String()),
			)
		case reflect.Ptr:
			ToLower(field.Interface())
		}
	}
}

func mustbePtr(in interface{}) bool {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		return true
	}
	return false
}

func typeEqual(v reflect.Value, kind reflect.Kind) bool {
	return v.Type().Kind() == kind
}
