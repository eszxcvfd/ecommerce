package account

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

// openTestSQLite opens a temporary SQLite database for testing.
func openTestSQLite(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// setupSQLiteTestServer creates a test server backed by a SQLite database.
func setupSQLiteTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	db := openTestSQLite(t)

	// Run migrations
	if err := RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	repo := NewSQLiteRepo(db)
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func TestSQLiteDangKy_Success(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	body := `{"email":"user@sqlite.com","password":"secret123","ten":"SQLite User"}`
	res, err := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}

	var resp struct {
		TaiKhoan TaiKhoanPublic `json:"tai_khoan"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if resp.TaiKhoan.Email != "user@sqlite.com" {
		t.Errorf("expected user@sqlite.com, got %s", resp.TaiKhoan.Email)
	}
	if resp.TaiKhoan.VaiTro != VaiTroNguoiMua {
		t.Errorf("expected %s, got %s", VaiTroNguoiMua, resp.TaiKhoan.VaiTro)
	}
}

func TestSQLiteDangKy_Duplicate(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	body := `{"email":"dup@sqlite.com","password":"secret123","ten":"Dup"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(body))
	res.Body.Close()

	res, err := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", res.StatusCode)
	}
}

func TestSQLiteDangNhap_Success(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	// Register
	regBody := `{"email":"login@sqlite.com","password":"secret123","ten":"SQL Login"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	// Login
	loginBody := `{"email":"login@sqlite.com","password":"secret123"}`
	res, err := http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(loginBody))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var resp DangNhapResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode login: %v", err)
	}

	if resp.TaiKhoan.Email != "login@sqlite.com" {
		t.Errorf("expected login@sqlite.com, got %s", resp.TaiKhoan.Email)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}

	// Use token to access protected endpoint
	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/tai-khoan/me", nil)
	req.Header.Set("Authorization", "Bearer "+resp.Token)
	meRes, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer meRes.Body.Close()

	if meRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from /me, got %d", meRes.StatusCode)
	}
}

func TestSQLiteDangXuat_InvalidatesToken(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	// Register and login
	regBody := `{"email":"logout@sqlite.com","password":"secret123","ten":"Logout"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	loginBody := `{"email":"logout@sqlite.com","password":"secret123"}`
	res, _ = http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(loginBody))
	var loginResp DangNhapResponse
	json.NewDecoder(res.Body).Decode(&loginResp)
	res.Body.Close()

	// Logout
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/dang-xuat", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	// Token should be invalid
	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/tai-khoan/me", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 after logout, got %d", res.StatusCode)
	}
}

func TestSQLiteHOSoBan_FullCycle(t *testing.T) {
	ts := setupSQLiteTestServer(t)

	// Register and login
	regBody := `{"email":"seller@sqlite.com","password":"secret123","ten":"Seller"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	loginBody := `{"email":"seller@sqlite.com","password":"secret123"}`
	res, _ = http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(loginBody))
	var loginResp DangNhapResponse
	json.NewDecoder(res.Body).Decode(&loginResp)
	res.Body.Close()

	// Activate seller profile
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/ho-so-nguoi-ban", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}

	// Fetch seller profile
	req, _ = http.NewRequest("GET", ts.URL+"/api/v1/ho-so-nguoi-ban", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var hs HOSoBan
	if err := json.NewDecoder(res.Body).Decode(&hs); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if hs.TrangThai != TrangThaiHOSoBanKichHoat {
		t.Errorf("expected %s, got %s", TrangThaiHOSoBanKichHoat, hs.TrangThai)
	}

	// Duplicate activation should fail
	req, _ = http.NewRequest("POST", ts.URL+"/api/v1/ho-so-nguoi-ban", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate, got %d", res.StatusCode)
	}
}
