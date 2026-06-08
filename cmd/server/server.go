package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	server *http.Server
}

func NewServer(bind string, router *chi.Mux) Server {
	server := &http.Server{
		Addr:    bind,
		Handler: router,
	}
	return Server{server: server}
}

func (s *Server) Start() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve: %v", err)
		}
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
