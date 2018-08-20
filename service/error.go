package service

import "errors"

// ErrFailedToAssert error is returned on failure of assertion
var ErrFailedToAssert = errors.New("service: failed to assert")

// NotFoundError holds the name of the resource that is not found
type NotFoundError struct {
	name string
}

func (e NotFoundError) Error() string {
	return e.name + " not found"
}

// ErrProductNotFound error is returned when a product not found
var ErrProductNotFound = NotFoundError{"product"}

type noOpLogger struct{}

func (l *noOpLogger) Print(...interface{}) {
}

func (l *noOpLogger) Println(...interface{}) {
}

func (l *noOpLogger) Printf(string, ...interface{}) {
}
