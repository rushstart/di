package di

import "errors"

var (
	ErrDefinitionNotFound       = errors.New("DI: definition not found")
	ErrInvalidDefinitionBinding = errors.New("DI: invalid definition binding")
)
