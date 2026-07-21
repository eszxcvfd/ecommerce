package account

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// setupTestServer creates a test server with account routes and in-memory storage.
func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	repo := NewMemoryRepo()
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// setupTestServerWithAdmin creates a test server seeded with admin credentials.
func setupTestServerWithAdmin(t *testing.T) *httptest.Server {
	t.Helper()
	repo := NewMemoryRepo()
	WithAdminAccount(repo, "admin@test.com", "admin123", "Quản trị viên")
	mux := http.NewServeMux()
	RegisterRoutes(mux, repo)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func TestDangKy_Success(t *testing.T) {
	ts := setupTestServer(t)

	body := `{"email":"user@test.com","password":"secret123","ten":"Nguyễn Văn A"}`
	res, err := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", res.StatusCode)
	}

	var resp struct {
		TaiKhoan TaiKhoanPublic `json:"tai_khoan"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.TaiKhoan.Email != "user@test.com" {
		t.Errorf("expected email user@test.com, got %s", resp.TaiKhoan.Email)
	}
	if resp.TaiKhoan.Ten != "Nguyễn Văn A" {
		t.Errorf("expected ten Nguyễn Văn A, got %s", resp.TaiKhoan.Ten)
	}
	if resp.TaiKhoan.VaiTro != VaiTroNguoiMua {
		t.Errorf("expected default role %s, got %s", VaiTroNguoiMua, resp.TaiKhoan.VaiTro)
	}
	if resp.TaiKhoan.ID == "" {
		t.Error("expected non-empty account ID")
	}
}

func TestDangKy_DuplicateEmail(t *testing.T) {
	ts := setupTestServer(t)

	body := `{"email":"dup@test.com","password":"secret123","ten":"User"}`
	res, err := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	// Second registration with same email
	res, err = http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 Conflict for duplicate email, got %d", res.StatusCode)
	}
}

func TestDangKy_InvalidBody(t *testing.T) {
	ts := setupTestServer(t)

	tests := []struct {
		name string
		body string
	}{
		{"missing email", `{"password":"secret123","ten":"User"}`},
		{"missing password", `{"email":"user@test.com","ten":"User"}`},
		{"empty email", `{"email":"","password":"secret123","ten":"User"}`},
		{"empty password", `{"email":"user@test.com","password":"","ten":"User"}`},
		{"malformed JSON", `not-json`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", res.StatusCode)
			}
		})
	}
}

func TestDangNhap_Success(t *testing.T) {
	// Register first
	ts := setupTestServer(t)
	regBody := `{"email":"login@test.com","password":"secret123","ten":"User"}`
	res, err := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	// Login
	loginBody := `{"email":"login@test.com","password":"secret123"}`
	res, err = http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(loginBody))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.StatusCode)
	}

	var resp DangNhapResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.TaiKhoan.Email != "login@test.com" {
		t.Errorf("expected login@test.com, got %s", resp.TaiKhoan.Email)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
}

func TestDangNhap_InvalidCredentials(t *testing.T) {
	ts := setupTestServer(t)

	// Register
	regBody := `{"email":"auth@test.com","password":"correct","ten":"User"}`
	res, err := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	tests := []struct {
		name string
		body string
	}{
		{"wrong password", `{"email":"auth@test.com","password":"wrong"}`},
		{"unknown email", `{"email":"nobody@test.com","password":"correct"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected 401, got %d", res.StatusCode)
			}
		})
	}
}

func TestTaiKhoanMe_RequiresAuth(t *testing.T) {
	ts := setupTestServer(t)

	// Without token
	res, err := http.Get(ts.URL + "/api/v1/tai-khoan/me")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 without auth, got %d", res.StatusCode)
	}
}

