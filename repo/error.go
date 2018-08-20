package repo

import "errors"

// ErrUnsupportedType is returned when unsupported struct type data is passed
var ErrUnsupportedType = errors.New("repo: unsupported type")
