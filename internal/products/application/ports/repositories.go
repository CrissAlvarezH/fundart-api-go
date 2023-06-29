package ports

import (
	"errors"
	"github.com/CrissAlvarezH/fundart-api/internal/products/domain"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/Rhymond/go-money"
	"time"
)

var (
	PhoneCaseDoesNotExists           = errors.New("phone case does not exists")
	DiscountDoesNotExists            = errors.New("discount does not exists")
	PhoneBrandDoesNotExists          = errors.New("phone brand does not exists")
	PhoneBrandReferenceDoesNotExists = errors.New("phone brand reference does not exists")
	BrandAlreadyExists               = errors.New("brand already exists")
	CaseTypeDoesNotExists            = errors.New("case type does not exists")
	CaseTypeImageDoesNotExists       = errors.New("case type image does not exists")
)

type PhoneCaseRepository interface {
	ListPhoneCases(filters map[string]string, limit int, offset int) ([]domain.PhoneCase, int)
	GetPhoneCaseByID(id domain.PhoneCaseID) (domain.PhoneCase, bool)
	CreatePhoneCase(
		price money.Money, scaffoldImgPath string, inventoryStatus domain.InventoryStatus,
		phoneBrandRefID domain.PhoneBrandReferenceID, caseTypeID domain.CaseTypeID,
		createdBy users.UserID,
	) (domain.PhoneCase, error)
	UpdatePhoneCase(
		ID domain.PhoneCaseID, price money.Money, scaffoldImgPath string, inventoryStatus domain.InventoryStatus,
		phoneBrandRefID domain.PhoneBrandReferenceID, caseTypeID domain.CaseTypeID,
	) (domain.PhoneCase, error)
	DeletePhoneCase(ID domain.PhoneCaseID) error

	AttachDiscount(caseID domain.PhoneCaseID, discountID domain.DiscountID) error

	ListAllDiscounts() []domain.Discount
	GetDiscountByID(ID domain.DiscountID) (domain.Discount, bool)
	CreateDiscount(name string, rate int, validUntil time.Time, createdBy users.UserID) (domain.Discount, error)
	UpdateDiscount(ID domain.DiscountID, rate int, validUntil time.Time) (domain.Discount, error)
	DeleteDiscount(ID domain.DiscountID) error
}

type PhoneBrandRepository interface {
	ListAllBrands() []domain.PhoneBrand
	UpdateBrand(brand domain.PhoneBrand, renameTo string) error
	CreateBrand(brand domain.PhoneBrand) error

	ListBrandReferences(brand domain.PhoneBrand) []domain.PhoneBrandReference
	GetBrandReferenceByID(ID domain.PhoneBrandReferenceID) (domain.PhoneBrandReference, bool)
	UpdateBrandReference(ID domain.PhoneBrandReferenceID, name string) (domain.PhoneBrandReference, error)
	CreateBrandReference(name string, brand domain.PhoneBrand) (domain.PhoneBrandReference, error)
}

type CaseTypeRepository interface {
	ListCaseTypes() []domain.CaseType
	GetCaseTypeByID(ID domain.CaseTypeID) (domain.CaseType, bool)
	UpdateCaseType(ID domain.CaseTypeID, name string) (domain.CaseType, error)
	CreateCaseType(name string, iconPath string) (domain.CaseType, error)

	ListCaseTypeImages(typeID domain.CaseTypeID) []domain.CaseTypeImage
	CreateCaseTypeImage(typeID domain.CaseTypeID, path string, orderPriority int) (domain.CaseTypeImage, error)
	UpdateCaseTypeImage(ID domain.CaseTypeImageID, path string) (domain.CaseTypeImage, error)
	DeleteCaseTypeImage(ID domain.CaseTypeImageID) error
}
