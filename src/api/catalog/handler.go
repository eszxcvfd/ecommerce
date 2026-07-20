package catalog

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
)

// RegisterRoutes mounts catalog HTTP endpoints on the given mux.
func RegisterRoutes(mux *http.ServeMux, repo CatalogRepository) {
	h := &handler{repo: repo}
	mux.HandleFunc("GET /api/v1/danh-muc", h.handleDanhMuc)
	mux.HandleFunc("GET /api/v1/san-pham", h.handleSanPham)
	mux.HandleFunc("GET /api/v1/dinh-dang", h.handleDinhDang)
}

type handler struct {
	repo CatalogRepository
}

func (h *handler) handleDanhMuc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]DanhMuc{"danh_muc": AllDanhMuc})
}

func (h *handler) handleSanPham(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	query := CatalogQuery{
		Q:        q.Get("q"),
		DanhMuc:  q.Get("danh_muc"),
		DinhDang: q.Get("dinh_dang"),
	}

	// Validate danh_muc enum
	if query.DanhMuc != "" {
		valid := false
		for _, dm := range AllDanhMuc {
			if string(dm) == query.DanhMuc {
				valid = true
				break
			}
		}
		if !valid {
			writeError(w, http.StatusBadRequest, "invalid_filter", "Danh mục không hợp lệ: "+query.DanhMuc)
			return
		}
	}

	// Validate dinh_dang against known formats from approved products
	if query.DinhDang != "" {
		formats := deriveFormats(h.repo.Products())
		valid := false
		for _, f := range formats {
			if f == query.DinhDang {
				valid = true
				break
			}
		}
		if !valid {
			writeError(w, http.StatusBadRequest, "invalid_filter", "Định dạng không hợp lệ: "+query.DinhDang)
			return
		}
	}

	// Validate sort enum
	sortStr := q.Get("sort")
	if sortStr != "" {
		valid := false
		for _, s := range ValidSortOrders {
			if string(s) == sortStr {
				valid = true
				break
			}
		}
		if !valid {
			writeError(w, http.StatusBadRequest, "invalid_filter", "Sắp xếp không hợp lệ: "+sortStr)
			return
		}
	}
	query.Sort = sortStr

	// Parse and validate price range
	if minStr := q.Get("min_xu"); minStr != "" {
		minVal, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil || minVal < 0 {
			writeError(w, http.StatusBadRequest, "invalid_filter", "min_xu không hợp lệ")
			return
		}
		query.MinXu = &minVal
	}
	if maxStr := q.Get("max_xu"); maxStr != "" {
		maxVal, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil || maxVal < 0 {
			writeError(w, http.StatusBadRequest, "invalid_filter", "max_xu không hợp lệ")
			return
		}
		query.MaxXu = &maxVal
	}

	// Validate min_xu <= max_xu when both are present
	if query.MinXu != nil && query.MaxXu != nil && *query.MinXu > *query.MaxXu {
		writeError(w, http.StatusBadRequest, "invalid_filter", "min_xu không được lớn hơn max_xu")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	products := h.repo.Search(query)
	json.NewEncoder(w).Encode(map[string][]SanPhamSo{"san_pham": products})
}

func (h *handler) handleDinhDang(w http.ResponseWriter, r *http.Request) {
	products := h.repo.Products()
	formats := deriveFormats(products)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"dinh_dang": formats})
}

// deriveFormats extracts distinct lowercase extensions from products, sorted.
func deriveFormats(products []SanPhamSo) []string {
	seen := make(map[string]bool)
	for _, p := range products {
		for _, ext := range p.DinhDang {
			seen[ext] = true
		}
	}
	var result []string
	for ext := range seen {
		result = append(result, ext)
	}
	sort.Strings(result)
	return result
}

// StartServer creates the mux, registers catalog routes, and starts the HTTP server.
// This is a convenience for both production and test mains.
func StartServer(addr string, repo CatalogRepository) error {
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	return http.ListenAndServe(addr, mux)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": code, "message": message})
}
