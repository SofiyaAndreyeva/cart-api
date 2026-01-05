package domain

import "errors"

var (
	ErrServerRun        = errors.New("failed to run http server")
	ErrLimitCart        = errors.New("cart limit reached")
	ErrCartNotFound     = errors.New("cart not found")
	ErrEmptyProduct     = errors.New("empty product")
	ErrInvalidPrice     = errors.New("invalid price")
	ErrCartItemNotFound = errors.New("cart item not found")
)
