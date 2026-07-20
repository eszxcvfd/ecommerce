package catalog

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

// closeRecorder wraps a close callback to track shutdown ordering.
type closeRecorder struct {
	counter *atomic.Int64 // shared counter — incremented once per Close
	seq     int64         // captured sequence at Close time
	label   string
	err     error
}

func (c *closeRecorder) Close() error {
	c.seq = c.counter.Add(1)
	return c.err
}

func TestGracefulShutdown_ClosesHTTPBeforeSQLite(t *testing.T) {
	var counter atomic.Int64

	httpCloser := &closeRecorder{counter: &counter, label: "http"}
	sqliteCloser := &closeRecorder{counter: &counter, label: "sqlite"}

	server := &Server{
		httpServer: &http.Server{Addr: ":0"},
		db:         sqliteCloser,
		httpCloser: httpCloser,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	// HTTP must close before SQLite
	if httpCloser.seq != 1 || sqliteCloser.seq != 2 {
		t.Errorf("expected HTTP close (seq=1) before SQLite close (seq=2), got HTTP=%d SQLite=%d",
			httpCloser.seq, sqliteCloser.seq)
	}
}

func TestGracefulShutdown_ReturnsHTTPErrorFirst(t *testing.T) {
	var counter atomic.Int64

	httpCloser := &closeRecorder{counter: &counter, label: "http", err: errors.New("http error")}
	sqliteCloser := &closeRecorder{counter: &counter, label: "sqlite"}

	server := &Server{
		httpServer: &http.Server{Addr: ":0"},
		db:         sqliteCloser,
		httpCloser: httpCloser,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err == nil || err.Error() != "http error" {
		t.Fatalf("expected http error, got %v", err)
	}
}

func TestServer_ServesAndShutsDown(t *testing.T) {
	repo := NewMemoryRepo(SeedData())
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	RegisterHealthRoutes(mux, nil)

	srv := &http.Server{Addr: ":0", Handler: mux}

	// Start in background
	go func() {
		srv.ListenAndServe()
	}()

	// Give it time to start — polling would be better but this is acceptable for a simple smoke test
	time.Sleep(50 * time.Millisecond)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
}
