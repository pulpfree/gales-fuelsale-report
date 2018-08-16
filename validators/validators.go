package validators

import (
	"reflect"
	"time"

	validator "gopkg.in/go-playground/validator.v8"
)

// ValidateDate function
// Ensures that the requested date is not in the future
func ValidateDate(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if date, ok := field.Interface().(time.Time); ok {
		today := time.Now()
		if today.Unix() < date.Unix() {
			return false
		}
	}
	return true
}
