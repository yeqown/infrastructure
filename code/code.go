// Package code to define some code
package code

import (
	"reflect"
)

const (
	// CodeOk 0
	CodeOk = iota
	// CodeParamInvalid ...
	CodeParamInvalid
	// CodeSystemErr ...
	CodeSystemErr
	// CodeNoPermission ...
	CodeNoPermission
	// CodeServerTimeout ...
	CodeServerTimeout
	// CodeResourceNotFound ...
	CodeResourceNotFound
	// CodeIllegeOP ...
	CodeIllegeOP
	// ErrNoSuchCode ...
	ErrNoSuchCode = "错误码未定义"
)

var messages = map[int]string{
	CodeOk:               "成功",
	CodeParamInvalid:     "参数非法",
	CodeSystemErr:        "系统错误",
	CodeNoPermission:     "没有权限",
	CodeServerTimeout:    "服务超时",
	CodeResourceNotFound: "资源未找到",
	CodeIllegeOP:         "非法操作",
}

// CodeInfo define a CodeInfo type
type CodeInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewCodeInfo create a new *CodeInfo
func NewCodeInfo(code int, message string) *CodeInfo {
	if message == "" {
		message = GetMessage(code)
	}

	return &CodeInfo{
		Code:    code,
		Message: message,
	}
}

// GetCodeInfo get CodeInfo with specified code
func GetCodeInfo(code int) *CodeInfo {
	return &CodeInfo{
		Code:    code,
		Message: GetMessage(code),
	}
}

// GetMessage get code desc from messages
func GetMessage(code int) string {
	v, ok := messages[code]
	if !ok {
		return ErrNoSuchCode
	}
	return v
}

// FillCodeInfo ... fill a response struct will *CodeInfo
// TODO: validate v
func FillCodeInfo(v interface{}, ci *CodeInfo) interface{} {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		panic("v must be ptr")
	}
	field := reflect.ValueOf(v).Elem().
		FieldByName("CodeInfo")

	// set field
	field.FieldByName("Code").SetInt(int64(ci.Code))
	field.FieldByName("Message").SetString(ci.Message)

	return v
}
