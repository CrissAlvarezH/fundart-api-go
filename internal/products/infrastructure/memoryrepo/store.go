package memoryrepo

type MemoryStore struct {
	PhoneCases      []MemoryPhoneCase
	Discounts       []MemoryDiscount
	Brands          []MemoryBrand
	BrandReferences []MemoryPhoneBrandReference
	CaseTypes       []MemoryCaseType
	CaseTypeImage   []MemoryCaseTypeImage
}

func NewMemoryStore(
	phoneCases []MemoryPhoneCase, discounts []MemoryDiscount, brands []MemoryBrand,
	brandRefs []MemoryPhoneBrandReference, caseTypes []MemoryCaseType,
) MemoryStore {
	return MemoryStore{
		PhoneCases:      phoneCases,
		Discounts:       discounts,
		Brands:          brands,
		BrandReferences: brandRefs,
		CaseTypes:       caseTypes,
	}
}
