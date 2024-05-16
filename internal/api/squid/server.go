package squid

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type HandlerFunc func(*Context) error

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Err: err.Error()}
}

type ErrorResponse struct {
	Err string `json:"error"`
}

func NewServer(oraDB *oracledb.Client, handlerLogger *zap.Logger) *Server {
	return &Server{
		handlerLogger: handlerLogger,
		router:        chi.NewRouter(),
		db:            oraDB,
	}
}

type Server struct {
	handlerLogger *zap.Logger
	router        chi.Router
	db            *oracledb.Client
}

func (s *Server) getContext(req *http.Request, resp http.ResponseWriter) (*Context, error) {
	tx, err := s.db.Tx(req.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to open transaction: %w", err)
	}

	return NewContext(req, resp, tx, s.handlerLogger), nil
}

func (s *Server) wrapHandler(handler HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		var err error
		var ctx *Context

		defer func() {
			if err != nil {
				errResponse := NewErrorResponse(err)

				marshalled, _ := json.Marshal(errResponse)

				resp.WriteHeader(500)
				resp.Write(marshalled)
			}
		}()

		ctx, err = s.getContext(req, resp)
		if err != nil {
			ctx.Logger.Error("failed to get transaction", zap.Error(err))
			return
		}

		defer ctx.DB.Commit()

		err = handler(ctx)
		if err != nil {
			ctx.Logger.Error("function returned error", zap.Error(err))
		}
	}
}

func (s *Server) Get(pattern string, handler HandlerFunc) {
	s.router.Get(pattern, s.wrapHandler(handler))
}

func (s *Server) Serve(listen string) error {
	return http.ListenAndServe(listen, s.router)
}
