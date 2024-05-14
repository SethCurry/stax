package squid

import (
	"net/http"

	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type HandlerFunc func(*Context) error

func NewServer(oraDB *oracledb.Client, handlerLogger *zap.Logger) *Server {
	return &Server{
		handlerLogger: handlerLogger,
		router:        chi.NewRouter(),
	}
}

type Server struct {
	handlerLogger *zap.Logger
	router        chi.Router
	db            *oracledb.Client
}

func (s *Server) getContext(req *http.Request, resp http.ResponseWriter) *Context {
	return NewContext(req, resp, s.db, s.handlerLogger)
}

func (s *Server) wrapHandler(handler HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		ctx := s.getContext(req, resp)

		err := handler(ctx)
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
