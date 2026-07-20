package catalog

import (
	"encoding/json"
	"net/http"
)

// RegisterHealthRoutes registers /healthz and /readyz endpoints on the given mux.
// If checkReady is nil, /readyz will always return 200 (for processes without DB).
func RegisterHealthRoutes(mux *http.ServeMux, checkReady func() error) {
	mux.HandleFunc("GET /healthz", handleHealthz)
	if checkReady != nil {
		mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
			handleReadyz(w, r, checkReady)
		})
	} else {
		mux.HandleFunc("GET /readyz", handleHealthz)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleReadyz(w http.ResponseWriter, r *http.Request, check func() error) {
	w.Header().Set("Content-Type", "application/json")
	if err := check(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "not_ready", "error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
