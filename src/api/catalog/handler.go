package catalog

import (
	"encoding/json"
	"net/http"
)

// RegisterRoutes mounts catalog HTTP endpoints on the given mux.
func RegisterRoutes(mux *http.ServeMux, repo CatalogRepository) {
	h := &handler{repo: repo}
	mux.HandleFunc("GET /api/v1/danh-muc", h.handleDanhMuc)
	mux.HandleFunc("GET /api/v1/san-pham", h.handleSanPham)
}

type handler struct {
	repo CatalogRepository
}

func (h *handler) handleDanhMuc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]DanhMuc{"danh_muc": AllDanhMuc})
}

func (h *handler) handleSanPham(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	products := h.repo.Products()
	json.NewEncoder(w).Encode(map[string][]SanPhamSo{"san_pham": products})
}

// StartServer creates the mux, registers catalog routes, and starts the HTTP server.
// This is a convenience for both production and test mains.
func StartServer(addr string, repo CatalogRepository) error {
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	return http.ListenAndServe(addr, mux)
}
