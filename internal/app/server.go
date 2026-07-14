package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	wrapped   http.Server
	serverErr chan error
}

func StartServer(bind string, router http.Handler) *Server {
	server := &Server{
		wrapped: http.Server{
			Addr:    bind,
			Handler: router,
		},
		serverErr: make(chan error, 1),
	}

	go func() {
		err := server.wrapped.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.serverErr <- err
		}
	}()
}

// Wait waits for error, ctl+C, or other interrupt.
func (s *Server) Wait() error {
	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-stopCtx.Done():
		return nil
	case err := <-s.serverErr:
		return err
	}
}

func (s *Server) Close() error {
	stopCtx, stopStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopStop()

	return s.wrapped.Shutdown(stopCtx)
}
