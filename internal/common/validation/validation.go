package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var msg strings.Builder
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()

			msg.WriteString(fmt.Sprintf("Поле '%s' не прошло проверку: %s\n", field, tag))
		}
		return msg.String()
	}
	return err.Error()
}
