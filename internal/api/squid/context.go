package squid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/SethCurry/stax/internal/bones"
	"go.uber.org/zap"
)

var formDecoder = schema.NewDecoder()

// RequestContext captures an HTTP request.  It's mostly a wrapper
// to provide additional functionality.
type RequestContext struct {
	req *http.Request
}

// Context returns the context.Context of the *http.Request underneat
// this RequestContext.
func (r *RequestContext) Context() context.Context {
	return r.req.Context()
}

// UnmarshalJSON unmarshals the contents of the request body into the
// provided struct.
func (r *RequestContext) UnmarshalJSON(into interface{}) error {
	defer r.req.Body.Close()

	decoder := json.NewDecoder(r.req.Body)

	return decoder.Decode(into)
}

// UnmarshalForm unmarshals any form data from an HTTP POST request into
// the provided struct.
//
// This uses github.com/gorilla/schema so see their documentation for how to
// use struct tags and what not.
func (r *RequestContext) UnmarshalForm(into interface{}) error {
	err := r.req.ParseForm()
	if err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	err = formDecoder.Decode(into, r.req.PostForm)
	if err != nil {
		return fmt.Errorf("failed to decode form data: %w", err)
	}

	return nil
}

// UnmarshalQuery unmarshals the query parameters from an HTTP request into
// the provided struct.
//
// This uses github.com/gorilla/schema so see their documentation for how to
// use the struct tags and what not.
func (r *RequestContext) UnmarshalQuery(into interface{}) error {
	err := formDecoder.Decode(into, r.req.URL.Query())
	if err != nil {
		return fmt.Errorf("failed to parse URL query: %w", err)
	}

	return nil
}

// ResponseContext captures an HTTP response.  It's mostly a wrapper
// around an http.ResponseWriter to provide additional functionality.
type ResponseContext struct {
	resp http.ResponseWriter
}

// WriteJSON writes the provided status code on the response, and then
// writes a JSON-marshalled copy of the provided data to the body of that response.
// It also sets the Content-Type to application/json.
//
// Use this any time you want to return JSONified data.
func (r *ResponseContext) WriteJSON(status int, data interface{}) error {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r.resp.Header().Set("Content-Type", "application/json")
	r.resp.WriteHeader(status)

	_, err = r.resp.Write(marshalled)
	return err
}

// NewContext initializes a new *Context object.
func NewContext(req *http.Request, resp http.ResponseWriter, oraDB *bones.Tx, logger *zap.Logger) *Context {
	return &Context{
		Request: &RequestContext{
			req: req,
		},
		Response: &ResponseContext{
			resp: resp,
		},
		Logger: logger,
		DB:     oraDB,
	}
}

// Context is a wrapper around the entirety of an HTTP request's context.
// It includes the request and response, as well as injected dependencies
// like a logger and a database transaction.
type Context struct {
	Request  *RequestContext
	Response *ResponseContext
	Logger   *zap.Logger
	DB       *bones.Tx
}
