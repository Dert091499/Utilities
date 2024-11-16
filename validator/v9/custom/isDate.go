package custom

import (
	"time"

	timeConstant "github.com/Dert091499/Utilities/common/constant/time"
	"github.com/Dert091499/Utilities/common/functions"
	v9 "gopkg.in/go-playground/validator.v9"
)

func IsDate(fl v9.FieldLevel) bool {
	checkedValue := functions.ConvertReflectValueToString(fl.Field())

	_, err := time.Parse(timeConstant.DateLayout, checkedValue)
	if err != nil {
		return false
	}

	return true
}
