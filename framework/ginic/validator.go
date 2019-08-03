package ginic

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
)

const (
	leftBracket  = "["
	rightBracket = "]"
)

// CustomValidationFunc ...
type CustomValidationFunc struct {
	Name string
	Func validator.Func
}

// BindCustomValidator binding custom validation funcs into validator
func BindCustomValidator(validationFuncs ...CustomValidationFunc) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); !ok {
		panic("bind custom validator failed")
	} else {
		// register validation
		v.RegisterValidation("mobile", Mobile)
		v.RegisterValidation("enum", Enum)
		v.RegisterValidation("ip", IP)
		for _, cvFunc := range validationFuncs {
			name := cvFunc.Name
			if name == "" {
				name = reflect.TypeOf(cvFunc.Func).Name()
			}
			v.RegisterValidation(name, cvFunc.Func)
			println("reg validation func success", name)
		}
	}
}

// Mobile validate mobile string
func Mobile(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if s, ok := field.Interface().(string); ok {
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
	if fieldKind != reflect.String {
		return false
	}

	if s, ok := field.Interface().(string); ok {
		leftBracketIdx := strings.Index(param, leftBracket)
		rightBracketIdx := strings.Index(param, rightBracket)

		if leftBracketIdx == -1 || rightBracketIdx == -1 {
			return false
		}

		// if strings.Contains(s, substr)
		enumArr := strings.Split(param[leftBracketIdx+1:rightBracketIdx], "/")
		for _, enumVal := range enumArr {
			if s == enumVal {
				return true
			}
		}
	}
	return false
}

// IP validator regexp ip string param
func IP(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if s, ok := field.Interface().(string); ok {
		rgx := regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
		return rgx.MatchString(s)
	}
	return false
}
