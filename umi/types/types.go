package types

import "net/http"

// Custom Handler that allow us to create a simple http function that return an error
type HandlerFunc func(Context) error

// Context contains all we need for http functions
type Context struct {
	http.ResponseWriter
	*http.Request
}
