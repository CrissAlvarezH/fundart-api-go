package memoryrepo

import (
	"github.com/CrissAlvarezH/fundart-api/internal/products/application/ports"
	"github.com/CrissAlvarezH/fundart-api/internal/products/domain"
)

type MemoryCaseTypeRepository struct {
	store MemoryStore
}

func NewMemoryCaseTypeRepository(store MemoryStore) MemoryCaseTypeRepository {
	return MemoryCaseTypeRepository{
		store: store,
	}
}

func (m *MemoryCaseTypeRepository) ListCaseTypes() []domain.CaseType {
	l := make([]domain.CaseType, 0, len(m.store.CaseTypes))
	for _, e := range m.store.CaseTypes {
		l = append(l, mapToCaseType(e))
	}
	return l
}

func (m *MemoryCaseTypeRepository) GetCaseTypeByID(ID domain.CaseTypeID) (domain.CaseType, bool) {
	for _, e := range m.store.CaseTypes {
		if e.ID == int(ID) {
			caseType := mapToCaseType(e)
			caseType.Images = m.ListCaseTypeImages(ID)
			return caseType, true
		}
	}
	return domain.CaseType{}, false
}

func (m *MemoryCaseTypeRepository) UpdateCaseType(ID domain.CaseTypeID, name string) (domain.CaseType, error) {
	for i, e := range m.store.CaseTypes {
		if e.ID == int(ID) {
			m.store.CaseTypes[i].Name = name
			return mapToCaseType(m.store.CaseTypes[i]), nil
		}
	}

	return domain.CaseType{}, ports.CaseTypeDoesNotExists
}

func (m *MemoryCaseTypeRepository) CreateCaseType(name string, iconPath string) (domain.CaseType, error) {
	lastID := 0
	if len(m.store.CaseTypes) > 0 {
		lastID = m.store.CaseTypes[len(m.store.CaseTypes)-1].ID
	}

	mC := MemoryCaseType{
		ID:          lastID + 1,
		Name:        name,
		IconImgPath: iconPath,
		Images:      []MemoryCaseTypeImage{},
	}

	m.store.CaseTypes = append(m.store.CaseTypes, mC)
	return mapToCaseType(mC), nil
}

func (m *MemoryCaseTypeRepository) ListCaseTypeImages(typeID domain.CaseTypeID) []domain.CaseTypeImage {
	l := make([]domain.CaseTypeImage, 0, len(m.store.CaseTypeImage))
	for _, e := range m.store.CaseTypeImage {
		if e.CaseTypeID == int(typeID) {
			l = append(l, mapToCaseTypeImage(e))
		}
	}
	return l
}

func (m *MemoryCaseTypeRepository) CreateCaseTypeImage(
	typeID domain.CaseTypeID, path string, orderPriority int,
) (domain.CaseTypeImage, error) {
	_, ok := m.GetCaseTypeByID(typeID)
	if !ok {
		return domain.CaseTypeImage{}, ports.CaseTypeDoesNotExists
	}

	lastID := 0
	if len(m.store.CaseTypeImage) > 0 {
		lastID = int(m.store.CaseTypeImage[len(m.store.CaseTypeImage)-1].ID)
	}

	caseImg := MemoryCaseTypeImage{
		ID:            lastID + 1,
		Path:          path,
		OrderPriority: orderPriority,
		CaseTypeID:    int(typeID),
	}
	m.store.CaseTypeImage = append(m.store.CaseTypeImage, caseImg)

	return mapToCaseTypeImage(caseImg), nil
}

func (m *MemoryCaseTypeRepository) UpdateCaseTypeImage(
	ID domain.CaseTypeImageID, path string,
) (domain.CaseTypeImage, error) {
	for i, e := range m.store.CaseTypeImage {
		if e.ID == int(ID) {
			m.store.CaseTypeImage[i].Path = path
			return mapToCaseTypeImage(e), nil
		}
	}
	return domain.CaseTypeImage{}, ports.CaseTypeImageDoesNotExists
}

func (m *MemoryCaseTypeRepository) DeleteCaseTypeImage(ID domain.CaseTypeImageID) error {
	l := make([]MemoryCaseTypeImage, 0)
	found := false
	for _, e := range m.store.CaseTypeImage {
		if e.ID != int(ID) {
			l = append(l, e)
			found = true
		}
	}
	if !found {
		return ports.CaseTypeImageDoesNotExists
	}

	m.store.CaseTypeImage = l
	return nil
}
