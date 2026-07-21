package catalog

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"
)

//go:embed seed_data.json
var seedDataJSON []byte

// SeedData returns a deterministic dataset for local/demo API and e2e tests.
// It covers:
//   - At least two approved products per category
//   - Free and paid products
//   - Products with and without ratings
//   - Products with and without new optional metadata fields
//   - Non-approved products that must not appear (draft, pending, rejected, hidden)
func SeedData() []SanPhamSo {
	return []SanPhamSo{
		// --- Approved products (should appear) ---
		{
			ID: "sp-001", Ten: "Bản vẽ nhà phố 3 tầng",
			MoTa: "Bản vẽ kiến trúc nhà phố 3 tầng, bao gồm mặt bằng, mặt đứng, mặt cắt.",
			MoTaChiTiet: "Bộ bản vẽ đầy đủ cho nhà phố 3 tầng phong cách hiện đại.",
			AnhDemo: "/images/nha-pho.jpg",
			Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"dwg", "pdf"},
			DiemDanhGia:    4.5, SoLuongDanhGia: 12,
			NgayTao:  time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 120,
			GiayPhep:        "Giấy phép Sử dụng Cá nhân",
			NguoiBanHienThi: "Kiến Trúc Sư Nguyễn Văn A",
			Tep: []Tep{
				{TenTep: "MB-TangTret.dwg", DinhDang: "dwg", DungLuongBytes: 2500000},
				{TenTep: "MatDung.dwg", DinhDang: "dwg", DungLuongBytes: 1800000},
				{TenTep: "ThuyetMinh.pdf", DinhDang: "pdf", DungLuongBytes: 1200000},
			},
		},
		{
			ID: "sp-017", Ten: "Mẫu vách CNC đồng tiền hiện đại",
			MoTa:    "Mẫu vách CNC trang trí nội thất với hoa văn đồng tiền hiện đại, tệp DXF.",
			MoTaChiTiet: "Mẫu vách CNC trang trí nội thất phong cách hiện đại với họa tiết đồng tiền cách điệu.",
			AnhDemo: "https://www.filethietke.vn/FilesUpload/Code/mau-vach-cnc-dong-tien-hien-dai-16434.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 100},
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"dxf"},
			DiemDanhGia:    5.0, SoLuongDanhGia: 0,
			NgayTao:  time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 45,
			GiayPhep:        "Giấy phép Sử dụng Thương mại",
			NguoiBanHienThi: "Xưởng CNC Hoàng Gia",
			Tep: []Tep{
				{TenTep: "Vach-CNC-DongTien.dxf", DinhDang: "dxf", DungLuongBytes: 500000},
			},
		},
		{
			ID: "sp-003", Ten: "Sơ đồ mạch Arduino điều khiển LED",
			MoTa: "Sơ đồ nguyên lý và bố trí bo mạch Arduino cho dự án điều khiển LED RGB.",
			MoTaChiTiet: "Dự án Arduino điều khiển LED RGB qua Bluetooth, gồm sơ đồ nguyên lý, PCB layout, code mẫu.",
			AnhDemo: "/images/arduino-led.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"brd", "sch"},
			DiemDanhGia:    3.8, SoLuongDanhGia: 5,
			NgayTao:  time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 230,
			GiayPhep:        "Giấy phép Nguồn Mở MIT",
			NguoiBanHienThi: "DIY Electronics",
			Tep: []Tep{
				{TenTep: "Arduino-LED.brd", DinhDang: "brd", DungLuongBytes: 800000},
				{TenTep: "Arduino-LED.sch", DinhDang: "sch", DungLuongBytes: 300000},
			},
		},
		{
			ID: "sp-004", Ten: "Bộ icon phong cách tối giản",
			MoTa: "Bộ 200 icon phong cách tối giản, phù hợp cho web và ứng dụng di động.",
			MoTaChiTiet: "Bộ 200 biểu tượng thiết kế theo phong cách tối giản hiện đại, bao gồm SVG, PNG, AI.",
			AnhDemo: "/images/icon-set.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 5000},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"svg", "png", "ai"},
			DiemDanhGia:    4.2, SoLuongDanhGia: 8,
			NgayTao:  time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 8, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 310,
			GiayPhep:        "Giấy phép Sử dụng Thương mại",
			NguoiBanHienThi: "Design Studio",
			Tep: []Tep{
				{TenTep: "Icons-Minimal.ai", DinhDang: "ai", DungLuongBytes: 3500000},
				{TenTep: "Icons-SVG.zip", DinhDang: "zip", DungLuongBytes: 1500000},
			},
		},
		{
			ID: "sp-005", Ten: "Đồ án thiết kế cầu dầm BTCT",
			MoTa: "Đồ án tốt nghiệp thiết kế cầu dầm bê tông cốt thép, bao gồm bản vẽ và thuyết minh.",
			AnhDemo: "/images/cau-dam.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDoAn, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"pdf", "dwg"},
			DiemDanhGia:    0, SoLuongDanhGia: 0,
			NgayTao:  time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 89,
			GiayPhep:        "Giấy phép Tham khảo Học thuật",
			NguoiBanHienThi: "Sinh Viên Xây Dựng",
			Tep: []Tep{
				{TenTep: "CauDam-BTCT.pdf", DinhDang: "pdf", DungLuongBytes: 3500000},
				{TenTep: "BanVe-CauDam.dwg", DinhDang: "dwg", DungLuongBytes: 4200000},
			},
		},
		{
			ID: "sp-006", Ten: "Luận văn thạc sĩ AI trong xây dựng",
			MoTa: "Luận văn thạc sĩ về ứng dụng trí tuệ nhân tạo trong quản lý dự án xây dựng.",
			MoTaChiTiet: "Luận văn thạc sĩ chuyên ngành Quản lý Xây dựng, nghiên cứu ứng dụng AI trong dự đoán rủi ro tiến độ.",
			AnhDemo: "/images/luanvan-ai.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 30000},
			DanhMuc: DanhMucLuanVan, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"pdf", "docx"},
			DiemDanhGia:    5.0, SoLuongDanhGia: 3,
			NgayTao:  time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 67,
			GiayPhep:        "Giấy phép Tham khảo Học thuật",
			NguoiBanHienThi: "Thạc sĩ Lê Văn B",
			Tep: []Tep{
				{TenTep: "LuanVan-AI-XayDung.pdf", DinhDang: "pdf", DungLuongBytes: 5000000},
				{TenTep: "LuanVan-AI-XayDung.docx", DinhDang: "docx", DungLuongBytes: 2800000},
			},
		},
		// --- Second batch of approved products (2 per category) ---
		{
			ID: "sp-011", Ten: "Phối cảnh khu nghỉ dưỡng",
			MoTa:    "Phối cảnh 3D khu nghỉ dưỡng cao cấp, file SketchUp và hình ảnh render.",
			MoTaChiTiet: "Phối cảnh 3D khu nghỉ dưỡng cao cấp ven biển, mô hình SketchUp và ảnh render chất lượng cao.",
			AnhDemo: "/images/khu-nghi-duong.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 25000},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"skp", "png"},
			DiemDanhGia:    4.0, SoLuongDanhGia: 7,
			NgayTao:  time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 43,
			GiayPhep:        "Giấy phép Sử dụng Thương mại",
			NguoiBanHienThi: "Kiến Trúc Sư Trần Văn C",
			Tep: []Tep{
				{TenTep: "KhuNghiDuong.skp", DinhDang: "skp", DungLuongBytes: 5000000},
				{TenTep: "Render-MatTien.png", DinhDang: "png", DungLuongBytes: 3500000},
			},
		},
		{
			ID: "sp-018", Ten: "Mẫu vách cổng CNC cây nghệ thuật",
			MoTa:    "Mẫu vách cổng CNC thiết kế cây nghệ thuật, phù hợp trang trí sân vườn.",
			MoTaChiTiet: "Mẫu vách cổng CNC với họa tiết cây nghệ thuật, phong cách thiên nhiên.",
			AnhDemo: "https://www.filethietke.vn/FilesUpload/Code/mau-vach-cong-cnc-cay-nghe-thuat-164126.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 100},
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"dxf"},
			DiemDanhGia:    5.0, SoLuongDanhGia: 0,
			NgayTao:  time.Date(2026, 7, 8, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 12, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 32,
			GiayPhep:        "Giấy phép Sử dụng Thương mại",
			NguoiBanHienThi: "Xưởng CNC Hoàng Gia",
			Tep: []Tep{
				{TenTep: "VachCong-CayNT.dxf", DinhDang: "dxf", DungLuongBytes: 750000},
			},
		},
		{
			ID: "sp-013", Ten: "Sơ đồ nguyên lý nguồn 5V",
			MoTa:    "Sơ đồ nguyên lý mạch nguồn 5V ổn áp, bao gồm sơ đồ mạch in và danh sách linh kiện.",
			MoTaChiTiet: "Thiết kế mạch nguồn 5V ổn áp dùng LM2596.",
			AnhDemo: "/images/nguon-5v.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 8000},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"pdf", "sch"},
			DiemDanhGia:    0, SoLuongDanhGia: 0,
			NgayTao:  time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 13, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 156,
			GiayPhep:        "Giấy phép Nguồn Mở",
			NguoiBanHienThi: "DIY Electronics",
			Tep: []Tep{
				{TenTep: "Nguon-5V-PCB.pdf", DinhDang: "pdf", DungLuongBytes: 900000},
				{TenTep: "Nguon-5V.sch", DinhDang: "sch", DungLuongBytes: 250000},
			},
		},
		{
			ID: "sp-014", Ten: "Template thiết kế brochure",
			MoTa:    "Template brochure đa năng cho doanh nghiệp, thiết kế chuyên nghiệp dễ chỉnh sửa.",
			MoTaChiTiet: "Brochure template đa năng phù hợp cho doanh nghiệp ở nhiều lĩnh vực.",
			AnhDemo: "/images/brochure.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"ai", "psd"},
			DiemDanhGia:    4.8, SoLuongDanhGia: 15,
			NgayTao:  time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 520,
			GiayPhep:        "Giấy phép Sử dụng Cá nhân và Thương mại",
			NguoiBanHienThi: "Design Studio",
			Tep: []Tep{
				{TenTep: "Brochure-Template.ai", DinhDang: "ai", DungLuongBytes: 4500000},
				{TenTep: "Brochure-Template.psd", DinhDang: "psd", DungLuongBytes: 12000000},
			},
		},
		{
			ID: "sp-015", Ten: "Đồ án tốt nghiệp phần mềm quản lý thư viện",
			MoTa:    "Đồ án tốt nghiệp xây dựng phần mềm quản lý thư viện, đầy đủ báo cáo và mã nguồn.",
			MoTaChiTiet: "Đồ án tốt nghiệp hệ thống quản lý thư viện với các chức năng quản lý sách, độc giả, mượn trả.",
			AnhDemo: "/images/qly-thu-vien.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 20000},
			DanhMuc: DanhMucDoAn, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"pdf"},
			DiemDanhGia:    3.5, SoLuongDanhGia: 2,
			NgayTao:  time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 78,
			GiayPhep:        "Giấy phép Tham khảo Học thuật",
			NguoiBanHienThi: "Sinh Viên CNTT",
			Tep: []Tep{
				{TenTep: "DoAn-QLTV.pdf", DinhDang: "pdf", DungLuongBytes: 4500000},
			},
		},
		{
			ID: "sp-016", Ten: "Luận văn cử nhân kinh tế xây dựng",
			MoTa:    "Luận văn cử nhân ngành kinh tế xây dựng, phân tích thị trường bất động sản.",
			AnhDemo: "/images/luanvan-kinh-te.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucLuanVan, TrangThai: TrangThaiDaDuyet,
			DinhDang:       []string{"pdf", "docx"},
			DiemDanhGia:    0, SoLuongDanhGia: 0,
			NgayTao:  time.Date(2026, 7, 12, 0, 0, 0, 0, time.UTC),
			NgayDang:  time.Date(2026, 7, 16, 0, 0, 0, 0, time.UTC),
			SoLuotTai: 94,
			Tep: []Tep{
				{TenTep: "LuanVan-KTXD.pdf", DinhDang: "pdf", DungLuongBytes: 3800000},
				{TenTep: "LuanVan-KTXD.docx", DinhDang: "docx", DungLuongBytes: 2100000},
			},
		},
		// --- Non-approved products (must NOT appear) ---
		{
			ID: "sp-007", Ten: "Bản nháp chưa duyệt",
			MoTa: "", AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDangSoan,
		},
		{
			ID: "sp-008", Ten: "Đang chờ duyệt",
			MoTa: "", AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiChoDuyet,
		},
		{
			ID: "sp-009", Ten: "Sản phẩm bị từ chối",
			MoTa: "", AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiBiTuChoi,
		},
		{
			ID: "sp-010", Ten: "Sản phẩm bị ẩn sau vi phạm",
			MoTa: "", AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiBiAn,
		},
	}
}

