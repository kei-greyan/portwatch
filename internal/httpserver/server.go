// Package httpserver provides a lightweight HTTP server exposing
// internal endpoints such as metrics and health.
package httpserver

import (
	"context"
	"net/http"
	"time"
)

// Server wraps an http.Server with graceful shutdown support.
type Server struct {
	srv *http.Server
}

// New creates a Server listening on addr with the provided handler.
func New(addr string, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}

// Start begins listening in a goroutine. It returns immediately.
func (s *Server) Start() error {
	go func() { _ = s.srv.ListenAndServe() }()
	return nil
}

// Shutdown gracefully stops the server, waiting up to 5 seconds.
func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
