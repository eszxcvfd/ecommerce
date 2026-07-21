package catalog

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ecommerce/api/account"
)

// RegisterSellerRoutes mounts seller draft endpoints on the given mux.
// It uses the account module's auth middleware to protect all routes.
func RegisterSellerRoutes(mux *http.ServeMux, repo CatalogRepository, accountRepo account.AccountRepository) {
	h := &sellerHandler{
		repo:        repo,
		accountRepo: accountRepo,
	}
	auth := account.RequireAuth(accountRepo)

	// Draft CRUD — all require auth + seller profile
	mux.HandleFunc("POST /api/v1/seller/san-pham", auth(h.requireSeller(h.handleCreateDraft)))
	mux.HandleFunc("GET /api/v1/seller/san-pham", auth(h.requireSeller(h.handleListDrafts)))
	mux.HandleFunc("GET /api/v1/seller/san-pham/{id}", auth(h.requireSeller(h.handleGetDraft)))
	mux.HandleFunc("PUT /api/v1/seller/san-pham/{id}", auth(h.requireSeller(h.handleUpdateDraft)))
	mux.HandleFunc("DELETE /api/v1/seller/san-pham/{id}", auth(h.requireSeller(h.handleDeleteDraft)))
}

type sellerHandler struct {
	repo        CatalogRepository
	accountRepo account.AccountRepository
}

// requireSeller wraps a handler to require an activated seller profile.
func (h *sellerHandler) requireSeller(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taiKhoanID := account.TaiKhoanIDFromContext(r.Context())
		if taiKhoanID == "" {
			writeError(w, http.StatusUnauthorized, "thieu_token", "Thiếu thông tin xác thực")
			return
		}

		hs, err := h.accountRepo.HOSoBanByTaiKhoanID(r.Context(), taiKhoanID)
		if err != nil {
			if errors.Is(err, account.ErrHOSoBanNotFound) {
				writeError(w, http.StatusForbidden, "khong_co_ho_so_ban", "Cần kích hoạt hồ sơ người bán")
				return
			}
			log.Printf("check seller profile error: %v", err)
			writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi kiểm tra hồ sơ người bán")
			return
		}
		if hs == nil {
			writeError(w, http.StatusForbidden, "khong_co_ho_so_ban", "Cần kích hoạt hồ sơ người bán")
			return
		}

		next(w, r)
	}
}

func (h *sellerHandler) handleCreateDraft(w http.ResponseWriter, r *http.Request) {
	taiKhoanID := account.TaiKhoanIDFromContext(r.Context())

	var input DraftInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "json_khong_hop_le", "Dữ liệu JSON không hợp lệ")
		return
	}

	if err := ValidateDraftInput(input); err != nil {
		writeError(w, http.StatusBadRequest, "du_lieu_khong_hop_le", err.Error())
		return
	}

	sp, err := h.repo.CreateDraft(input, taiKhoanID)
	if err != nil {
		log.Printf("create draft error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi tạo bản nháp")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]SanPhamSo{"san_pham": *sp})
}

func (h *sellerHandler) handleListDrafts(w http.ResponseWriter, r *http.Request) {
	taiKhoanID := account.TaiKhoanIDFromContext(r.Context())

	drafts, err := h.repo.DraftsBySeller(taiKhoanID)
	if err != nil {
		log.Printf("list drafts error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi lấy danh sách bản nháp")
		return
	}
	if drafts == nil {
		drafts = []SanPhamSo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"san_pham": drafts})
}

func (h *sellerHandler) handleGetDraft(w http.ResponseWriter, r *http.Request) {
	taiKhoanID := account.TaiKhoanIDFromContext(r.Context())
	id := r.PathValue("id")

	sp, err := h.repo.DraftByID(id, taiKhoanID)
	if err != nil {
		log.Printf("get draft error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi lấy bản nháp")
		return
	}
	if sp == nil {
		writeError(w, http.StatusNotFound, "khong_tim_thay", "Không tìm thấy bản nháp")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]SanPhamSo{"san_pham": *sp})
}

func (h *sellerHandler) handleUpdateDraft(w http.ResponseWriter, r *http.Request) {
	taiKhoanID := account.TaiKhoanIDFromContext(r.Context())
	id := r.PathValue("id")

	var input DraftUpdateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "json_khong_hop_le", "Dữ liệu JSON không hợp lệ")
		return
	}

	if err := ValidateDraftUpdateInput(input); err != nil {
		writeError(w, http.StatusBadRequest, "du_lieu_khong_hop_le", err.Error())
		return
	}

	sp, err := h.repo.UpdateDraft(id, taiKhoanID, input)
	if err != nil {
		log.Printf("update draft error: %v", err)
		writeError(w, http.StatusInternalServerError, "loi_he_thong", "Lỗi cập nhật bản nháp")
		return
	}
	if sp == nil {
		writeError(w, http.StatusNotFound, "khong_tim_thay", "Không tìm thấy bản nháp")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]SanPhamSo{"san_pham": *sp})
}

func (h *sellerHandler) handleDeleteDraft(w http.ResponseWriter, r *http.Request) {
	taiKhoanID := account.TaiKhoanIDFromContext(r.Context())
	id := r.PathValue("id")

	if err := h.repo.DeleteDraft(id, taiKhoanID); err != nil {
		log.Printf("delete draft error: %v", err)
		writeError(w, http.StatusNotFound, "khong_tim_thay", "Không tìm thấy bản nháp")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
