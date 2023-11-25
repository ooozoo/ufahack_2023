package valid

import (
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

// Singleton
var v *validator.Validate

func GetValidator() *validator.Validate {
	if v == nil {
		sync.OnceFunc(func() {
			v = validator.New()
			v.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
				if name == "-" {
					return ""
				}
				return name
			})
		})()
	}

	return v
}
