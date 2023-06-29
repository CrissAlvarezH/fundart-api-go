package domain

import (
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/Rhymond/go-money"
	"time"
)

type PhoneCaseID int
type InventoryStatus string

const (
	PhoneCaseAvailable  InventoryStatus = "AVAILABLE"
	PhoneCaseOutOfStock InventoryStatus = "OUT_OF_STOCK"
	PhoneCaseInactive   InventoryStatus = "INACTIVE"
)

type PhoneCase struct {
	ID                  PhoneCaseID
	Price               money.Money
	CaseScaffoldImgPah  string
	InventoryStatus     InventoryStatus
	Discount            Discount
	CreatedAt           time.Time
	PhoneBrandReference PhoneBrandReference
	CaseType            CaseType
	CreatedBy           users.User
}

type PhoneBrand string
type PhoneBrandReferenceID int

type PhoneBrandReference struct {
	ID    PhoneBrandReferenceID
	Brand PhoneBrand
	Name  string
}

type CaseTypeID int

type CaseType struct {
	ID          CaseTypeID
	Name        string
	IconImgPath string
	Images      []CaseTypeImage
}

type CaseTypeImageID int

type CaseTypeImage struct {
	ID            CaseTypeImageID
	Path          string
	OrderPriority int
}

type DiscountID int

type Discount struct {
	ID         DiscountID
	Name       string
	Rate       int
	ValidUntil time.Time
	CreatedBy  users.User
	CreatedAt  time.Time
}
