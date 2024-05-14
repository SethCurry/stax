package squid

import (
	"encoding/json"
	"net/http"

	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"go.uber.org/zap"
)

type RequestContext struct {
	req *http.Request
}

func (r *RequestContext) UnmarshalJSON(into interface{}) error {
	defer r.req.Body.Close()

	decoder := json.NewDecoder(r.req.Body)

	return decoder.Decode(into)
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

func NewContext(req *http.Request, resp http.ResponseWriter, oraDB *oracledb.Client, logger *zap.Logger) *Context {
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
	DB       *oracledb.Client
}
