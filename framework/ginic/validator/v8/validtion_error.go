package validator

// import (
// 	"fmt"

// 	validator "gopkg.in/go-playground/validator.v8"
// )

// const (
// 	// ErrLayout 错误格式
// 	ErrLayout = "%v不满足%s的%s校验条件"
// )

// // FormValidationErrors ...
// type FormValidationErrors struct {
// 	errs []error
// }

// func (f FormValidationErrors) Error() string {
// 	errMsg := ""
// 	for idx, err := range f.errs {
// 		if idx == 0 {
// 			errMsg = errMsg + err.Error()
// 			continue
// 		}
// 		errMsg = errMsg + ";" + err.Error()
// 	}
// 	return errMsg
// }

// // HdlValidationErrors ...
// // 帮助处理gin框架的表单校验异常信息
// func HdlValidationErrors(err error) error {
// 	valErrs, ok := err.(validator.ValidationErrors)
// 	if !ok {
// 		return err
// 	}

// 	fves := FormValidationErrors{}
// 	for _, ve := range valErrs {
// 		newErr := fmt.Errorf(ErrLayout, ve.Value, ve.Field, ve.Tag)
// 		fves.errs = append(fves.errs, newErr)
// 	}
// 	return fves
// }
