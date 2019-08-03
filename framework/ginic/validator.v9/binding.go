package validator

import (
	vali "gopkg.in/go-playground/validator.v9"
)

// NewResourceCheck .
func NewResourceCheck(chks ...ResourceChecker) vali.Func {
	var _checkers = make(map[string]ResourceChecker)

	for _, chk := range chks {
		_checkers[chk.Tag()] = chk
	}

	return func(fl vali.FieldLevel) bool {
		return resourceCheck(_checkers, fl)
	}
}
