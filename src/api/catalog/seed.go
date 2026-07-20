package catalog

import "time"

// SeedData returns a deterministic dataset for local/demo API and e2e tests.
// It covers:
//   - At least two approved products per category
//   - Free and paid products
//   - Products with and without ratings
//   - Non-approved products that must not appear (draft, pending, rejected, hidden)
func SeedData() []SanPhamSo {
	return []SanPhamSo{
		// --- Approved products (should appear) ---
		{
			ID: "sp-001", Ten: "Bản vẽ nhà phố 3 tầng",
			MoTa:    "Bản vẽ kiến trúc nhà phố 3 tầng, bao gồm mặt bằng, mặt đứng, mặt cắt.",
			AnhDemo: "/images/nha-pho.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"dwg", "pdf"},
			DiemDanhGia: 4.5, SoLuongDanhGia: 12,
			NgayTao: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), SoLuotTai: 120,
		},
		{
			// Source: filethietke.vn — https://www.filethietke.vn/file-thiet-ke/mau-vach-cnc-dong-tien-hien-dai-222776.htm
			ID: "sp-017", Ten: "Mẫu vách CNC đồng tiền hiện đại",
			MoTa:    "Mẫu vách CNC trang trí nội thất với hoa văn đồng tiền hiện đại, tệp DXF.",
			AnhDemo: "https://www.filethietke.vn/FilesUpload/Code/mau-vach-cnc-dong-tien-hien-dai-16434.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 100},
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"dxf"},
			DiemDanhGia: 5.0, SoLuongDanhGia: 0,
			NgayTao: time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC), SoLuotTai: 45,
		},
		{
			ID: "sp-003", Ten: "Sơ đồ mạch Arduino điều khiển LED",
			MoTa:    "Sơ đồ nguyên lý và bố trí bo mạch Arduino cho dự án điều khiển LED RGB.",
			AnhDemo: "/images/arduino-led.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"brd", "sch"},
			DiemDanhGia: 3.8, SoLuongDanhGia: 5,
			NgayTao: time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC), SoLuotTai: 230,
		},
		{
			ID: "sp-004", Ten: "Bộ icon phong cách tối giản",
			MoTa:    "Bộ 200 icon phong cách tối giản, phù hợp cho web và ứng dụng di động.",
			AnhDemo: "/images/icon-set.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 5000},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"svg", "png", "ai"},
			DiemDanhGia: 4.2, SoLuongDanhGia: 8,
			NgayTao: time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC), SoLuotTai: 310,
		},
		{
			ID: "sp-005", Ten: "Đồ án thiết kế cầu dầm BTCT",
			MoTa:    "Đồ án tốt nghiệp thiết kế cầu dầm bê tông cốt thép, bao gồm bản vẽ và thuyết minh.",
			AnhDemo: "/images/cau-dam.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDoAn, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"pdf", "dwg"},
			DiemDanhGia: 0, SoLuongDanhGia: 0,
			NgayTao: time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC), SoLuotTai: 89,
		},
		{
			ID: "sp-006", Ten: "Luận văn thạc sĩ AI trong xây dựng",
			MoTa:    "Luận văn thạc sĩ về ứng dụng trí tuệ nhân tạo trong quản lý dự án xây dựng.",
			AnhDemo: "/images/luanvan-ai.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 30000},
			DanhMuc: DanhMucLuanVan, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"pdf", "docx"},
			DiemDanhGia: 5.0, SoLuongDanhGia: 3,
			NgayTao: time.Date(2026, 7, 6, 0, 0, 0, 0, time.UTC), SoLuotTai: 67,
		},
		// --- Second batch of approved products (2 per category) ---
		{
			ID: "sp-011", Ten: "Phối cảnh khu nghỉ dưỡng",
			MoTa:    "Phối cảnh 3D khu nghỉ dưỡng cao cấp, file SketchUp và hình ảnh render.",
			AnhDemo: "/images/khu-nghi-duong.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 25000},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"skp", "png"},
			DiemDanhGia: 4.0, SoLuongDanhGia: 7,
			NgayTao: time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC), SoLuotTai: 43,
		},
		{
			// Source: filethietke.vn — https://www.filethietke.vn/file-thiet-ke/mau-vach-cong-cnc-cay-nghe-thuat-222775.htm
			ID: "sp-018", Ten: "Mẫu vách cổng CNC cây nghệ thuật",
			MoTa:    "Mẫu vách cổng CNC thiết kế cây nghệ thuật, phù hợp trang trí sân vườn.",
			AnhDemo: "https://www.filethietke.vn/FilesUpload/Code/mau-vach-cong-cnc-cay-nghe-thuat-164126.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 100},
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"dxf"},
			DiemDanhGia: 5.0, SoLuongDanhGia: 0,
			NgayTao: time.Date(2026, 7, 8, 0, 0, 0, 0, time.UTC), SoLuotTai: 32,
		},
		{
			ID: "sp-013", Ten: "Sơ đồ nguyên lý nguồn 5V",
			MoTa:    "Sơ đồ nguyên lý mạch nguồn 5V ổn áp, bao gồm sơ đồ mạch in và danh sách linh kiện.",
			AnhDemo: "/images/nguon-5v.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 8000},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"pdf", "sch"},
			DiemDanhGia: 0, SoLuongDanhGia: 0,
			NgayTao: time.Date(2026, 7, 9, 0, 0, 0, 0, time.UTC), SoLuotTai: 156,
		},
		{
			ID: "sp-014", Ten: "Template thiết kế brochure",
			MoTa:    "Template brochure đa năng cho doanh nghiệp, thiết kế chuyên nghiệp dễ chỉnh sửa.",
			AnhDemo: "/images/brochure.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"ai", "psd"},
			DiemDanhGia: 4.8, SoLuongDanhGia: 15,
			NgayTao: time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC), SoLuotTai: 520,
		},
		{
			ID: "sp-015", Ten: "Đồ án tốt nghiệp phần mềm quản lý thư viện",
			MoTa:    "Đồ án tốt nghiệp xây dựng phần mềm quản lý thư viện, đầy đủ báo cáo và mã nguồn.",
			AnhDemo: "/images/qly-thu-vien.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 20000},
			DanhMuc: DanhMucDoAn, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"pdf"},
			DiemDanhGia: 3.5, SoLuongDanhGia: 2,
			NgayTao: time.Date(2026, 7, 11, 0, 0, 0, 0, time.UTC), SoLuotTai: 78,
		},
		{
			ID: "sp-016", Ten: "Luận văn cử nhân kinh tế xây dựng",
			MoTa:    "Luận văn cử nhân ngành kinh tế xây dựng, phân tích thị trường bất động sản.",
			AnhDemo: "/images/luanvan-kinh-te.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucLuanVan, TrangThai: TrangThaiDaDuyet,
			DinhDang:    []string{"pdf", "docx"},
			DiemDanhGia: 0, SoLuongDanhGia: 0,
			NgayTao: time.Date(2026, 7, 12, 0, 0, 0, 0, time.UTC), SoLuotTai: 94,
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
