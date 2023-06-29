package memoryrepo

import (
	"github.com/CrissAlvarezH/fundart-api/internal/products/domain"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/Rhymond/go-money"
	"strconv"
)

func mapToPhoneCase(
	m MemoryPhoneCase, discount domain.Discount, brandRef domain.PhoneBrandReference,
	caseType domain.CaseType,
) domain.PhoneCase {
	price, _ := strconv.Atoi(m.Price)
	return domain.PhoneCase{
		ID:                  domain.PhoneCaseID(m.ID),
		Price:               *money.New(int64(price), money.COP),
		CaseScaffoldImgPah:  m.CaseScaffoldImgPath,
		InventoryStatus:     domain.InventoryStatus(m.InventoryStatus),
		Discount:            discount,
		CreatedAt:           m.CreatedAt,
		PhoneBrandReference: brandRef,
		CaseType:            caseType,
		CreatedBy:           users.User{},
	}
}

func mapToDiscounts(m []MemoryDiscount) []domain.Discount {
	res := make([]domain.Discount, 0, len(m))
	for _, d := range m {
		res = append(res, mapToDiscount(d))
	}
	return res
}

func mapToDiscount(m MemoryDiscount) domain.Discount {
	return domain.Discount{
		ID:         domain.DiscountID(m.ID),
		Name:       m.Name,
		Rate:       m.Rate,
		ValidUntil: m.ValidUntil,
		CreatedBy:  users.User{},
		CreatedAt:  m.CreatedAt,
	}
}

func mapToPhoneBrands(m []MemoryBrand) []domain.PhoneBrand {
	b := make([]domain.PhoneBrand, 0, len(m))
	for _, e := range m {
		b = append(b, domain.PhoneBrand(e))
	}
	return b
}

func mapToPhoneBrandRefs(m []MemoryPhoneBrandReference) []domain.PhoneBrandReference {
	l := make([]domain.PhoneBrandReference, 0, len(m))
	for _, e := range m {
		l = append(l, mapToPhoneBrandRef(e))
	}
	return l
}

func mapToPhoneBrandRef(m MemoryPhoneBrandReference) domain.PhoneBrandReference {
	return domain.PhoneBrandReference{
		ID:    domain.PhoneBrandReferenceID(m.ID),
		Brand: domain.PhoneBrand(m.Brand),
		Name:  m.Name,
	}
}

func mapToCaseTypes(m []MemoryCaseType) []domain.CaseType {
	l := make([]domain.CaseType, 0, len(m))
	for _, e := range m {
		l = append(l, mapToCaseType(e))
	}
	return l
}

func mapToCaseType(m MemoryCaseType) domain.CaseType {
	return domain.CaseType{
		ID:          domain.CaseTypeID(m.ID),
		Name:        m.Name,
		IconImgPath: m.IconImgPath,
		Images:      mapToCaseTypeImages(m.Images),
	}
}

func mapToCaseTypeImages(m []MemoryCaseTypeImage) []domain.CaseTypeImage {
	res := make([]domain.CaseTypeImage, 0, len(m))
	for _, e := range m {
		res = append(res, mapToCaseTypeImage(e))
	}
	return res
}

func mapToCaseTypeImage(m MemoryCaseTypeImage) domain.CaseTypeImage {
	return domain.CaseTypeImage{
		ID:            domain.CaseTypeImageID(m.ID),
		Path:          m.Path,
		OrderPriority: m.OrderPriority,
	}
}
