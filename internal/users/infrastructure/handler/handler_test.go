package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/services"
	"github.com/CrissAlvarezH/fundart-api/internal/users/infrastructure"
	"github.com/CrissAlvarezH/fundart-api/internal/users/infrastructure/handler"
	"github.com/CrissAlvarezH/fundart-api/internal/users/infrastructure/memoryrepo"
	"github.com/CrissAlvarezH/fundart-api/internal/users/infrastructure/notifications"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func createMockUserService(
	userData []memoryrepo.MemoryUser, addressData []memoryrepo.MemoryAddress,
) (services.UserService, notifications.MockVerificationCodeManager) {
	userMemoRepo := memoryrepo.NewMemoryUserRepository(userData)

	addressMemoRepo := memoryrepo.NewMemoryAddressRepository(addressData)
	mockPassManager := infrastructure.NewMockPasswordManager()
	mockVerifyCode := notifications.NewMockVerificationCodeManager()
	mockJWTManager := infrastructure.NewMockJWTManager(userMemoRepo)

	userService := services.NewUserService(
		userMemoRepo, addressMemoRepo, mockVerifyCode,
		mockPassManager, mockJWTManager,
	)
	return userService, *mockVerifyCode
}

func TestUserHandler_List(t *testing.T) {
	userData := []memoryrepo.MemoryUser{
		{
			ID:        1,
			Name:      "Cristian",
			Email:     "cristian@email.com",
			Password:  "23456_encrypt",
			Phone:     "320684398",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        2,
			Name:      "Yuli",
			Email:     "yuli@email.com",
			Password:  "ddd_encrypt",
			Phone:     "442546536",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        3,
			Name:      "Juan",
			Email:     "juan@email.com",
			Password:  "ccc_encrypt",
			Phone:     "44444444",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
	}
	userService, _ := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	apiV1Routes := router.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	// TESTING PAGINATION DATA
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("status code != 200, code:", w.Code, "response:", w.Body.String())
	}

	resBody := make(map[string]interface{})
	if err := json.Unmarshal([]byte(w.Body.String()), &resBody); err != nil {
		t.Error("error parsing to json:", w.Body.String())
	}

	totalUsers := resBody["pagination"].(map[string]interface{})["total"].(float64)
	if totalUsers != float64(len(userData)) {
		t.Error("incorrect total users:", totalUsers)
	}

	// test total pages, must be 1 because default page size is 10
	totalPages := resBody["pagination"].(map[string]interface{})["total_pages"].(float64)
	if totalPages != 1 {
		t.Error("incorrect total pages:", totalPages)
	}

	bodyResult := resBody["result"].([]interface{})
	if len(bodyResult) != len(userData) {
		t.Error("incorrect len of result data", len(bodyResult))
	}

	// TEST GO THROUGH PAGINATION FOR PAGES
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users?page=2&page_size=2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// must there is just one item because page size is 2, total items are 3, and we are in second page
	resBody = make(map[string]interface{})
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("error parsing to json:", w.Body.String())
	}

	bodyResult = resBody["result"].([]interface{})
	if len(bodyResult) != 1 {
		t.Error("incorrect len of users in second page with page size 2:", len(bodyResult))
	}
	currentPage := resBody["pagination"].(map[string]interface{})["page"].(float64)
	if currentPage != 2 {
		t.Error("incorrect current page, must be 2, is:", totalUsers)
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	cristianUser := memoryrepo.MemoryUser{
		ID:        1,
		Name:      "Cristian",
		Email:     "cristian@email.com",
		Password:  "23456_encrypt",
		Phone:     "320684398",
		IsActive:  true,
		CreatedAt: time.Time{},
		Addresses: nil,
		Scopes:    nil,
	}
	userData := []memoryrepo.MemoryUser{cristianUser}
	userService, _ := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	apiV1Routes := router.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	// TESTING HAPPY PATH
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+strconv.Itoa(int(cristianUser.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("status code invalid:", w.Code)
	}

	resBody := make(map[string]interface{})
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("error parsing to json:", w.Body.String())
	}

	if resBody["id"] != float64(cristianUser.ID) {
		t.Error("incorrect user id:", resBody["id"], "expected:", cristianUser.ID)
	}
	if resBody["name"] != cristianUser.Name {
		t.Error("incorrect user name:", resBody["name"], "expected:", cristianUser.Name)
	}

	// TESTING UNHAPPY PATH
	// this user does not exist
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/users/10", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("incorrect status code:", w.Code, "expected not found")
	}
}

func TestUserHandler_Login(t *testing.T) {
	cristianUser := memoryrepo.MemoryUser{
		ID:        1,
		Name:      "Cristian",
		Email:     "cristian@email.com",
		Password:  "23456_encrypt", // the mock password manager add _encrypt to the passwords
		Phone:     "320684398",
		IsActive:  true,
		CreatedAt: time.Time{},
		Addresses: nil,
		Scopes:    nil,
	}
	userData := []memoryrepo.MemoryUser{cristianUser}
	userService, _ := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	apiV1Routes := router.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	// TEST HAPPY PATH
	reqBody := bytes.NewReader([]byte(`{"email": "cristian@email.com", "password": "23456"}`))
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/login", reqBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("login status code:", w.Code, "expected", http.StatusOK)
		return
	}

	resBody := make(map[string]string)
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("error parsing json login response:", err)
	}

	if _, ok := resBody["access_token"]; ok == false {
		t.Error("login does not return 'access_token', body:", w.Body.String())
	}
	if _, ok := resBody["refresh"]; ok == false {
		t.Error("login does not return 'refresh', body:", w.Body.String())
	}

	// TEST UNHAPPY PATH
	reqBody = bytes.NewReader([]byte(`{"email": "cristian@email.com", "password": "23456__"}`))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/login", reqBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("login status code:", w.Code, "expected", http.StatusBadRequest)
		return
	}
}
