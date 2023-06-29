package memoryrepo

import (
	"github.com/CrissAlvarezH/fundart-api/internal/products/application/ports"
	"github.com/CrissAlvarezH/fundart-api/internal/products/domain"
)

type MemoryPhoneBrandRepository struct {
	store MemoryStore
}

func NewMemoryPhoneBrandRepository(store MemoryStore) MemoryPhoneBrandRepository {
	return MemoryPhoneBrandRepository{
		store: store,
	}
}

func (m *MemoryPhoneBrandRepository) ListAllBrands() []domain.PhoneBrand {
	return mapToPhoneBrands(m.store.Brands)
}

func (m *MemoryPhoneBrandRepository) UpdateBrand(brand domain.PhoneBrand, renameTo string) error {
	for i, b := range m.store.Brands {
		if b == MemoryBrand(brand) {
			m.store.Brands[i] = MemoryBrand(renameTo)
			return nil
		}
	}
	return ports.PhoneBrandDoesNotExists
}

func (m *MemoryPhoneBrandRepository) CreateBrand(brand domain.PhoneBrand) error {
	for _, e := range m.store.Brands {
		if e == MemoryBrand(brand) {
			return ports.BrandAlreadyExists
		}
	}

	m.store.Brands = append(m.store.Brands, MemoryBrand(brand))
	return nil
}

func (m *MemoryPhoneBrandRepository) ListBrandReferences(brand domain.PhoneBrand) []domain.PhoneBrandReference {
	r := make([]MemoryPhoneBrandReference, 0)
	for _, e := range m.store.BrandReferences {
		if e.Brand == MemoryBrand(brand) {
			r = append(r, e)
		}
	}
	return mapToPhoneBrandRefs(r)
}

func (m *MemoryPhoneBrandRepository) GetBrandReferenceByID(
	ID domain.PhoneBrandReferenceID,
) (domain.PhoneBrandReference, bool) {
	for _, e := range m.store.BrandReferences {
		if e.ID == int(ID) {
			return mapToPhoneBrandRef(e), true
		}
	}
	return domain.PhoneBrandReference{}, false
}

func (m *MemoryPhoneBrandRepository) UpdateBrandReference(
	ID domain.PhoneBrandReferenceID, name string,
) (domain.PhoneBrandReference, error) {
	for i, e := range m.store.BrandReferences {
		if e.ID == int(ID) {
			m.store.BrandReferences[i].Name = name
			return mapToPhoneBrandRef(m.store.BrandReferences[i]), nil
		}
	}
	return domain.PhoneBrandReference{}, ports.PhoneBrandReferenceDoesNotExists
}

func (m *MemoryPhoneBrandRepository) CreateBrandReference(
	name string, brand domain.PhoneBrand,
) (domain.PhoneBrandReference, error) {
	lastID := 0
	if len(m.store.BrandReferences) > 0 {
		lastID = m.store.BrandReferences[len(m.store.BrandReferences)-1].ID
	}

	mR := MemoryPhoneBrandReference{
		ID:    lastID + 1,
		Brand: MemoryBrand(brand),
		Name:  name,
	}

	m.store.BrandReferences = append(m.store.BrandReferences, mR)
	return mapToPhoneBrandRef(mR), nil
}