// SeedFromJSON parses the embedded versioned JSON contract and returns the seed products.
// It validates the JSON and reports errors for corrupt embedded data.
func SeedFromJSON() ([]SanPhamSo, error) {
	cf, err := ValidateCatalogJSON(seedDataJSON)
	if err != nil {
		return nil, fmt.Errorf("embedded seed JSON: %w", err)
	}
	products := make([]SanPhamSo, len(cf.Products))
	for i, cp := range cf.Products {
		sp, err := cp.ToSanPhamSo()
		if err != nil {
			return nil, fmt.Errorf("seed product %s: %w", cp.ID, err)
		}
		products[i] = sp
	}
	return products, nil
}

// SeedSQLite inserts seed data into the given SQLite database when it is empty.
// It uses the embedded versioned JSON contract. If the database already contains
// products, it skips seeding (idempotent on non-empty databases).
func SeedSQLite(db *sql.DB) error {
	// Check if the database already has products.
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM san_pham_so").Scan(&count); err != nil {
		return fmt.Errorf("check seed count: %w", err)
	}
	if count > 0 {
		return nil // already seeded, skip
	}

	products, err := SeedFromJSON()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin seed tx: %w", err)
	}
	defer tx.Rollback()

	for _, p := range products {
		var ngayDangStr string
		if !p.NgayDang.IsZero() {
			ngayDangStr = p.NgayDang.Format(time.RFC3339)
		}

		_, err := tx.Exec(`
			INSERT INTO san_pham_so (id, ten, mo_ta, mo_ta_chi_tiet, anh_demo, mien_phi, so_xu, danh_muc,
			                         diem_danh_gia, so_luong_danh_gia, ngay_tao, so_luot_tai, trang_thai,
			                         giay_phep, nguoi_ban_hien_thi, ngay_dang,
			                         ten_search, mo_ta_search)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, p.ID, p.Ten, p.MoTa, p.MoTaChiTiet, p.AnhDemo, boolToInt(p.Gia.MienPhi), p.Gia.SoXu,
			string(p.DanhMuc), p.DiemDanhGia, p.SoLuongDanhGia,
			p.NgayTao.Format(time.RFC3339), p.SoLuotTai, string(p.TrangThai),
			p.GiayPhep, p.NguoiBanHienThi, ngayDangStr,
			normalizeSearch(p.Ten), normalizeSearch(p.MoTa),
		)
		if err != nil {
			return fmt.Errorf("insert seed product %s: %w", p.ID, err)
		}

		for _, ext := range p.DinhDang {
			_, err := tx.Exec(
				"INSERT OR IGNORE INTO san_pham_dinh_dang (san_pham_id, dinh_dang) VALUES (?, ?)",
				p.ID, ext,
			)
			if err != nil {
				return fmt.Errorf("insert format %s for %s: %w", ext, p.ID, err)
			}
		}

		for _, f := range p.Tep {
			_, err := tx.Exec(
				"INSERT OR IGNORE INTO san_pham_tep (san_pham_id, ten_tep, dinh_dang, dung_luong_bytes) VALUES (?, ?, ?, ?)",
				p.ID, f.TenTep, f.DinhDang, f.DungLuongBytes,
			)
			if err != nil {
				return fmt.Errorf("insert file %s for %s: %w", f.TenTep, p.ID, err)
			}
		}
	}

	return tx.Commit()
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
