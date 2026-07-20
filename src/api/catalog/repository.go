package catalog

import (
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// CatalogRepository is the seam for product data access.
// The API behavior test injects a seeded in-memory implementation.
type CatalogRepository interface {
	// Products returns only approved (public) products.
	Products() []SanPhamSo
	// Search returns approved products matching the given query, sorted accordingly.
	Search(query CatalogQuery) []SanPhamSo
}

// memoryRepo holds in-memory product data.
type memoryRepo struct {
	products []SanPhamSo
}

// NewMemoryRepo creates an in-memory repository pre-loaded with the given products.
func NewMemoryRepo(products []SanPhamSo) CatalogRepository {
	return &memoryRepo{products: products}
}

func (r *memoryRepo) Products() []SanPhamSo {
	// Return only approved/public products
	var result []SanPhamSo
	for _, p := range r.products {
		if p.TrangThai == TrangThaiDaDuyet {
			result = append(result, p)
		}
	}
	return result
}

// Search returns approved products filtered and sorted by the given query.
func (r *memoryRepo) Search(query CatalogQuery) []SanPhamSo {
	// Start with all approved products
	candidates := r.Products()

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

	return candidates
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
