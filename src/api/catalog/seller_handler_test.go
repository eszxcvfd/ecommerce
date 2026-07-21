package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce/api/account"
)

// setupSellerTestServer creates a test server with seller draft routes,
// seeded account repo (with a seller user and session), and empty catalog.
// Returns the server, the account repo (for creating additional test users),
// the seller's account, and the session token.
func setupSellerTestServer(t *testing.T) (*httptest.Server, account.AccountRepository, *account.TaiKhoan, string) {
	t.Helper()
	accRepo := account.NewMemoryRepo()

	acc, err := accRepo.CreateTaiKhoan(context.Background(), "seller@test.com", "pass123", "Người Bán")
	if err != nil {
		t.Fatalf("create account: %v", err)
	}

	_, err = accRepo.TaoHOSoBan(context.Background(), acc.ID)
	if err != nil {
		t.Fatalf("create seller profile: %v", err)
	}

	session, err := accRepo.TaoPhienDangNhap(context.Background(), acc.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	catRepo := NewMemoryRepo(nil) // empty catalog

	mux := http.NewServeMux()
	RegisterSellerRoutes(mux, catRepo, accRepo)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	return ts, accRepo, acc, session.Token
}

// setupSellerTestServerWithNonSeller creates a test server with a user that has no seller profile.
func setupSellerTestServerWithNonSeller(t *testing.T) (*httptest.Server, string) {
	t.Helper()
	accRepo := account.NewMemoryRepo()

	acc, err := accRepo.CreateTaiKhoan(context.Background(), "buyer@test.com", "pass123", "Người Mua")
	if err != nil {
		t.Fatalf("create account: %v", err)
	}

	session, err := accRepo.TaoPhienDangNhap(context.Background(), acc.ID)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	catRepo := NewMemoryRepo(nil)

	mux := http.NewServeMux()
	RegisterSellerRoutes(mux, catRepo, accRepo)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	return ts, session.Token
}

func validDraftInput() DraftInput {
	return DraftInput{
		Ten:         "Bản vẽ nhà phố 3 tầng",
		MoTa:        "Bản vẽ kiến trúc nhà phố đầy đủ chi tiết",
		MoTaChiTiet: "Bản vẽ bao gồm mặt bằng, mặt đứng, mặt cắt và chi tiết cấu tạo.",
		AnhDemo:     "https://example.com/demo.jpg",
		MienPhi:     false,
		SoXu:        5000,
		DanhMuc:     DanhMucKienTruc,
		GiayPhep:    "Giấy phép tiêu chuẩn",
		Tep: []TepInput{
			{TenTep: "mat-bang.dwg", DinhDang: "dwg", DungLuongBytes: 2048000},
			{TenTep: "3d-model.skp", DinhDang: "skp", DungLuongBytes: 5120000},
		},
	}
}

func TestSellerDraft_Create_Success(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}

	var resp map[string]SanPhamSo
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	sp := resp["san_pham"]
	if sp.Ten != "Bản vẽ nhà phố 3 tầng" {
		t.Errorf("expected ten 'Bản vẽ nhà phố 3 tầng', got %q", sp.Ten)
	}
	if sp.ID == "" {
		t.Error("expected non-empty ID")
	}
	if len(sp.Tep) != 2 {
		t.Errorf("expected 2 files, got %d", len(sp.Tep))
	}
}

