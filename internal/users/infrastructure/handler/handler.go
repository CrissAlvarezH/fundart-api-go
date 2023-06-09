package handler

import (
	"errors"
	"github.com/CrissAlvarezH/fundart-api/internal/common"
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/ports"
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
	// users
	e.GET("/api/v1/users", h.List)
	e.GET("/api/v1/users/:id", h.GetByID)
	e.POST("/api/v1/users/", h.Register)
	e.POST("/api/v1/users/login", h.Login)
	e.POST("/api/v1/users/:id/verification-code/", h.ValidateAccountVerificationCode)
	e.PUT("/api/v1/users/:id/", h.Update)
	e.DELETE("/api/v1/users/:id/", h.Delete)
	e.PUT("/api/v1/users/:id/password/", h.ChangePassword)

	// recovery password
	e.POST("/api/v1/users/recovery-password/request", h.RequestRecoveryPassword)
	e.POST("/api/v1/users/recovery-password/", h.RecoveryPassword)

	// addresses
	e.GET("/api/v1/users/:id/addresses", h.ListAddresses)
	e.POST("/api/v1/users/:id/addresses/", h.AddAddress)
	e.PUT("/api/v1/users/:id/addresses/:address_id/", h.UpdateAddress)
	e.DELETE("/api/v1/users/:id/addresses/:address_id/", h.DeleteAddress)
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

	c.JSON(http.StatusOK, MapToRetrieveUserDTO(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var body LoginUserDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"access_token": token.AccessToken, "refresh": token.RefreshToken},
	)
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

	err = h.service.SendAccountVerificationCode(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MapToRetrieveUserDTO(user))
}

func (h *UserHandler) ValidateAccountVerificationCode(c *gin.Context) {
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

	userWasActivated := h.service.ValidateAccountVerificationCode(users.UserID(userID), body.Code)

	if userWasActivated == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"details": "User was activated successfully"})
}

func (h *UserHandler) RequestRecoveryPassword(c *gin.Context) {
	var body RequestRecoveryPasswordDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SendRecoveryPasswordRequest(body.Email)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ports.UserDoesNotExists) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"details": "message sent successfully"})
}

func (h *UserHandler) RecoveryPassword(c *gin.Context) {
	var body RecoveryPasswordDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RecoveryPassword(body.Email, body.NewPassword, body.Code)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ports.InvalidValidationCode) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"details": "recovery password successfully"})
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

	c.JSON(http.StatusOK, MapToRetrieveUserDTO(user))
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var body ChangePasswordDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	err = h.service.ChangePassword(users.UserID(userID), body.CurrentPassword, body.NewPassword)
	if err != nil {
		status := http.StatusInternalServerError
		message := err.Error()
		if errors.Is(err, ports.InvalidCredentials) || errors.Is(err, ports.UserDoesNotExists) {
			status = http.StatusBadRequest
			message = "invalid credentials"
		}
		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"details": "password changed successfully"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	err = h.service.Deactivate(users.UserID(ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) ListAddresses(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	addresses := h.service.ListAddresses(users.UserID(userID))
	addressesDTO := MapToListAddressesDTO(addresses)
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
	addressesDTO := MapToListAddressesDTO(addresses)
	c.JSON(http.StatusCreated, addressesDTO)
}

func (h *UserHandler) UpdateAddress(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	addressID, err := strconv.Atoi(c.Param("address_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address id is not a valid number"})
		return
	}

	var body CreateAddressDTO
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.service.UpdateAddress(
		users.AddressID(addressID), body.Department, body.City,
		body.Address, body.ReceiverPhone, body.ReceiverName,
	)
	if err != nil {
		statusError := http.StatusInternalServerError
		if errors.Is(err, ports.AddressDoesNotExists) {
			statusError = http.StatusNotFound
		}
		c.JSON(statusError, gin.H{"error": err.Error()})
		return
	}

	addresses := h.service.ListAddresses(users.UserID(userID))
	addressesDTO := MapToListAddressesDTO(addresses)
	c.JSON(http.StatusOK, addressesDTO)
}

func (h *UserHandler) DeleteAddress(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is not a valid number"})
		return
	}

	addressID, err := strconv.Atoi(c.Param("address_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address id is not a valid number"})
		return
	}

	err = h.service.DeleteAddress(users.UserID(userID), users.AddressID(addressID))
	if err != nil {
		statusError := http.StatusInternalServerError
		if errors.Is(err, ports.AddressDoesNotExists) {
			statusError = http.StatusNotFound
		}
		c.JSON(statusError, gin.H{"error": err.Error()})
		return
	}

	addresses := h.service.ListAddresses(users.UserID(userID))
	addressesDTO := MapToListAddressesDTO(addresses)
	c.JSON(http.StatusOK, addressesDTO)
}
