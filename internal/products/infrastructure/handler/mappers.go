package handler

import (
	"github.com/CrissAlvarezH/fundart-api/internal/products/domain"
)

func mapToPhoneCasesListDTO(l []domain.PhoneCase) []PhoneCaseListDTO {
	result := make([]PhoneCaseListDTO, 0, len(l))
	for _, c := range l {
		result = append(result, mapToPhoneCaseListDTO(c))
	}
	return result
}

func mapToPhoneCaseListDTO(p domain.PhoneCase) PhoneCaseListDTO {
	discountDTO := mapToDiscountListDTO(p.Discount)
	phoneBrandRefDTO := mapToPhoneBrandRefListDTO(p.PhoneBrandReference)
	caseTypeDTO := mapToCaseTypeListDTO(p.CaseType)

	return PhoneCaseListDTO{
		ID:                  int(p.ID),
		Price:               p.Price,
		CaseScaffoldImgPath: p.CaseScaffoldImgPah,
		InventoryStatus:     string(p.InventoryStatus),
		Discount:            discountDTO,
		CreatedAt:           p.CreatedAt.Format("2006-01-02 15:04:05"),
		PhoneBrandReference: phoneBrandRefDTO,
		CaseType:            caseTypeDTO,
		CreatedBy:           p.CreatedBy.Email,
	}
}

func mapToDiscountListDTO(d domain.Discount) DiscountListDTO {
	return DiscountListDTO{
		ID:         int(d.ID),
		Name:       d.Name,
		Rate:       d.Rate,
		ValidUntil: d.ValidUntil.Format("2006-01-02 15:04:05"),
	}
}

func mapToPhoneBrandRefListDTO(r domain.PhoneBrandReference) PhoneBrandReferenceListDTO {
	return PhoneBrandReferenceListDTO{
		ID:    int(r.ID),
		Brand: string(r.Brand),
		Name:  r.Name,
	}
}

func mapToCaseTypeListDTO(c domain.CaseType) CaseTypeListDTO {
	return CaseTypeListDTO{
		ID:          int(c.ID),
		Name:        c.Name,
		IconImgPath: c.IconImgPath,
	}
}