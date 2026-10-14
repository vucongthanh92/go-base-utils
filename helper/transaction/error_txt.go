package transaction

import (
	"errors"

	errHandler "github.com/vucongthanh92/go-base-utils/helper/http_error"
)

// Error carries an ErrorBuilder through go-base-utils' error-only transaction
// callback. This preserves the original API status/code while ensuring GORM
// always receives a non-nil error and rolls the transaction back.
type TxnError struct {
	Builder *errHandler.ErrorBuilder
	cause   error
}

func (e *TxnError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return "transaction aborted by payment domain error"
}

func (e *TxnError) Unwrap() error {
	return e.cause
}

// Rollback annotates a business error so the shared transaction manager rolls
// back even when the ErrorBuilder has no underlying infrastructure error.
func Rollback(builder *errHandler.ErrorBuilder) *errHandler.ErrorBuilder {
	if builder == nil {
		return nil
	}

	builder.SetLogError(&TxnError{
		Builder: builder,
		cause:   builder.LogError,
	})
	return builder
}

// BuilderFromError recovers the original domain error at the transaction
// boundary, allowing handlers to retain the intended HTTP response.
func BuilderFromError(err error) *errHandler.ErrorBuilder {
	var transactionErr *TxnError
	if errors.As(err, &transactionErr) {
		return transactionErr.Builder
	}
	return nil
}
