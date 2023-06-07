package handler

import (
	"github.com/CrissAlvarezH/fundart-api/internal/common"
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/services"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return UserHandler{service: service}
}

func (h *UserHandler) AddRoutes(e *gin.Engine) {
	e.GET("/api/v1/users", h.List)
	e.GET("/api/v1/users/:id", h.GetByID)
	e.POST("/api/v1/users/", h.Register)
	e.POST("/api/v1/users/:id/verification-code/", h.ValidateVerificationCode)
	e.PUT("/api/v1/users/:id/", h.Update)

	e.GET("/api/v1/users/:id/addresses", h.ListAddresses)
	e.POST("/api/v1/users/:id/addresses/", h.AddAddress)
}

func (h *UserHandler) List(c *gin.Context) {
	pageParams, err := common.GetPaginationParams(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResult, userCount := h.service.List(
		pageParams.Filters, pageParams.Limit, pageParams.Offset,
	)

	usersDTO := make([]ListUserDTO, 0, len(userResult))
	for _, user := range userResult {
		usersDTO = append(usersDTO, MapToListUserDTO(user))
	}

	c.JSON(http.StatusOK, gin.H{
		"pagination": common.PaginationJson(userCount, pageParams.Page, pageParams.PageSize),
		"result":     usersDTO,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not valid number"})
		return
	}
	user, ok := h.service.GetByID(users.UserID(id))
	if ok == false {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, MapToListUserDTO(user))
}

func (h *UserHandler) Register(c *gin.Context) {
	var body RegisterUserDTO
	if err := c.BindJSON(&body); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Add(
		body.Name, body.Email, body.Password, body.Phone,
		false, []users.ScopeName{},
	)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SendVerificationCode(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapToListUserDTO(user))
}

func (h *UserHandler) ValidateVerificationCode(c *gin.Context) {
	var body ValidateVerificationCodeDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	userWasActivated := h.service.ValidateVerificationCode(users.UserID(userID), body.Code)

	if userWasActivated == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"details": "User was activated successfully"})
}

func (h *UserHandler) Update(c *gin.Context) {
	var body UpdateUserDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	user, err := h.service.Update(
		users.UserID(userID), body.Name, body.Email, body.Phone, body.Scopes,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, MapToListUserDTO(user))
}

func (h *UserHandler) ListAddresses(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	addresses := h.service.ListAddresses(users.UserID(userID))

	addressesDTO := make([]ListAddressDTO, 0, len(addresses))
	for _, a := range addresses {
		addressesDTO = append(addressesDTO, MapToListAddressDTO(a))
	}

	c.JSON(http.StatusOK, addressesDTO)
}

func (h *UserHandler) AddAddress(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	var body CreateAddressDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.service.AddAddress(
		users.UserID(userID), body.Department, body.City,
		body.Address, body.ReceiverPhone, body.ReceiverName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	addresses := h.service.ListAddresses(users.UserID(userID))
	c.JSON(http.StatusCreated, addresses)
}
