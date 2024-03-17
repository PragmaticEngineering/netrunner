package encoding

import "fmt"

// confirms that EncodingError implements the error interface
var _ error = (*EncodingError)(nil)

type sentinelError string

const (
	WriteResponse sentinelError = "write response error"
	EncodeJSON    sentinelError = "encode json error"
	DecodeJSON    sentinelError = "decode json error"
)

// EncodingError is a custom error type for encoding errors.
type EncodingError struct {
	err            error
	sentinelError  sentinelError
	additionalInfo []string
}

// Error returns the error message.
func (e *EncodingError) Error() string {
	errMsg := fmt.Sprintf("%s error: %s.", e.sentinelError, e.err.Error())

	if len(e.additionalInfo) == 0 {
		return errMsg
	}

	return fmt.Sprintf("%s AdditionalInformation: %v", errMsg, e.additionalInfo)
}

// Unwrap returns the wrapped error.
func (e *EncodingError) Unwrap() error {
	return e.err
}

// NewEncodingError returns a new EncodingError.
func newEncodingError(sentinelError sentinelError, err error, additionalInfo ...string) *EncodingError {
	return &EncodingError{
		err:            err,
		sentinelError:  sentinelError,
		additionalInfo: additionalInfo,
	}
}
