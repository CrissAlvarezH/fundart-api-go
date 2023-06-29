package memoryrepo

import (
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"time"
)

type MemoryPhoneCase struct {
	ID                    int
	Price                 string
	CaseScaffoldImgPath   string
	InventoryStatus       string
	DiscountID            int
	CreatedAt             time.Time
	PhoneBrandReferenceID int
	CaseTypeID            int
	CreatedBy             users.User
}

type MemoryDiscount struct {
	ID         int
	Name       string
	Rate       int
	ValidUntil time.Time
	CreateByID int
	CreatedAt  time.Time
}

type MemoryBrand string

type MemoryPhoneBrandReference struct {
	ID    int
	Brand MemoryBrand
	Name  string
}

type MemoryCaseType struct {
	ID          int
	Name        string
	IconImgPath string
	Images      []MemoryCaseTypeImage
}

type MemoryCaseTypeImage struct {
	ID            int
	Path          string
	OrderPriority int
	CaseTypeID    int
}