func TestTaiKhoanMe_Success(t *testing.T) {
	ts := setupTestServer(t)

	// Register and login
	regBody := `{"email":"me@test.com","password":"secret123","ten":"User Me"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	loginBody := `{"email":"me@test.com","password":"secret123"}`
	res, _ = http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(loginBody))
	var loginResp DangNhapResponse
	json.NewDecoder(res.Body).Decode(&loginResp)
	res.Body.Close()

	// Fetch profile with token
	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/tai-khoan/me", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var profile TaiKhoanPublic
	if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if profile.Email != "me@test.com" {
		t.Errorf("expected me@test.com, got %s", profile.Email)
	}
}

func TestDangXuat_Success(t *testing.T) {
	ts := setupTestServer(t)

	// Register and login
	regBody := `{"email":"logout@test.com","password":"secret123","ten":"User"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	loginBody := `{"email":"logout@test.com","password":"secret123"}`
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
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	// Token should be invalid now
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

func TestHOSoBan_Activate_Success(t *testing.T) {
	ts := setupTestServer(t)

	// Register and login
	regBody := `{"email":"seller@test.com","password":"secret123","ten":"Seller"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	loginBody := `{"email":"seller@test.com","password":"secret123"}`
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

	var hs HOSoBan
	if err := json.NewDecoder(res.Body).Decode(&hs); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if hs.TaiKhoanID == "" {
		t.Error("expected non-empty tai_khoan_id")
	}
	if hs.TrangThai != TrangThaiHOSoBanKichHoat {
		t.Errorf("expected %s, got %s", TrangThaiHOSoBanKichHoat, hs.TrangThai)
	}
}

func TestHOSoBan_RequiresAuth(t *testing.T) {
	ts := setupTestServer(t)

	res, err := http.Post(ts.URL+"/api/v1/ho-so-nguoi-ban", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 without auth, got %d", res.StatusCode)
	}
}

func TestHOSoBan_Duplicate(t *testing.T) {
	ts := setupTestServer(t)

	// Register and login
	regBody := `{"email":"dubsell@test.com","password":"secret123","ten":"DupSeller"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	loginBody := `{"email":"dubsell@test.com","password":"secret123"}`
	res, _ = http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(loginBody))
	var loginResp DangNhapResponse
	json.NewDecoder(res.Body).Decode(&loginResp)
	res.Body.Close()

	// First activation
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/ho-so-nguoi-ban", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, _ = http.DefaultClient.Do(req)
	res.Body.Close()

	// Second activation → should fail
	req, _ = http.NewRequest("POST", ts.URL+"/api/v1/ho-so-nguoi-ban", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate seller profile, got %d", res.StatusCode)
	}
}

func TestAdminRole_AdminAccessGranted(t *testing.T) {
	ts := setupTestServerWithAdmin(t)

	// Login as admin
	loginBody := `{"email":"admin@test.com","password":"admin123"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(loginBody))
	var loginResp DangNhapResponse
	json.NewDecoder(res.Body).Decode(&loginResp)
	res.Body.Close()

	// Verify admin role
	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/tai-khoan/me", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var profile TaiKhoanPublic
	json.NewDecoder(res.Body).Decode(&profile)

	if profile.VaiTro != VaiTroAdmin {
		t.Errorf("expected admin role, got %s", profile.VaiTro)
	}
}

func TestProtectedEndpoint_InvalidToken(t *testing.T) {
	ts := setupTestServer(t)

	tests := []struct {
		name  string
		token string
	}{
		{"malformed token", "not-a-bearer"},
		{"missing Authorization header", ""},
		{"garbage token", "Bearer this-is-not-a-valid-token"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", ts.URL+"/api/v1/tai-khoan/me", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected 401, got %d", res.StatusCode)
			}
		})
	}
}

func TestDangNhap_ResponseHasSetCookie(t *testing.T) {
	ts := setupTestServer(t)

	// Register
	regBody := `{"email":"cookie@test.com","password":"secret123","ten":"Cookie"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	// Login
	loginBody := `{"email":"cookie@test.com","password":"secret123"}`
	res, err := http.Post(ts.URL+"/api/v1/dang-nhap", "application/json", strings.NewReader(loginBody))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Check for session cookie
	cookies := res.Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "session_token" {
			found = true
			if c.Value == "" {
				t.Error("session cookie value is empty")
			}
			if !c.HttpOnly {
				t.Error("session cookie should be HttpOnly")
			}
			break
		}
	}
	if !found {
		t.Error("expected session_token cookie in response")
	}
}

func TestDangXuat_ClearsCookie(t *testing.T) {
	ts := setupTestServer(t)

	// Register and login
	regBody := `{"email":"clearcookie@test.com","password":"secret123","ten":"Clear"}`
	res, _ := http.Post(ts.URL+"/api/v1/dang-ky", "application/json", strings.NewReader(regBody))
	res.Body.Close()

	loginBody := `{"email":"clearcookie@test.com","password":"secret123"}`
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
	defer res.Body.Close()

	// Check cookie is cleared
	cookies := res.Cookies()
	for _, c := range cookies {
		if c.Name == "session_token" {
			if c.Value != "" {
				t.Error("expected empty value for cleared session cookie")
			}
			if c.MaxAge >= 0 {
				t.Errorf("expected negative MaxAge for cleared cookie, got %d", c.MaxAge)
			}
			break
		}
	}
}
