package domain

import "errors"

var (
	ErrServerRun = errors.New("failed to run http server")
)
