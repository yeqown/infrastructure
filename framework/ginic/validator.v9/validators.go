package validator

import (
	"log"
	"reflect"
	"regexp"
	"strings"

	vali "gopkg.in/go-playground/validator.v9"
)

const (
	leftBracket  = "["
	rightBracket = "]"
)

// Enum validate the val is in enum type or not, only support string type
// use like this: `form:"user_type" binding:"required,enum=[01/02/03]"`
func Enum(fld vali.FieldLevel) bool {
	if s := fld.Field().String(); s != "" {
		param := fld.Param()
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

var (
	rgxMobile = regexp.MustCompile(`^1([3-9][0-9]|14[57]|5[^4])\d{8}$`)
	rgxIP     = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
)

// Mobile validate mobile string
func Mobile(fld vali.FieldLevel) bool {
	if s := fld.Field().String(); s != "" {
		return rgxMobile.MatchString(s)
	}
	return false
}

// IP validator regexp ip string param
func IP(fld vali.FieldLevel) bool {
	if s := fld.Field().String(); s != "" {
		return rgxIP.MatchString(s)
	}
	return false
}

// ResourceCheck to check resource id in request form
func ResourceCheck(fld vali.FieldLevel) bool {
	chk, ok := _checkers[fld.Param()]
	if !ok {
		panic(fld.Param() + " not registered")
	}

	switch k := fld.Field().Kind(); k {
	case reflect.String:
		id := fld.Field().String()
		if err := chk.Check(id); err != nil {
			log.Printf("check(%s) error: %v", id, err)
			return false
		}
		return true

	case reflect.Int64, reflect.Int:

		id := fld.Field().Int()
		if err := chk.CheckInt64(id); err != nil {
			log.Printf("check(%d) error: %v", id, err)
			return false
		}
		return true
	case reflect.Uint64, reflect.Uint:
		id := fld.Field().Uint()
		if err := chk.CheckInt64(int64(id)); err != nil {
			log.Printf("check(%d) error: %v", id, err)
			return false
		}
		return true
	}

	return false
}
