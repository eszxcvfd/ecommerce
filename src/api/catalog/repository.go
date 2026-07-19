package catalog

// CatalogRepository is the seam for product data access.
// The API behavior test injects a seeded in-memory implementation.
type CatalogRepository interface {
	// Products returns only approved (public) products.
	Products() []SanPhamSo
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
