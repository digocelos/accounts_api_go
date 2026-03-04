package httpapi

import (
	"fmt"

	"github.com/digocelo/account-api/internal/account"
)

func validationf(format string, args ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{account.ErrValidation}, args...)...)
}
