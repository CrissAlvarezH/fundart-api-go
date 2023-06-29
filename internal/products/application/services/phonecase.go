package services

import (
	"github.com/CrissAlvarezH/fundart-api/internal/products/application/ports"
	"github.com/CrissAlvarezH/fundart-api/internal/products/domain"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/Rhymond/go-money"
	"time"
)

type PhoneCaseService struct {
	caseRepo     ports.PhoneCaseRepository
	brandRepo    ports.PhoneBrandRepository
	caseTypeRepo ports.CaseTypeRepository
}

func NewPhoneCaseService(
	caseRepo ports.PhoneCaseRepository, brandRepo ports.PhoneBrandRepository,
	caseTypeRepo ports.CaseTypeRepository,
) PhoneCaseService {
	return PhoneCaseService{
		caseRepo:     caseRepo,
		brandRepo:    brandRepo,
		caseTypeRepo: caseTypeRepo,
	}
}

func (s *PhoneCaseService) ListPhoneCases(filters map[string]string, limit int, offset int) ([]domain.PhoneCase, int) {
	return s.caseRepo.ListPhoneCases(filters, limit, offset)
}

func (s *PhoneCaseService) GetPhoneCaseByID(ID domain.PhoneCaseID) (domain.PhoneCase, bool) {
	return s.caseRepo.GetPhoneCaseByID(ID)
}

func (s *PhoneCaseService) CreatePhoneCase(
	price money.Money, scaffoldImgPath string, inventoryStatus domain.InventoryStatus,
	phoneBrandRefID domain.PhoneBrandReferenceID, caseTypeID domain.CaseTypeID,
	createdBy users.UserID,
) (domain.PhoneCase, error) {
	return s.caseRepo.CreatePhoneCase(price, scaffoldImgPath, inventoryStatus, phoneBrandRefID, caseTypeID, createdBy)
}

func (s *PhoneCaseService) UpdatePhoneCase(
	ID domain.PhoneCaseID, price money.Money, scaffoldImgPath string, inventoryStatus domain.InventoryStatus,
	phoneBrandRefID domain.PhoneBrandReferenceID, caseTypeID domain.CaseTypeID,
) (domain.PhoneCase, error) {
	return s.caseRepo.UpdatePhoneCase(ID, price, scaffoldImgPath, inventoryStatus, phoneBrandRefID, caseTypeID)
}

func (s *PhoneCaseService) DeletePhoneCase(ID domain.PhoneCaseID) error {
	return s.caseRepo.DeletePhoneCase(ID)
}

func (s *PhoneCaseService) AttachDiscount(caseID domain.PhoneCaseID, discountID domain.DiscountID) error {
	return s.caseRepo.AttachDiscount(caseID, discountID)
}

func (s *PhoneCaseService) ListAllDiscounts() []domain.Discount {
	return s.caseRepo.ListAllDiscounts()
}

func (s *PhoneCaseService) CreateDiscount(name string, rate int, validUntil time.Time, createdBy users.UserID) (domain.Discount, error) {
	return s.caseRepo.CreateDiscount(name, rate, validUntil, createdBy)
}

func (s *PhoneCaseService) UpdateDiscount(ID domain.DiscountID, rate int, validUntil time.Time, user users.User) (domain.Discount, error) {
	return s.caseRepo.UpdateDiscount(ID, rate, validUntil)
}

func (s *PhoneCaseService) DeleteDiscount(ID domain.DiscountID) error {
	return s.caseRepo.DeleteDiscount(ID)
}

func (s *PhoneCaseService) ListPhoneBrands() []domain.PhoneBrand {
	return s.brandRepo.ListAllBrands()
}

func (s *PhoneCaseService) UpdatePhoneBrand(brand domain.PhoneBrand, renameTo string) error {
	return s.brandRepo.UpdateBrand(brand, renameTo)
}

func (s *PhoneCaseService) CreatePhoneBrand(brand domain.PhoneBrand) error {
	return s.brandRepo.CreateBrand(brand)
}

func (s *PhoneCaseService) ListBrandReferences(brand domain.PhoneBrand) []domain.PhoneBrandReference {
	return s.brandRepo.ListBrandReferences(brand)
}

func (s *PhoneCaseService) UpdateBrandReference(ID domain.PhoneBrandReferenceID, name string) (domain.PhoneBrandReference, error) {
	return s.brandRepo.UpdateBrandReference(ID, name)
}

func (s *PhoneCaseService) CreateBrandReference(name string, brand domain.PhoneBrand) (domain.PhoneBrandReference, error) {
	return s.brandRepo.CreateBrandReference(name, brand)
}

func (s *PhoneCaseService) ListCaseTypes() []domain.CaseType {
	return s.caseTypeRepo.ListCaseTypes()
}

func (s *PhoneCaseService) UpdateCaseType(ID domain.CaseTypeID, name string) (domain.CaseType, error) {
	return s.caseTypeRepo.UpdateCaseType(ID, name)
}

func (s *PhoneCaseService) CreateCaseType(name string, iconPath string) (domain.CaseType, error) {
	return s.caseTypeRepo.CreateCaseType(name, iconPath)
}

func (s *PhoneCaseService) ListCaseTypeImages(typeID domain.CaseTypeID) []domain.CaseTypeImage {
	return s.caseTypeRepo.ListCaseTypeImages(typeID)
}

func (s *PhoneCaseService) CreateCaseTypeImage(typeID domain.CaseTypeID, path string, orderPriority int) (domain.CaseTypeImage, error) {
	return s.caseTypeRepo.CreateCaseTypeImage(typeID, path, orderPriority)
}

func (s *PhoneCaseService) UpdateCaseTypeImage(ID domain.CaseTypeImageID, path string) (domain.CaseTypeImage, error) {
	return s.caseTypeRepo.UpdateCaseTypeImage(ID, path)
}

func (s *PhoneCaseService) DeleteCaseTypeImage(ID domain.CaseTypeImageID) error {
	return s.caseTypeRepo.DeleteCaseTypeImage(ID)
}
