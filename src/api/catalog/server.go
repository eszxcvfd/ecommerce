package catalog

import (
	"context"
	"database/sql"
	"io"
	"net/http"
)

// Server wraps an HTTP server with a SQLite database connection for graceful shutdown.
// On Shutdown, it drains the HTTP server first, then closes the database connection.
type Server struct {
	httpServer *http.Server
	db         io.Closer
	httpCloser io.Closer // for testing: wraps httpServer.Shutdown
}

// NewServer creates a Server with catalog routes and health endpoints registered.
// The db parameter is optional (nil is allowed for testing without SQLite).
func NewServer(addr string, repo CatalogRepository, db *sql.DB, checkReady func() error) *Server {
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	RegisterHealthRoutes(mux, checkReady)


	return &Server{
		httpServer: &http.Server{Addr: addr, Handler: mux},
		db:         db,
	}
}

// NewServerWithHandler creates a Server with the given pre-built HTTP handler (mux).
// This is used when the caller wants to register routes from multiple modules.
// The db parameter is optional (nil is allowed for testing without SQLite).
func NewServerWithHandler(addr string, handler http.Handler, db *sql.DB) *Server {
	return &Server{
		httpServer: &http.Server{Addr: addr, Handler: handler},
		db:         db,
	}
}

// ListenAndServe starts the HTTP server and blocks until it fails.
func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server: first drains HTTP connections,
// then closes the database connection. It respects the given context for
// the HTTP drain timeout.
func (s *Server) Shutdown(ctx context.Context) error {
	// Close HTTP first — drain in-flight requests
	var httpErr error
	if s.httpCloser != nil {
		httpErr = s.httpCloser.Close()
	} else {
		httpErr = s.httpServer.Shutdown(ctx)
	}

	// Then close SQLite
	if s.db != nil {
		if err := s.db.Close(); err != nil && httpErr == nil {
			httpErr = err
		}
	}

	return httpErr
}
