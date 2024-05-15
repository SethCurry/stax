package squid

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"go.uber.org/zap"
)

var formDecoder = schema.NewDecoder()

type RequestContext struct {
	req *http.Request
}

func (r *RequestContext) Context() context.Context {
	return r.req.Context()
}

func (r *RequestContext) UnmarshalJSON(into interface{}) error {
	defer r.req.Body.Close()

	decoder := json.NewDecoder(r.req.Body)

	return decoder.Decode(into)
}

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

type ResponseContext struct {
	resp http.ResponseWriter
}

func (r *ResponseContext) WriteJSON(status int, data interface{}) error {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r.resp.WriteHeader(status)

	_, err = r.resp.Write(marshalled)
	return err
}

func NewContext(req *http.Request, resp http.ResponseWriter, oraDB *oracledb.Tx, logger *zap.Logger) *Context {
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

type Context struct {
	Request  *RequestContext
	Response *ResponseContext
	Logger   *zap.Logger
	DB       *oracledb.Tx
}
