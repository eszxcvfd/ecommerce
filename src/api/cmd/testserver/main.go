// Command testserver starts the catalog API with seeded data for e2e tests.
package main

import (
	"log"
	"os"

	"ecommerce/api/catalog"
)

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	repo := catalog.NewMemoryRepo(seed())
	log.Printf("Test API server listening on :%s", port)
	if err := catalog.StartServer(":"+port, repo); err != nil {
		log.Fatal(err)
	}
}

func seed() []catalog.SanPhamSo {
	return []catalog.SanPhamSo{
		{
			ID: "sp-001", Ten: "Bản vẽ nhà phố 3 tầng",
			AnhDemo: "/images/nha-pho.jpg",
			Gia:     catalog.Gia{MienPhi: true},
			DanhMuc: catalog.DanhMucKienTruc, TrangThai: catalog.TrangThaiDaDuyet,
			DiemDanhGia: 4.5, SoLuongDanhGia: 12,
		},
		{
			ID: "sp-002", Ten: "Mô hình khung thép tiền chế",
			AnhDemo: "/images/khung-thep.jpg",
			Gia:     catalog.Gia{MienPhi: false, SoXu: 15000},
			DanhMuc: catalog.DanhMucCoKhi, TrangThai: catalog.TrangThaiDaDuyet,
			DiemDanhGia: 0, SoLuongDanhGia: 0,
		},
		{
			ID: "sp-003", Ten: "Sơ đồ mạch Arduino điều khiển LED",
			AnhDemo: "/images/arduino-led.jpg",
			Gia:     catalog.Gia{MienPhi: true},
			DanhMuc: catalog.DanhMucDienTu, TrangThai: catalog.TrangThaiDaDuyet,
			DiemDanhGia: 3.8, SoLuongDanhGia: 5,
		},
		{
			ID: "sp-004", Ten: "Bộ icon phong cách tối giản",
			AnhDemo: "/images/icon-set.jpg",
			Gia:     catalog.Gia{MienPhi: false, SoXu: 5000},
			DanhMuc: catalog.DanhMucDoHoa, TrangThai: catalog.TrangThaiDaDuyet,
			DiemDanhGia: 4.2, SoLuongDanhGia: 8,
		},
		{
			ID: "sp-005", Ten: "Đồ án thiết kế cầu dầm BTCT",
			AnhDemo: "/images/cau-dam.jpg",
			Gia:     catalog.Gia{MienPhi: true},
			DanhMuc: catalog.DanhMucDoAn, TrangThai: catalog.TrangThaiDaDuyet,
			DiemDanhGia: 0, SoLuongDanhGia: 0,
		},
		{
			ID: "sp-006", Ten: "Luận văn thạc sĩ AI trong xây dựng",
			AnhDemo: "/images/luanvan-ai.jpg",
			Gia:     catalog.Gia{MienPhi: false, SoXu: 30000},
			DanhMuc: catalog.DanhMucLuanVan, TrangThai: catalog.TrangThaiDaDuyet,
			DiemDanhGia: 5.0, SoLuongDanhGia: 3,
		},
		// Non-approved — must not appear
		{
			ID: "sp-007", Ten: "Bản nháp chưa duyệt",
			AnhDemo: "", Gia: catalog.Gia{MienPhi: true},
			DanhMuc: catalog.DanhMucKienTruc, TrangThai: catalog.TrangThaiDangSoan,
		},
		{
			ID: "sp-008", Ten: "Đang chờ duyệt",
			AnhDemo: "", Gia: catalog.Gia{MienPhi: true},
			DanhMuc: catalog.DanhMucCoKhi, TrangThai: catalog.TrangThaiChoDuyet,
		},
	}
}
