package api

import (
	"Project/utils"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
   if curency , ok := fieldLevel.Field().Interface().(string); ok{
        return utils.IsSupportedCurrency(curency)
   }
   return false
}
