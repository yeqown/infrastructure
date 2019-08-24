package ginic

import (
	"log"
	"reflect"
	"regexp"
	"strings"

	// "github.com/gin-gonic/gin/binding"
	cusVali "github.com/yeqown/infrastructure/framework/ginic/validator"
	validator "gopkg.in/go-playground/validator.v8"
)

const (
	leftBracket  = "["
	rightBracket = "]"
)

// default inner checkers
var _checkers map[string]cusVali.ResourceChecker

func init() {
	_checkers = make(map[string]cusVali.ResourceChecker)
}

type fieldLevel struct {
	V                    *validator.Validate
	TopStruct            reflect.Value
	CurrentStructOrField reflect.Value
	Field                reflect.Value
	FieldType            reflect.Type
	FieldKind            reflect.Kind
	Param                string
}

// // CustomValidationFunc ...
// type CustomValidationFunc struct {
// 	Name string
// 	Func validator.Func
// }

// // BindCustomValidator binding custom validation funcs into validator
// func BindCustomValidator(validationFuncs ...CustomValidationFunc) {
// 	if v, ok := binding.Validator.Engine().(*validator.Validate); !ok {
// 		panic("bind custom validator failed")
// 	} else {
// 		// register validation
// 		v.RegisterValidation("mobile", Mobile)
// 		v.RegisterValidation("enum", Enum)
// 		v.RegisterValidation("ip", IP)
// 		for _, cvFunc := range validationFuncs {
// 			name := cvFunc.Name
// 			if name == "" {
// 				name = reflect.TypeOf(cvFunc.Func).Name()
// 			}
// 			v.RegisterValidation(name, cvFunc.Func)
// 			println("reg validation func success", name)
// 		}
// 	}
// }

// Mobile validate mobile string
func Mobile(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	// if s, ok := field.Interface().(string); ok {
	// 	rgx := regexp.MustCompile(`^1([3-9][0-9]|14[57]|5[^4])\d{8}$`)
	// 	return rgx.MatchString(s)
	// }
	// return false
	return mobile(&fieldLevel{
		V:                    v,
		TopStruct:            topStruct,
		CurrentStructOrField: currentStructOrField,
		Field:                field,
		FieldType:            fieldType,
		FieldKind:            fieldKind,
		Param:                param,
	})
}

func mobile(fl *fieldLevel) bool {
	if s, ok := fl.Field.Interface().(string); ok {
		rgx := regexp.MustCompile(`^1([3-9][0-9]|14[57]|5[^4])\d{8}$`)
		return rgx.MatchString(s)
	}
	return false
}

// Enum validate the val is in enum type or not, only support string type
// use like this: `form:"user_type" binding:"required,enum=[01/02/03]"`
func Enum(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	return enum(&fieldLevel{
		V:                    v,
		TopStruct:            topStruct,
		CurrentStructOrField: currentStructOrField,
		Field:                field,
		FieldType:            fieldType,
		FieldKind:            fieldKind,
		Param:                param,
	})
}

func enum(fl *fieldLevel) bool {
	if fl.FieldKind != reflect.String {
		return false
	}

	s := fl.Field.String()
	leftBracketIdx := strings.Index(fl.Param, leftBracket)
	rightBracketIdx := strings.Index(fl.Param, rightBracket)

	if leftBracketIdx == -1 || rightBracketIdx == -1 {
		return false
	}

	// if strings.Contains(s, substr)
	enumArr := strings.Split(fl.Param[leftBracketIdx+1:rightBracketIdx], "/")
	for _, enumVal := range enumArr {
		if s == enumVal {
			return true
		}
	}
	return false
}

// IP validator regexp ip string param
func IP(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {

	return ip(&fieldLevel{
		V:                    v,
		TopStruct:            topStruct,
		CurrentStructOrField: currentStructOrField,
		Field:                field,
		FieldType:            fieldType,
		FieldKind:            fieldKind,
		Param:                param,
	})
}

func ip(fl *fieldLevel) bool {
	if fl.FieldKind != reflect.String {
		return false
	}

	s := fl.Field.String()
	rgx := regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
	return rgx.MatchString(s)
}

func resourceCheck(checkers map[string]cusVali.ResourceChecker, fl *fieldLevel) bool {
	chk, ok := checkers[fl.Param]
	if !ok {
		panic(fl.Param + " not registered")
	}

	switch k := fl.FieldKind; k {
	case reflect.String:
		id := fl.Field.String()
		if err := chk.Check(id); err != nil {
			log.Printf("check(%s) error: %v", id, err)
			return false
		}
		return true

	case reflect.Int64, reflect.Int:

		id := fl.Field.Int()
		if err := chk.CheckInt64(id); err != nil {
			log.Printf("check(%d) error: %v", id, err)
			return false
		}
		return true
	case reflect.Uint64, reflect.Uint:
		id := fl.Field.Uint()
		if err := chk.CheckInt64(int64(id)); err != nil {
			log.Printf("check(%d) error: %v", id, err)
			return false
		}
		return true
	}

	return false
}

// RegisterResChk to bind name with checker
func RegisterResChk(name string, ic cusVali.ResourceChecker) {
	_checkers[name] = ic
}

// DefaultResourceCheck to check resource id in request form
func DefaultResourceCheck(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	return resourceCheck(_checkers, &fieldLevel{
		V:                    v,
		TopStruct:            topStruct,
		CurrentStructOrField: currentStructOrField,
		Field:                field,
		FieldType:            fieldType,
		FieldKind:            fieldKind,
		Param:                param,
	})
}
