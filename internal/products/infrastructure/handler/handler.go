package handler

import (
	"github.com/CrissAlvarezH/fundart-api/internal/common"
	"github.com/CrissAlvarezH/fundart-api/internal/products/application/services"
	"github.com/CrissAlvarezH/fundart-api/internal/products/domain"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PhoneCaseHandler struct {
	service services.PhoneCaseService
}

func NewPhoneCaseHandler(service services.PhoneCaseService) PhoneCaseHandler {
	return PhoneCaseHandler{
		service: service,
	}
}

func (h *PhoneCaseHandler) AddRoutes(g *gin.RouterGroup) {
	g.GET("/cases", h.List)
	g.POST("/cases/", h.Create)
	g.GET("/cases/:id", h.GetByID)
	g.PUT("/cases/:id/", h.Update)
	g.DELETE("/cases/:id/", h.Delete)
}

func (h *PhoneCaseHandler) Create(c *gin.Context) {
	var body PhoneCaseCreateDTO
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authU, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "auth error"})
		return
	}
	user := authU.(users.User)

	phoneCase, err := h.service.CreatePhoneCase(
		body.Price,
		body.ScaffoldImgPath,
		domain.InventoryStatus(body.InventoryStatus),
		domain.PhoneBrandReferenceID(body.PhoneBrandRefID),
		domain.CaseTypeID(body.CaseTypeID),
		user.ID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapToPhoneCaseListDTO(phoneCase))
}

func (h *PhoneCaseHandler) List(c *gin.Context) {
	pageParams, err := common.GetPaginationParams(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phoneCases, count := h.service.ListPhoneCases(
		pageParams.Filters, pageParams.Limit, pageParams.Offset,
	)

	phoneCasesDTO := mapToPhoneCasesListDTO(phoneCases)

	c.JSON(http.StatusOK, gin.H{
		"pagination": common.PaginationJson(count, pageParams.Page, pageParams.Page),
		"result":     phoneCasesDTO,
	})
}

func (h *PhoneCaseHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not valid number"})
		return
	}

	phoneCase, ok := h.service.GetPhoneCaseByID(domain.PhoneCaseID(id))
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone case does not exists"})
		return
	}

	c.JSON(http.StatusOK, mapToPhoneCaseListDTO(phoneCase))
}

func (h *PhoneCaseHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not valid number"})
		return
	}

	var body PhoneCaseUpdateDTO
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phoneCase, err := h.service.UpdatePhoneCase(
		domain.PhoneCaseID(id),
		body.Price,
		body.ScaffoldImgPath,
		domain.InventoryStatus(body.InventoryStatus),
		domain.PhoneBrandReferenceID(body.PhoneBrandRefID),
		domain.CaseTypeID(body.CaseTypeID),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapToPhoneCaseListDTO(phoneCase))
}

func (h *PhoneCaseHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not valid number"})
		return
	}

	if err := h.service.DeletePhoneCase(domain.PhoneCaseID(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

type DiscountHandler struct {
	service services.PhoneCaseService
}

func (h *DiscountHandler) AddRoutes(g *gin.RouterGroup) {
	g.GET("/discounts", h.List)
	g.GET("/discounts/:id", h.GetByID)
	g.PUT("/discounts/:id/", h.Update)
	g.DELETE("/discounts/:id/", h.Delete)
}

func (h *DiscountHandler) List(c *gin.Context) {

}

func (h *DiscountHandler) GetByID(c *gin.Context) {

}

func (h *DiscountHandler) Update(c *gin.Context) {

}

func (h *DiscountHandler) Delete(c *gin.Context) {

}

type PhoneBrandHandler struct {
	service services.PhoneCaseService
}

func (h *PhoneBrandHandler) AddRoutes(g *gin.RouterGroup) {
	g.GET("/brands", h.ListBrands)
	g.GET("/brands/:id", h.GetBrandByID)
	g.GET("/brands/:id/references", h.ListBrandReferences)
	g.PUT("/brands/:id/", h.UpdateBrand)
	g.DELETE("/brands/:id/", h.DeleteBrand)

	g.GET("/brand-references", h.ListAllBrandReferences)
	g.GET("/brand-references/:id", h.GetBrandReferencesByID)
	g.PUT("/brand-references/:id/", h.UpdateBrandReference)
	g.DELETE("/brand-references/:id/", h.DeleteBrandReference)
}

func (h *PhoneBrandHandler) ListBrands(c *gin.Context) {

}

func (h *PhoneBrandHandler) GetBrandByID(c *gin.Context) {

}

func (h *PhoneBrandHandler) ListBrandReferences(c *gin.Context) {

}

func (h *PhoneBrandHandler) UpdateBrand(c *gin.Context) {

}

func (h *PhoneBrandHandler) DeleteBrand(c *gin.Context) {

}

func (h *PhoneBrandHandler) ListAllBrandReferences(c *gin.Context) {

}

func (h *PhoneBrandHandler) GetBrandReferencesByID(c *gin.Context) {

}

func (h *PhoneBrandHandler) UpdateBrandReference(c *gin.Context) {

}

func (h *PhoneBrandHandler) DeleteBrandReference(c *gin.Context) {

}

type CaseTypeHandler struct {
	service services.PhoneCaseService
}

func (h *CaseTypeHandler) AddRoutes(g *gin.RouterGroup) {
	g.GET("/case-types", h.List)
	g.GET("/case-types/:id", h.GetByID)
	g.PUT("/case-types/:id/", h.Update)
	g.DELETE("/case-types/:id/", h.Delete)

	g.PUT("/case-types/:id/images", h.UpdateImages)
	g.DELETE("/case-types/:id/images/:img_id", h.DeleteImages)
}

func (h *CaseTypeHandler) List(c *gin.Context) {

}

func (h *CaseTypeHandler) GetByID(c *gin.Context) {

}

func (h *CaseTypeHandler) Update(c *gin.Context) {

}

func (h *CaseTypeHandler) Delete(c *gin.Context) {

}

func (h *CaseTypeHandler) UpdateImages(c *gin.Context) {

}

func (h *CaseTypeHandler) DeleteImages(c *gin.Context) {

}
