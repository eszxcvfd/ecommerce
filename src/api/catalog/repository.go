package catalog

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// CatalogRepository is the seam for product data access.
type CatalogRepository interface {
	// Products returns only approved (public) products.
	Products() ([]SanPhamSo, error)
	// Search returns approved products matching the given query, sorted accordingly.
	Search(query CatalogQuery) ([]SanPhamSo, error)
	// ProductByID returns one approved product by its ID, or nil if not found or not public.
	ProductByID(id string) (*SanPhamSo, error)
	// ProductsByCategory returns approved products in the given category,
	// excluding the given product ID, ordered by publish date descending
	// (newest first) with ID tie-break, limited to max results.
	ProductsByCategory(category DanhMuc, excludeID string, max int) ([]SanPhamSo, error)
	// DraftsBySeller returns all drafts owned by the given seller.
	DraftsBySeller(sellerID string) ([]SanPhamSo, error)
	// DraftByID returns one draft by ID, scoped to the given seller.
	DraftByID(id, sellerID string) (*SanPhamSo, error)
	// CreateDraft creates a new draft product owned by the given seller.
	CreateDraft(input DraftInput, sellerID string) (*SanPhamSo, error)
	// UpdateDraft updates an existing draft, scoped to the seller.
	UpdateDraft(id, sellerID string, input DraftUpdateInput) (*SanPhamSo, error)
	// DeleteDraft deletes a draft by ID, scoped to the seller.
	DeleteDraft(id, sellerID string) error
}


// memoryRepo holds in-memory product data.
type memoryRepo struct {
	products []SanPhamSo
}

// NewMemoryRepo creates an in-memory repository pre-loaded with the given products.
func NewMemoryRepo(products []SanPhamSo) CatalogRepository {
	return &memoryRepo{products: products}
}

