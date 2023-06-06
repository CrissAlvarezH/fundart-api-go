package handler

import (
	"github.com/CrissAlvarezH/fundart-api/internal/common"
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return UserHandler{service: service}
}

func (h *UserHandler) AddRoutes(e *gin.Engine) {
	e.GET("/users", h.List)
}

func (h *UserHandler) List(c *gin.Context) {
	pageParams, err := common.GetPaginationParams(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResult, userCount, err := h.service.List(
		pageParams.Filters, pageParams.Limit, pageParams.Offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "on list users"})
		return
	}

	usersDTO := make([]ListUserDTO, 0, len(userResult))
	for _, user := range userResult {
		usersDTO = append(usersDTO, MapToListUserDTO(user))
	}

	c.JSON(http.StatusOK, gin.H{
		"pagination": common.PaginationJson(userCount, pageParams.Page, pageParams.PageSize),
		"result":     usersDTO,
	})
}
