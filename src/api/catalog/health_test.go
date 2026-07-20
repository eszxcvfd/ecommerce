package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

	t.Run("GET /readyz with unhealthy check returns 503 without leaking raw error text", func(t *testing.T) {
		mux := http.NewServeMux()
		rawErr := "sql: /var/data/db.sqlite3: disk I/O error"
		RegisterHealthRoutes(mux, func() error { return errors.New(rawErr) })
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
		var body map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["status"] != "not_ready" {
			t.Errorf("expected status 'not_ready', got %q", body["status"])
		}
		if errStr, ok := body["error"].(string); ok {
			if errStr == rawErr {
				t.Errorf("/readyz leaked raw check error text: %q", rawErr)
			}
			if strings.Contains(errStr, "/var/data/") {
				t.Errorf("/readyz leaked filesystem path in error: %q", errStr)
			}
		}
	})

}
