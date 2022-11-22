package httpserver

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(r *mux.Router, port string) *Server {
	httpServer := &http.Server{
		Addr:    port,
		Handler: r,
	}

	s := &Server{
		server: httpServer,
	}

	return s
}

func (s *Server) Start() {
	err := s.server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Server) Stop(ctx context.Context) {
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
