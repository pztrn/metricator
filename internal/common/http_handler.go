package common

import "context"

// HTTPHandlerFunc describes signature of HTTP requests handling function.
type HTTPHandlerFunc func(context.Context) string
