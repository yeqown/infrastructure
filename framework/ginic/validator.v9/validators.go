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

// default inner checkers
var _checkers map[string]ResourceChecker

func init() {
	_checkers = make(map[string]ResourceChecker)
}

// Enum validate the val is in enum type or not, only support string type
// use like this: `form:"user_type" binding:"required,enum=[01/02/03]"`
func Enum(fl vali.FieldLevel) bool {
	if s := fl.Field().String(); s != "" {
		param := fl.Param()
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
func Mobile(fl vali.FieldLevel) bool {
	if s := fl.Field().String(); s != "" {
		return rgxMobile.MatchString(s)
	}
	return false
}

// IP validator regexp ip string param
func IP(fl vali.FieldLevel) bool {
	if s := fl.Field().String(); s != "" {
		return rgxIP.MatchString(s)
	}
	return false
}

func resourceCheck(checkers map[string]ResourceChecker, fl vali.FieldLevel) bool {
	chk, ok := checkers[fl.Param()]
	if !ok {
		panic(fl.Param() + " not registered")
	}

	switch k := fl.Field().Kind(); k {
	case reflect.String:
		id := fl.Field().String()
		if err := chk.Check(id); err != nil {
			log.Printf("check(%s) error: %v", id, err)
			return false
		}
		return true

	case reflect.Int64, reflect.Int:

		id := fl.Field().Int()
		if err := chk.CheckInt64(id); err != nil {
			log.Printf("check(%d) error: %v", id, err)
			return false
		}
		return true
	case reflect.Uint64, reflect.Uint:
		id := fl.Field().Uint()
		if err := chk.CheckInt64(int64(id)); err != nil {
			log.Printf("check(%d) error: %v", id, err)
			return false
		}
		return true
	}

	return false
}

// DefaultResourceCheck to check resource id in request form
func DefaultResourceCheck(fl vali.FieldLevel) bool {
	return resourceCheck(_checkers, fl)
}
