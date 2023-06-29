package handler

import (
	"github.com/Rhymond/go-money"
)

type PhoneBrandReferenceListDTO struct {
	ID    int    `json:"id"`
	Brand string `json:"brand"`
	Name  string `json:"name"`
}

type CaseTypeListDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IconImgPath string `json:"icon_img_path"`
}

type DiscountListDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Rate       int    `json:"rate"`
	ValidUntil string `json:"valid_until"`
}

type PhoneCaseListDTO struct {
	ID                  int                        `json:"id"`
	Price               money.Money                `json:"price"`
	CaseScaffoldImgPath string                     `json:"case_scaffold_img_path"`
	InventoryStatus     string                     `json:"inventory_status"`
	Discount            DiscountListDTO            `json:"discount"`
	CreatedAt           string                     `json:"created_at"`
	PhoneBrandReference PhoneBrandReferenceListDTO `json:"phone_brand_reference"`
	CaseType            CaseTypeListDTO            `json:"case_type"`
	CreatedBy           string                     `json:"created_by"`
}

type PhoneCaseUpdateDTO struct {
	Price           money.Money `json:"price"`
	ScaffoldImgPath string      `json:"scaffold_img_path"`
	InventoryStatus string      `json:"inventory_status"`
	PhoneBrandRefID int         `json:"phone_brand_ref_id"`
	CaseTypeID      int         `json:"case_type_id"`
}

type PhoneCaseCreateDTO struct {
	Price           money.Money `json:"price"`
	ScaffoldImgPath string      `json:"scaffold_img_path"`
	InventoryStatus string      `json:"inventory_status"`
	PhoneBrandRefID int         `json:"phone_brand_ref_id"`
	CaseTypeID      int         `json:"case_type_id"`
}
