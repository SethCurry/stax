package squid

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// HandlerFunc is a function for handling HTTP requests.
// The handler is expected to write its own response to the request,
// or to return an error (not both).
type HandlerFunc func(*Context) error

// NewErrorResponse creates a new *ErrorResponse object.
func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Err: err.Error()}
}

// ErrorResponse is an error that will be marshalled to JSON and
// returned to the client.
type ErrorResponse struct {
	Err string `json:"error"`
}

// NewServer creates a new *Server that can serve HTTP requests.
// It does not register any handlers, nor does it start serving.
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
		ctx, err := s.getContext(req, resp)
		if err != nil {
			ctx.Logger.Error("failed to get transaction", zap.Error(err))
			return
		}

		defer ctx.DB.Commit()

		err = handler(ctx)
		if err != nil {
			ctx.Logger.Error("function returned error", zap.Error(err))
			errorResponse(resp, err)
		}
	}
}

func (s *Server) Get(pattern string, handler HandlerFunc) {
	s.router.Get(pattern, s.wrapHandler(handler))
}

func (s *Server) Serve(listen string) error {
	return http.ListenAndServe(listen, s.router)
}

func errorResponse(w http.ResponseWriter, gotErr error) {
	statusCode := 500
	errResponse := NewErrorResponse(gotErr)

	if oracledb.IsNotFound(gotErr) {
		statusCode = 400
		errResponse = &ErrorResponse{
			Err: "no results found",
		}
	}

	marshalled, _ := json.Marshal(errResponse)
	w.WriteHeader(statusCode)
	w.Write(marshalled)
}
