package types

import (
	"encoding/json"
	"net/http"
)

// Custom Handler that allow us to create a simple http function that return an error
type HandlerFunc func(*Context) error

// Context contains all we need for http functions
type Context struct {
	http.ResponseWriter
	*http.Request
}

type Error struct {
	Err    error
	Status int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (c *Context) JSON(obj any) error {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.ResponseWriter).Encode(obj)
}

func (c *Context) FAIL(err error, status int) Error {
	return Error{err, status}
}