func TestSellerDraft_Create_RequiresAuth(t *testing.T) {
	ts, _, _, _ := setupSellerTestServer(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Create_RequiresSellerProfile(t *testing.T) {
	ts, token := setupSellerTestServerWithNonSeller(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Create_InvalidFormat(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	input := validDraftInput()
	input.Tep = []TepInput{
		{TenTep: "file.exe", DinhDang: "exe", DungLuongBytes: 1000},
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid format, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Create_EmptyName(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	input := validDraftInput()
	input.Ten = ""

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty name, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Create_NoFiles(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	input := validDraftInput()
	input.Tep = []TepInput{}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for no files, got %d", res.StatusCode)
	}
}

func TestSellerDraft_List_Success(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	// Create a draft first
	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	res.Body.Close()

	// List drafts
	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/seller/san-pham", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("list request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var resp map[string][]SanPhamSo
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp["san_pham"]) != 1 {
		t.Errorf("expected 1 draft, got %d", len(resp["san_pham"]))
	}
}

func TestSellerDraft_List_RequiresAuth(t *testing.T) {
	ts, _, _, _ := setupSellerTestServer(t)

	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/seller/san-pham", nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Get_Success(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}

	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID
	if draftID == "" {
		t.Fatal("draft ID is empty")
	}

	// Get draft by ID
	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/seller/san-pham/"+draftID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var getResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&getResp)
	sp := getResp["san_pham"]
	if sp.ID != draftID {
		t.Errorf("expected id %s, got %s", draftID, sp.ID)
	}
	if sp.Ten != "Bản vẽ nhà phố 3 tầng" {
		t.Errorf("expected ten, got %q", sp.Ten)
	}
}

func TestSellerDraft_Get_NotFound(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/seller/san-pham/nonexistent", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Ownership_Enforced(t *testing.T) {
	ts, accRepo, _, token1 := setupSellerTestServer(t)

	// Create a draft as seller 1
	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token1)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}

	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Create second seller through the SAME account repo (so server recognizes the token)
	acc2, err := accRepo.CreateTaiKhoan(context.Background(), "seller2@test.com", "pass123", "Người Bán 2")
	if err != nil {
		t.Fatalf("create account 2: %v", err)
	}
	_, err = accRepo.TaoHOSoBan(context.Background(), acc2.ID)
	if err != nil {
		t.Fatalf("create seller profile 2: %v", err)
	}
	session2, err := accRepo.TaoPhienDangNhap(context.Background(), acc2.ID)
	if err != nil {
		t.Fatalf("create session 2: %v", err)
	}

	// Seller 2 tries to read seller 1's draft
	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/seller/san-pham/"+draftID, nil)
	req.Header.Set("Authorization", "Bearer "+session2.Token)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 (ownership enforced), got %d", res.StatusCode)
	}
}

func TestSellerDraft_Update_Success(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	// Create a draft
	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}

	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Update: change name and price
	updateInput := map[string]interface{}{
		"ten":      "Bản vẽ nhà phố 5 tầng (đã sửa)",
		"so_xu":    8000,
		"mien_phi": false,
	}
	updateBody, _ := json.Marshal(updateInput)
	req, _ = http.NewRequest("PUT", ts.URL+"/api/v1/seller/san-pham/"+draftID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("update request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var updateResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&updateResp)
	sp := updateResp["san_pham"]
	if sp.Ten != "Bản vẽ nhà phố 5 tầng (đã sửa)" {
		t.Errorf("expected updated name, got %q", sp.Ten)
	}
	if sp.Gia.SoXu != 8000 {
		t.Errorf("expected so_xu 8000, got %d", sp.Gia.SoXu)
	}
}

func TestSellerDraft_Update_Ownership(t *testing.T) {
	ts, accRepo, _, token1 := setupSellerTestServer(t)

	// Create a draft as seller 1
	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token1)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}

	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Create second seller through the SAME account repo
	acc2, err := accRepo.CreateTaiKhoan(context.Background(), "seller2@test.com", "pass123", "Người Bán 2")
	if err != nil {
		t.Fatalf("create account 2: %v", err)
	}
	_, err = accRepo.TaoHOSoBan(context.Background(), acc2.ID)
	if err != nil {
		t.Fatalf("create seller profile 2: %v", err)
	}
	session2, err := accRepo.TaoPhienDangNhap(context.Background(), acc2.ID)
	if err != nil {
		t.Fatalf("create session 2: %v", err)
	}

	// Seller 2 tries to update seller 1's draft
	updateBody, _ := json.Marshal(map[string]string{"ten": "Hacked"})
	req, _ = http.NewRequest("PUT", ts.URL+"/api/v1/seller/san-pham/"+draftID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+session2.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 (ownership enforced), got %d", res.StatusCode)
	}
}

// ---------------------------------------------------------------------------
// Regression: update input validation
// ---------------------------------------------------------------------------

func TestSellerDraft_Update_EmptyName(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}
	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Update with empty name — should be rejected
	updateBody, _ := json.Marshal(map[string]string{"ten": ""})
	req, _ = http.NewRequest("PUT", ts.URL+"/api/v1/seller/san-pham/"+draftID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty name, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Update_InvalidCategory(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}
	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Update with invalid category
	cat := "không_tồn_tại"
	updateBody, _ := json.Marshal(map[string]string{"danh_muc": cat})
	req, _ = http.NewRequest("PUT", ts.URL+"/api/v1/seller/san-pham/"+draftID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid category, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Update_NegativePrice(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}
	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Update with negative price
	updateBody, _ := json.Marshal(map[string]int64{"so_xu": -100})
	req, _ = http.NewRequest("PUT", ts.URL+"/api/v1/seller/san-pham/"+draftID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for negative price, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Update_EmptyFileList(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}
	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Update with explicit empty file list — should be rejected
	updateBody, _ := json.Marshal(map[string]interface{}{"tep": []interface{}{}})
	req, _ = http.NewRequest("PUT", ts.URL+"/api/v1/seller/san-pham/"+draftID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty file list, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Update_Partial_NoChanges(t *testing.T) {
	// Sending just {} should succeed (valid no-op update)
	ts, _, _, token := setupSellerTestServer(t)

	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}
	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Empty update (no fields) — should succeed
	updateBody, _ := json.Marshal(map[string]interface{}{})
	req, _ = http.NewRequest("PUT", ts.URL+"/api/v1/seller/san-pham/"+draftID, bytes.NewReader(updateBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for empty update body, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Delete_Success(t *testing.T) {
	ts, _, _, token := setupSellerTestServer(t)

	// Create a draft
	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}

	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Delete
	req, _ = http.NewRequest("DELETE", ts.URL+"/api/v1/seller/san-pham/"+draftID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	// Verify deleted: should return 404
	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/seller/san-pham/"+draftID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("get after delete: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 after delete, got %d", res.StatusCode)
	}
}

func TestSellerDraft_Delete_Ownership(t *testing.T) {
	ts, accRepo, _, token1 := setupSellerTestServer(t)

	// Create a draft as seller 1
	body, _ := json.Marshal(validDraftInput())
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/seller/san-pham", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token1)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d", res.StatusCode)
	}

	var createResp map[string]SanPhamSo
	json.NewDecoder(res.Body).Decode(&createResp)
	draftID := createResp["san_pham"].ID

	// Create second seller through the SAME account repo
	acc2, err := accRepo.CreateTaiKhoan(context.Background(), "seller2@test.com", "pass123", "Người Bán 2")
	if err != nil {
		t.Fatalf("create account 2: %v", err)
	}
	_, err = accRepo.TaoHOSoBan(context.Background(), acc2.ID)
	if err != nil {
		t.Fatalf("create seller profile 2: %v", err)
	}
	session2, err := accRepo.TaoPhienDangNhap(context.Background(), acc2.ID)
	if err != nil {
		t.Fatalf("create session 2: %v", err)
	}

	// Seller 2 tries to delete seller 1's draft
	req, _ = http.NewRequest("DELETE", ts.URL+"/api/v1/seller/san-pham/"+draftID, nil)
	req.Header.Set("Authorization", "Bearer "+session2.Token)

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 (ownership enforced), got %d", res.StatusCode)
	}
}

func TestSellerDraft_ExcludedFromPublicCatalog(t *testing.T) {
	// Create a server with BOTH public catalog routes AND seller routes
	accRepo := account.NewMemoryRepo()
	catRepo := NewMemoryRepo(SeedData()) // seed with public approved products

	// Add a seller user
	acc, _ := accRepo.CreateTaiKhoan(context.Background(), "seller@test.com", "pass123", "Người Bán")
	accRepo.TaoHOSoBan(context.Background(), acc.ID)
	session, _ := accRepo.TaoPhienDangNhap(context.Background(), acc.ID)

	// Create a draft in the catalog repo directly
	draftInput := validDraftInput()
	created, err := catRepo.CreateDraft(draftInput, acc.ID)
	if err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if created == nil {
		t.Fatal("draft creation returned nil")
	}

	// Build mux with both public and seller routes
	mux := http.NewServeMux()
	RegisterRoutes(mux, catRepo)
	RegisterSellerRoutes(mux, catRepo, accRepo)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Public catalog should NOT include the draft
	res, err := http.Get(ts.URL + "/api/v1/san-pham")
	if err != nil {
		t.Fatalf("catalog request: %v", err)
	}
	defer res.Body.Close()

	var catalogResp map[string][]SanPhamSo
	json.NewDecoder(res.Body).Decode(&catalogResp)
	for _, p := range catalogResp["san_pham"] {
		if p.ID == created.ID {
			t.Errorf("public catalog contains draft product %s", p.ID)
		}
	}

	// Public detail should also NOT return draft
	res, err = http.Get(ts.URL + "/api/v1/san-pham/" + created.ID)
	if err != nil {
		t.Fatalf("detail request: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404 for draft in public detail, got %d", res.StatusCode)
	}

	// But seller can access their draft through seller endpoint
	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/seller/san-pham/"+created.ID, nil)
	req.Header.Set("Authorization", "Bearer "+session.Token)
	res2, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("seller get draft: %v", err)
	}
	defer res2.Body.Close()
	if res2.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for seller getting own draft, got %d", res2.StatusCode)
	}
}
