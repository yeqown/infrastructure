package v9

import (
	"github.com/yeqown/infrastructure/framework/ginic/validator"
	vali "gopkg.in/go-playground/validator.v9"
)

// NewResourceCheck .
func NewResourceCheck(chks ...validator.ResourceChecker) vali.Func {
	var _checkers = make(map[string]validator.ResourceChecker)

	for _, chk := range chks {
		_checkers[chk.Tag()] = chk
	}

	return func(fl vali.FieldLevel) bool {
		return resourceCheck(_checkers, fl)
	}
}
