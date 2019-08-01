package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
	vali "gopkg.in/go-playground/validator.v9"
)

// BindDataChecker .
func BindDataChecker(db *gorm.DB) error {
	validate := binding.Validator.Engine().(*vali.Validate)
	if err := validate.RegisterValidation("reschk", resourceChecker); err != nil {
		return err
	}
	return nil
}

func resourceChecker(fld vali.FieldLevel) bool {
	id := fld.Field().Int()
	if id == 0 {
		// true: recourse id == 0
		return false
	}

	chk, ok := _checkers[fld.Param()]
	if !ok {
		// TODO: panic or return fasle ?
		panic(fld.Param() + " not registered")
	}

	if err := chk.Check(id); err != nil {
		return false
	}

	return true
}
