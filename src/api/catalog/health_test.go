package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoints(t *testing.T) {
	t.Run("GET /healthz returns 200", func(t *testing.T) {
		mux := http.NewServeMux()
		RegisterHealthRoutes(mux, nil)
		ts := httptest.NewServer(mux)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/healthz")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		var body map[string]string
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["status"] != "ok" {
			t.Errorf("expected status 'ok', got %q", body["status"])
		}
	})

	t.Run("GET /readyz with healthy check returns 200", func(t *testing.T) {
		mux := http.NewServeMux()
		RegisterHealthRoutes(mux, func() error { return nil })
		ts := httptest.NewServer(mux)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/readyz")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", res.StatusCode)
		}
		var body map[string]string
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["status"] != "ok" {
			t.Errorf("expected status 'ok', got %q", body["status"])
		}
	})

	t.Run("GET /readyz with unhealthy check returns 503", func(t *testing.T) {
		mux := http.NewServeMux()
		RegisterHealthRoutes(mux, func() error { return errors.New("db not ready") })
		ts := httptest.NewServer(mux)
		defer ts.Close()

		res, err := http.Get(ts.URL + "/readyz")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("expected 503, got %d", res.StatusCode)
		}
		var body map[string]string
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["status"] != "not_ready" {
			t.Errorf("expected status 'not_ready', got %q", body["status"])
		}
	})
}
