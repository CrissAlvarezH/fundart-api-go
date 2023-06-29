package memoryrepo

import (
	"github.com/CrissAlvarezH/fundart-api/internal/products/application/ports"
	"github.com/CrissAlvarezH/fundart-api/internal/products/domain"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/Rhymond/go-money"
	"strconv"
	"time"
)

type MemoryPhoneCaseRepository struct {
	store     MemoryStore
	brandRepo ports.PhoneBrandRepository
	caseType  ports.CaseTypeRepository
}

func NewMemoryPhoneCaseRepository(
	store MemoryStore, brandRepo ports.PhoneBrandRepository, caseTypeRepo ports.CaseTypeRepository,
) MemoryPhoneCaseRepository {
	return MemoryPhoneCaseRepository{
		store:     store,
		brandRepo: brandRepo,
		caseType:  caseTypeRepo,
	}
}

func (r *MemoryPhoneCaseRepository) ListPhoneCases(
	filters map[string]string, limit int, offset int,
) ([]domain.PhoneCase, int) {
	filtered := make([]MemoryPhoneCase, 0, len(r.store.PhoneCases))

	if len(filters) == 0 {
		filtered = r.store.PhoneCases
	} else {
		for _, v := range r.store.PhoneCases {
			price, ok := filters["price"]
			if ok && v.Price == price {
				filtered = append(filtered, v)
			}

			status, ok := filters["inventory_status"]
			if ok && v.InventoryStatus == status {
				filtered = append(filtered, v)
			}
		}
	}

	data := make([]MemoryPhoneCase, 0, len(filtered))
	if len(filtered) >= limit && len(filtered) >= offset {
		data = filtered[offset:limit]
	} else if len(filtered) >= offset && len(filtered) < limit {
		data = filtered[offset:]
	}

	result := make([]domain.PhoneCase, 0, len(filtered))
	for _, p := range data {
		discount, _ := r.GetDiscountByID(domain.DiscountID(p.DiscountID))
		brandRef, _ := r.brandRepo.GetBrandReferenceByID(domain.PhoneBrandReferenceID(p.PhoneBrandReferenceID))
		caseType, _ := r.caseType.GetCaseTypeByID(domain.CaseTypeID(p.CaseTypeID))
		phoneCase := mapToPhoneCase(p, discount, brandRef, caseType)
		result = append(result, phoneCase)
	}

	return result, len(filtered)
}

func (r *MemoryPhoneCaseRepository) GetPhoneCaseByID(ID domain.PhoneCaseID) (domain.PhoneCase, bool) {
	for _, p := range r.store.PhoneCases {
		if p.ID == int(ID) {
			discount, _ := r.GetDiscountByID(domain.DiscountID(p.DiscountID))
			brandRef, _ := r.brandRepo.GetBrandReferenceByID(domain.PhoneBrandReferenceID(p.PhoneBrandReferenceID))
			caseType, _ := r.caseType.GetCaseTypeByID(domain.CaseTypeID(p.CaseTypeID))
			phoneCase := mapToPhoneCase(p, discount, brandRef, caseType)
			return phoneCase, true
		}
	}
	return domain.PhoneCase{}, false
}

func (r *MemoryPhoneCaseRepository) CreatePhoneCase(
	price money.Money, scaffoldImgPath string, inventoryStatus domain.InventoryStatus,
	phoneBrandRefID domain.PhoneBrandReferenceID, caseTypeID domain.CaseTypeID,
	createdBy users.UserID,
) (domain.PhoneCase, error) {
	lastID := 0
	if len(r.store.PhoneCases) > 0 {
		lastID = r.store.PhoneCases[len(r.store.PhoneCases)-1].ID
	}
	mP := MemoryPhoneCase{
		ID:                    lastID + 1,
		Price:                 strconv.Itoa(int(price.Amount())),
		CaseScaffoldImgPath:   scaffoldImgPath,
		InventoryStatus:       string(inventoryStatus),
		DiscountID:            0,
		CreatedAt:             time.Time{},
		PhoneBrandReferenceID: int(phoneBrandRefID),
		CaseTypeID:            int(caseTypeID),
		CreatedBy:             users.User{},
	}

	discount := domain.Discount{}
	brandRef, _ := r.brandRepo.GetBrandReferenceByID(phoneBrandRefID)
	caseType, _ := r.caseType.GetCaseTypeByID(caseTypeID)
	phoneCase := mapToPhoneCase(mP, discount, brandRef, caseType)
	return phoneCase, nil
}

