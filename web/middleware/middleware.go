package middleware

import (
	"net/http"
)

// Middleware represents http handler middleware
type Middleware func(http.Handler) http.Handler
