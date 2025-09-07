package errors

type ErrorCode string

const (
	ErrBadRequest       ErrorCode = "BAD_REQUEST"
	ErrNotFound         ErrorCode = "NOT_FOUND"
	ErrInternal         ErrorCode = "INTERNAL_ERROR"
	ErrUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrValidationFailed ErrorCode = "VALIDATION_FAILED"
)