func (r *MemoryPhoneCaseRepository) UpdatePhoneCase(
	ID domain.PhoneCaseID, price money.Money, scaffoldImgPath string, inventoryStatus domain.InventoryStatus,
	phoneBrandRefID domain.PhoneBrandReferenceID, caseTypeID domain.CaseTypeID,
) (domain.PhoneCase, error) {
	for i, p := range r.store.PhoneCases {
		if p.ID == int(ID) {
			r.store.PhoneCases[i].Price = strconv.Itoa(int(price.Amount()))
			r.store.PhoneCases[i].CaseScaffoldImgPath = scaffoldImgPath
			r.store.PhoneCases[i].InventoryStatus = string(inventoryStatus)
			r.store.PhoneCases[i].PhoneBrandReferenceID = int(phoneBrandRefID)

			discount, _ := r.GetDiscountByID(domain.DiscountID(p.DiscountID))
			brandRef, _ := r.brandRepo.GetBrandReferenceByID(domain.PhoneBrandReferenceID(p.PhoneBrandReferenceID))
			caseType, _ := r.caseType.GetCaseTypeByID(domain.CaseTypeID(p.CaseTypeID))
			phoneCase := mapToPhoneCase(p, discount, brandRef, caseType)
			return phoneCase, nil
		}
	}
	return domain.PhoneCase{}, ports.PhoneCaseDoesNotExists
}

func (r *MemoryPhoneCaseRepository) DeletePhoneCase(ID domain.PhoneCaseID) error {
	result := make([]MemoryPhoneCase, 0, len(r.store.PhoneCases))
	found := false

	for _, p := range r.store.PhoneCases {
		if p.ID != int(ID) {
			result = append(result, p)
			found = true
		}

	}

	if !found {
		return ports.PhoneCaseDoesNotExists
	}
	return nil
}

func (r *MemoryPhoneCaseRepository) AttachDiscount(caseID domain.PhoneCaseID, discountID domain.DiscountID) error {
	_, ok := r.GetDiscountByID(discountID)
	if !ok {
		return ports.DiscountDoesNotExists
	}

	for i, p := range r.store.PhoneCases {
		if p.ID == int(caseID) {
			r.store.PhoneCases[i].DiscountID = int(discountID)
			return nil
		}
	}

	return ports.PhoneCaseDoesNotExists
}

func (r *MemoryPhoneCaseRepository) ListAllDiscounts() []domain.Discount {
	return mapToDiscounts(r.store.Discounts)
}

func (r *MemoryPhoneCaseRepository) GetDiscountByID(ID domain.DiscountID) (domain.Discount, bool) {
	for _, d := range r.store.Discounts {
		if d.ID == int(ID) {
			return mapToDiscount(d), true
		}
	}
	return domain.Discount{}, false
}

func (r *MemoryPhoneCaseRepository) CreateDiscount(
	name string, rate int, validUntil time.Time, createdBy users.UserID,
) (domain.Discount, error) {
	lastID := 0
	if len(r.store.Discounts) > 0 {
		lastID = r.store.Discounts[len(r.store.Discounts)-1].ID
	}
	mD := MemoryDiscount{
		ID:         lastID + 1,
		Name:       name,
		Rate:       rate,
		ValidUntil: validUntil,
		CreateByID: int(createdBy),
		CreatedAt:  time.Time{},
	}
	r.store.Discounts = append(r.store.Discounts, mD)
	return mapToDiscount(mD), nil
}

func (r *MemoryPhoneCaseRepository) UpdateDiscount(
	ID domain.DiscountID, rate int, validUntil time.Time,
) (domain.Discount, error) {
	for i, d := range r.store.Discounts {
		if d.ID == int(ID) {
			r.store.Discounts[i].Rate = rate
			r.store.Discounts[i].ValidUntil = validUntil
			return mapToDiscount(d), nil
		}
	}
	return domain.Discount{}, ports.DiscountDoesNotExists
}

func (r *MemoryPhoneCaseRepository) DeleteDiscount(ID domain.DiscountID) error {
	result := make([]MemoryDiscount, 0, len(r.store.Discounts))
	found := false
	for _, d := range r.store.Discounts {
		if d.ID == int(ID) {
			result = append(result, d)
			found = true
		}
	}
	if !found {
		return ports.DiscountDoesNotExists
	}
	r.store.Discounts = result
	return nil
}
