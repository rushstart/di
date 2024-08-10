package di

import "errors"

var (
	ErrDefinitionNotFound  = errors.New("DI: definition not found")
	ErrDefinitionIsInvalid = errors.New("DI: definition is invalid")
)