func (r *memoryRepo) Products() ([]SanPhamSo, error) {
	// Return only approved/public products
	var result []SanPhamSo
	for _, p := range r.products {
		if p.TrangThai == TrangThaiDaDuyet {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *memoryRepo) ProductByID(id string) (*SanPhamSo, error) {
	for _, p := range r.products {
		if p.ID == id && p.TrangThai == TrangThaiDaDuyet {
			return &p, nil
		}
	}
	return nil, nil
}

// Search returns approved products filtered and sorted by the given query.
func (r *memoryRepo) Search(query CatalogQuery) ([]SanPhamSo, error) {
	// Start with all approved products
	candidates, err := r.Products()
	if err != nil {
		return nil, err
	}

	// Apply text search (q)
	if query.Q != "" {
		normalized := normalizeSearch(query.Q)
		var filtered []SanPhamSo
		for _, p := range candidates {
			haystack := normalizeSearch(p.Ten + " " + p.MoTa)
			if strings.Contains(haystack, normalized) {
				filtered = append(filtered, p)
			}
		}
		candidates = filtered
	}

	// Filter by danh_muc
	if query.DanhMuc != "" {
		var filtered []SanPhamSo
		for _, p := range candidates {
			if string(p.DanhMuc) == query.DanhMuc {
				filtered = append(filtered, p)
			}
		}
		candidates = filtered
	}

	// Filter by dinh_dang
	if query.DinhDang != "" {
		var filtered []SanPhamSo
		for _, p := range candidates {
			for _, ext := range p.DinhDang {
				if ext == query.DinhDang {
					filtered = append(filtered, p)
					break
				}
			}
		}
		candidates = filtered
	}

	// Filter by price range (inclusive)
	if query.MinXu != nil || query.MaxXu != nil {
		var filtered []SanPhamSo
		for _, p := range candidates {
			price := int64(0)
			if !p.Gia.MienPhi {
				price = p.Gia.SoXu
			}
			if query.MinXu != nil && price < *query.MinXu {
				continue
			}
			if query.MaxXu != nil && price > *query.MaxXu {
				continue
			}
			filtered = append(filtered, p)
		}
		candidates = filtered
	}

	// Apply sort
	sortOrder := SortOrder(query.Sort)
	if sortOrder == "" {
		sortOrder = SortNewest
	}

	switch sortOrder {
	case SortPopular:
		sort.SliceStable(candidates, func(i, j int) bool {
			if candidates[i].SoLuotTai != candidates[j].SoLuotTai {
				return candidates[i].SoLuotTai > candidates[j].SoLuotTai
			}
			return candidates[i].ID < candidates[j].ID
		})
	case SortPriceAsc:
		sort.SliceStable(candidates, func(i, j int) bool {
			pi := priceValue(candidates[i])
			pj := priceValue(candidates[j])
			if pi != pj {
				return pi < pj
			}
			return candidates[i].ID < candidates[j].ID
		})
	case SortPriceDesc:
		sort.SliceStable(candidates, func(i, j int) bool {
			pi := priceValue(candidates[i])
			pj := priceValue(candidates[j])
			if pi != pj {
				return pi > pj
			}
			return candidates[i].ID < candidates[j].ID
		})
	case SortRating:
		sort.SliceStable(candidates, func(i, j int) bool {
			// Unrated products (0 ratings) sort last
			iRated := candidates[i].SoLuongDanhGia > 0
			jRated := candidates[j].SoLuongDanhGia > 0
			if iRated != jRated {
				return iRated // rated before unrated
			}
			if !iRated && !jRated {
				return candidates[i].ID < candidates[j].ID
			}
			if candidates[i].DiemDanhGia != candidates[j].DiemDanhGia {
				return candidates[i].DiemDanhGia > candidates[j].DiemDanhGia
			}
			return candidates[i].ID < candidates[j].ID
		})
	default: // newest
		sort.SliceStable(candidates, func(i, j int) bool {
			if !candidates[i].NgayTao.Equal(candidates[j].NgayTao) {
				return candidates[i].NgayTao.After(candidates[j].NgayTao)
			}
			return candidates[i].ID < candidates[j].ID
		})
	}

	return candidates, nil
}

// ProductsByCategory returns approved products filtered by category, excluding one ID.
// Results are ordered by ngay_dang DESC, id ASC, limited to max.
func (r *memoryRepo) ProductsByCategory(category DanhMuc, excludeID string, max int) ([]SanPhamSo, error) {
	var result []SanPhamSo
	for _, p := range r.products {
		if p.TrangThai != TrangThaiDaDuyet {
			continue
		}
		if p.DanhMuc != category {
			continue
		}
		if p.ID == excludeID {
			continue
		}
		result = append(result, p)
	}

	// Sort by ngay_dang descending, then id ascending (stable tie-break)
	sort.SliceStable(result, func(i, j int) bool {
		if !result[i].NgayDang.Equal(result[j].NgayDang) {
			return result[i].NgayDang.After(result[j].NgayDang)
		}
		return result[i].ID < result[j].ID
	})

	// Limit
	if len(result) > max {
		result = result[:max]
	}

	return result, nil
}


// ---------------------------------------------------------------------------
// Memory repo: draft methods
// ---------------------------------------------------------------------------

func (r *memoryRepo) DraftsBySeller(sellerID string) ([]SanPhamSo, error) {
	var result []SanPhamSo
	for _, p := range r.products {
		if p.NguoiBanID == sellerID && p.TrangThai == TrangThaiDangSoan {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *memoryRepo) DraftByID(id, sellerID string) (*SanPhamSo, error) {
	for _, p := range r.products {
		if p.ID == id && p.NguoiBanID == sellerID && p.TrangThai == TrangThaiDangSoan {
			return &p, nil
		}
	}
	return nil, nil
}

func (r *memoryRepo) CreateDraft(input DraftInput, sellerID string) (*SanPhamSo, error) {
	id := fmt.Sprintf("draft_%d", time.Now().UnixNano())

	var dinhDang []string
	formatSet := make(map[string]bool)
	for _, tep := range input.Tep {
		if !formatSet[tep.DinhDang] {
			dinhDang = append(dinhDang, tep.DinhDang)
			formatSet[tep.DinhDang] = true
		}
	}

	var tepList []Tep
	for _, t := range input.Tep {
		tepList = append(tepList, Tep{
			TenTep:         t.TenTep,
			DinhDang:       t.DinhDang,
			DungLuongBytes: t.DungLuongBytes,
		})
	}

	sp := SanPhamSo{
		ID:        id,
		Ten:       input.Ten,
		MoTa:      input.MoTa,
		MoTaChiTiet: input.MoTaChiTiet,
		AnhDemo:   input.AnhDemo,
		Gia: Gia{
			MienPhi: input.MienPhi,
			SoXu:    input.SoXu,
		},
		DanhMuc:  input.DanhMuc,
		DinhDang: dinhDang,
		GiayPhep: input.GiayPhep,
		NgayTao:  time.Now(),
		NguoiBanID: sellerID,
		Tep:      tepList,
		TrangThai: TrangThaiDangSoan,
	}

	r.products = append(r.products, sp)
	return &sp, nil
}

func (r *memoryRepo) UpdateDraft(id, sellerID string, input DraftUpdateInput) (*SanPhamSo, error) {
	for i, p := range r.products {
		if p.ID == id && p.NguoiBanID == sellerID && p.TrangThai == TrangThaiDangSoan {
			if input.Ten != nil {
				r.products[i].Ten = *input.Ten
			}
			if input.MoTa != nil {
				r.products[i].MoTa = *input.MoTa
			}
			if input.MoTaChiTiet != nil {
				r.products[i].MoTaChiTiet = *input.MoTaChiTiet
			}
			if input.AnhDemo != nil {
				r.products[i].AnhDemo = *input.AnhDemo
			}
			if input.MienPhi != nil {
				r.products[i].Gia.MienPhi = *input.MienPhi
			}
			if input.SoXu != nil {
				r.products[i].Gia.SoXu = *input.SoXu
			}
			if input.DanhMuc != nil {
				r.products[i].DanhMuc = *input.DanhMuc
			}
			if input.GiayPhep != nil {
				r.products[i].GiayPhep = *input.GiayPhep
			}
			if input.Tep != nil {
				var tepList []Tep
				var dinhDang []string
				formatSet := make(map[string]bool)
				for _, t := range input.Tep {
					tepList = append(tepList, Tep{
						TenTep:         t.TenTep,
						DinhDang:       t.DinhDang,
						DungLuongBytes: t.DungLuongBytes,
					})
					if !formatSet[t.DinhDang] {
						dinhDang = append(dinhDang, t.DinhDang)
						formatSet[t.DinhDang] = true
					}
				}
				r.products[i].Tep = tepList
				r.products[i].DinhDang = dinhDang
			}
			result := r.products[i]
			return &result, nil
		}
	}
	return nil, nil
}

func (r *memoryRepo) DeleteDraft(id, sellerID string) error {
	for i, p := range r.products {
		if p.ID == id && p.NguoiBanID == sellerID && p.TrangThai == TrangThaiDangSoan {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("không tìm thấy bản nháp")
}
// normalizeSearch normalizes Vietnamese text for case-insensitive, accent-insensitive matching.
func normalizeSearch(s string) string {
	// Step 1: NFD decompose so combining marks become separate codepoints
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)))
	result, _, _ := transform.String(t, s)
	// Step 2: lowercase
	result = strings.ToLower(result)
	// Step 3: map đ/Đ to d
	result = strings.NewReplacer("đ", "d", "Đ", "d").Replace(result)
	return result
}

// priceValue returns the numeric price for sorting purposes (0 for free).
func priceValue(p SanPhamSo) int64 {
	if p.Gia.MienPhi {
		return 0
	}
	return p.Gia.SoXu
}
