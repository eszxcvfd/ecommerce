package catalog

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
			AnhDemo: "/images/nha-pho.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 4.5, SoLuongDanhGia: 12,
		},
		{
			// Source: filethietke.vn — https://www.filethietke.vn/file-thiet-ke/mau-vach-cnc-dong-tien-hien-dai-222776.htm
			ID: "sp-017", Ten: "Mẫu vách CNC đồng tiền hiện đại",
			AnhDemo: "https://www.filethietke.vn/FilesUpload/Code/mau-vach-cnc-dong-tien-hien-dai-16434.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 100}, // demo price based on source (100 Xu filethietke.vn)
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 5.0, SoLuongDanhGia: 0,
		},
		{
			ID: "sp-003", Ten: "Sơ đồ mạch Arduino điều khiển LED",
			AnhDemo: "/images/arduino-led.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 3.8, SoLuongDanhGia: 5,
		},
		{
			ID: "sp-004", Ten: "Bộ icon phong cách tối giản",
			AnhDemo: "/images/icon-set.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 5000},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 4.2, SoLuongDanhGia: 8,
		},
		{
			ID: "sp-005", Ten: "Đồ án thiết kế cầu dầm BTCT",
			AnhDemo: "/images/cau-dam.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDoAn, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 0, SoLuongDanhGia: 0,
		},
		{
			ID: "sp-006", Ten: "Luận văn thạc sĩ AI trong xây dựng",
			AnhDemo: "/images/luanvan-ai.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 30000},
			DanhMuc: DanhMucLuanVan, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 5.0, SoLuongDanhGia: 3,
		},
		// --- Second batch of approved products (2 per category) ---
		{
			ID: "sp-011", Ten: "Phối cảnh khu nghỉ dưỡng",
			AnhDemo: "/images/khu-nghi-duong.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 25000},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 4.0, SoLuongDanhGia: 7,
		},
		{
			// Source: filethietke.vn — https://www.filethietke.vn/file-thiet-ke/mau-vach-cong-cnc-cay-nghe-thuat-222775.htm
			ID: "sp-018", Ten: "Mẫu vách cổng CNC cây nghệ thuật",
			AnhDemo: "https://www.filethietke.vn/FilesUpload/Code/mau-vach-cong-cnc-cay-nghe-thuat-164126.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 100}, // demo price based on source (100 Xu filethietke.vn)
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 5.0, SoLuongDanhGia: 0,
		},
		{
			ID: "sp-013", Ten: "Sơ đồ nguyên lý nguồn 5V",
			AnhDemo: "/images/nguon-5v.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 8000},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 0, SoLuongDanhGia: 0,
		},
		{
			ID: "sp-014", Ten: "Template thiết kế brochure",
			AnhDemo: "/images/brochure.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 4.8, SoLuongDanhGia: 15,
		},
		{
			ID: "sp-015", Ten: "Đồ án tốt nghiệp phần mềm quản lý thư viện",
			AnhDemo: "/images/qly-thu-vien.jpg",
			Gia:     Gia{MienPhi: false, SoXu: 20000},
			DanhMuc: DanhMucDoAn, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 3.5, SoLuongDanhGia: 2,
		},
		{
			ID: "sp-016", Ten: "Luận văn cử nhân kinh tế xây dựng",
			AnhDemo: "/images/luanvan-kinh-te.jpg",
			Gia:     Gia{MienPhi: true},
			DanhMuc: DanhMucLuanVan, TrangThai: TrangThaiDaDuyet,
			DiemDanhGia: 0, SoLuongDanhGia: 0,
		},
		// --- Non-approved products (must NOT appear) ---
		{
			ID: "sp-007", Ten: "Bản nháp chưa duyệt",
			AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucKienTruc, TrangThai: TrangThaiDangSoan,
		},
		{
			ID: "sp-008", Ten: "Đang chờ duyệt",
			AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucCoKhi, TrangThai: TrangThaiChoDuyet,
		},
		{
			ID: "sp-009", Ten: "Sản phẩm bị từ chối",
			AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucDienTu, TrangThai: TrangThaiBiTuChoi,
		},
		{
			ID: "sp-010", Ten: "Sản phẩm bị ẩn sau vi phạm",
			AnhDemo: "", Gia: Gia{MienPhi: true},
			DanhMuc: DanhMucDoHoa, TrangThai: TrangThaiBiAn,
		},
	}
}
