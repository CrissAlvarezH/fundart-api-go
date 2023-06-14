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
) (
	services.UserService, notifications.MockVerificationCodeManager,
	memoryrepo.MemoryUserRepository,
) {
	userMemoRepo := memoryrepo.NewMemoryUserRepository(userData)

	addressMemoRepo := memoryrepo.NewMemoryAddressRepository(addressData)
	mockPassManager := infrastructure.NewMockPasswordManager()
	mockVerifyCode := notifications.NewMockVerificationCodeManager()
	mockJWTManager := infrastructure.NewMockJWTManager(userMemoRepo)

	userService := services.NewUserService(
		userMemoRepo, addressMemoRepo, mockVerifyCode,
		mockPassManager, mockJWTManager,
	)
	return userService, *mockVerifyCode, *userMemoRepo
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
	userService, _, _ := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
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
	userService, _, _ := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
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
	userService, _, _ := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
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

func TestUserHandler_Register_And_VerifyAccount(t *testing.T) {
	service, verifyCodeManager, _ := createMockUserService(
		make([]memoryrepo.MemoryUser, 0), make([]memoryrepo.MemoryAddress, 0),
	)
	userHandler := handler.NewUserHandler(service)

	routes := gin.New()
	apiV1Routes := routes.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	reqBodyRaw := `
		{ 
			"name": "Juan",
			"email": "juan@email.com",
			"phone": "3207846634",
			"password": "333333"
		}
	`
	reqBody := bytes.NewReader([]byte(reqBodyRaw))
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", reqBody)
	w := httptest.NewRecorder()
	routes.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("register user response code:", w.Code, "expected:", http.StatusCreated)
		t.Log("register body res:", w.Body.String())
		return
	}

	resBody := make(map[string]interface{})
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("error parsing register user body:", err)
		return
	}

	if resBody["email"].(string) != "juan@email.com" || resBody["phone"].(string) != "3207846634" {
		t.Error(
			"register user return email and phone incorrect:", w.Body.String(),
			"expected:", "juan@email.com, 3207846634",
		)
	}

	verifyCode, ok := verifyCodeManager.AccountCodes["juan@email.com"]

	if ok == false {
		t.Error("verification account code was not register")
		return
	}

	reqBody = bytes.NewReader([]byte(`{"code": "` + verifyCode + `"}`))
	userID := strconv.Itoa(int(resBody["id"].(float64)))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/"+userID+"/verification-code/", reqBody)
	w = httptest.NewRecorder()
	routes.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("account verification status code was:", w.Code, "expected:", http.StatusOK)
		return
	}
}

func TestUserHandler_Request_And_RecoveryPassword(t *testing.T) {
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
	userService, verifyCodeManager, _ := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	apiV1Routes := router.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	// REQUEST RECOVERY PASSWORD
	reqBody := bytes.NewReader([]byte(`{"email": "` + cristianUser.Email + `"}`))
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/recovery-password/request", reqBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("request recovery pasword status code:", w.Code, "expected:", http.StatusOK)
		t.Log("response body", w.Body.String())
		return
	}

	// USE VALIDATION CODE SENT TO CHANGE PASSWORD
	codeSent := verifyCodeManager.PassCodes[cristianUser.Email]
	newPass := "444444"
	reqBody = bytes.NewReader([]byte(`{
		"email": "` + cristianUser.Email + `",
		"new_password": "` + newPass + `",
		"code": "` + codeSent + `"
	}`))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/recovery-password/", reqBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("recovery password response code:", w.Code, "expected:", http.StatusOK)
		t.Log("response body:", w.Body.String())
		return
	}

	// TEST NEW PASSWORD WITH LOGIN
	reqBody = bytes.NewReader([]byte(`{
		"email": "` + cristianUser.Email + `",
		"password": "` + newPass + `"
	}`))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/login", reqBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("login after reset pass response code:", w.Code, "expected:", http.StatusOK)
		t.Log("login res body:", w.Body.String())
	}

	// TEST PREVIOUS PASSWORD WITH LOGIN
	reqBody = bytes.NewReader([]byte(`{
		"email": "` + cristianUser.Email + `",
		"password": "23456"
	}`))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/login", reqBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("login with previous pass response code:", w.Code, "expected:", http.StatusBadRequest)
		t.Log("login res body:", w.Body.String())
	}
}

func TestUserHandler_Update(t *testing.T) {
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
	userService, _, userRepo := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	apiV1Routes := router.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	// UPDATE USER
	newPhone := "3207774343"
	newEmail := "alvarez@email.com"
	reqBody := bytes.NewReader([]byte(`{
		"name": "Cristian Alvarez",
		"email": "` + newEmail + `",
		"phone": "` + newPhone + `",
		"scopes": []
	}`))
	req, _ := http.NewRequest(
		http.MethodPut, "/api/v1/users/"+strconv.Itoa(int(cristianUser.ID))+"/",
		reqBody,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("update user res status code:", w.Code, "expected:", http.StatusOK)
		t.Log("update user res body:", w.Body.String())
		return
	}

	// VALIDATE CHANGES IN USER REPOSITORY
	userInRepo := userRepo.Users[0]
	if userInRepo.Email != newEmail {
		t.Error("email updated incorrect:", userInRepo.Email, "expected:", newEmail)
	}
	if userInRepo.Phone != newPhone {
		t.Error("phone updated incorrect:", userInRepo.Phone, "expected:", newPhone)
	}
}

func TestUserHandler_ChangePassword(t *testing.T) {
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
	userService, _, _ := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	apiV1Routes := router.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	// UPDATE USER
	newPass := "888888"
	reqBody := bytes.NewReader([]byte(`{
		"current_password": "23456",
		"new_password": "` + newPass + `"
	}`))
	req, _ := http.NewRequest(
		http.MethodPut, "/api/v1/users/"+strconv.Itoa(int(cristianUser.ID))+"/password/",
		reqBody,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("update pass res status code:", w.Code, "expected:", http.StatusOK)
		t.Log("update pass res body:", w.Body.String())
		return
	}

	// TEST NEW PASSWORD WITH LOGIN
	reqBody = bytes.NewReader([]byte(`{
		"email": "` + cristianUser.Email + `",
		"password": "` + newPass + `"
	}`))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/login", reqBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("login after change pass response code:", w.Code, "expected:", http.StatusOK)
		t.Log("login res body:", w.Body.String())
	}

	// TEST PREVIOUS PASSWORD WITH LOGIN
	reqBody = bytes.NewReader([]byte(`{
		"email": "` + cristianUser.Email + `",
		"password": "23456"
	}`))
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/users/login", reqBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("login previous pass response code:", w.Code, "expected:", http.StatusBadRequest)
		t.Log("login res body:", w.Body.String())
	}
}

func TestUserHandler_Delete(t *testing.T) {
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
	userService, _, userRepo := createMockUserService(userData, make([]memoryrepo.MemoryAddress, 0))
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	apiV1Routes := router.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/"+strconv.Itoa(int(cristianUser.ID))+"/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Error("delete user code:", w.Code, "expected:", http.StatusNoContent)
		return
	}

	userInRepo := userRepo.Users[0]
	if userInRepo.IsActive == true {
		t.Error("user was not deleted")
		t.Log("users in repo:", userRepo.Users)
	}
}

func TestUserHandler_List_Add_Update_And_DeleteAddress(t *testing.T) {
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
	addresses := make([]memoryrepo.MemoryAddress, 0)
	userData := []memoryrepo.MemoryUser{cristianUser}
	userService, _, _ := createMockUserService(userData, addresses)
	userHandler := handler.NewUserHandler(userService)

	router := gin.New()
	apiV1Routes := router.Group("/api/v1")
	userHandler.AddRoutes(apiV1Routes)

	// CHECK THERE AREN'T ANY ADDRESS
	if len(userService.ListAddresses(cristianUser.ID)) != 0 {
		t.Error("address must be empty")
		t.Log("addresses:", userService.ListAddresses(cristianUser.ID))
		return
	}

	// ADD ADDRESS
	reqBody := bytes.NewReader([]byte(`{
		"department": "cordoba",
		"city": "planeta rica",
		"address": "calle 24 crr 45",
		"receiver_name": "Cristian Alvarez",
		"receiver_phone": "3205467456"
	}`))
	req, _ := http.NewRequest(
		http.MethodPost, "/api/v1/users/"+strconv.Itoa(int(cristianUser.ID))+"/addresses/", reqBody,
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("status code created address:", w.Code, "expected:", http.StatusCreated)
		t.Log("body:", w.Body.String())
		return
	}
	firstAddress := userService.ListAddresses(cristianUser.ID)[0]
	if len(userService.ListAddresses(cristianUser.ID)) != 1 {
		t.Error("address must one address")
		t.Log("addresses:", userService.ListAddresses(cristianUser.ID))
		return
	}

	// UPDATE ADDRESS
	newDep := "cundinamarca updated"
	newAddress := "cr 3 cll 4 updated"
	reqBody = bytes.NewReader([]byte(`{
		"department": "` + newDep + `",
		"city": "bogota",
		"address": "` + newAddress + `",
		"receiver_name": "Cristian Alvarez",
		"receiver_phone": "3497775888"
	}`))
	firstAddressID := strconv.Itoa(int(firstAddress.ID))
	req, _ = http.NewRequest(
		http.MethodPut,
		"/api/v1/users/"+strconv.Itoa(int(cristianUser.ID))+"/addresses/"+firstAddressID+"/",
		reqBody,
	)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("status code created address:", w.Code, "expected:", http.StatusCreated)
		t.Log("body:", w.Body.String())
		return
	}
	firstAddress = userService.ListAddresses(cristianUser.ID)[0]
	if firstAddress.Department != newDep || firstAddress.Address != newAddress {
		t.Error("address was not updated, id:", firstAddressID)
		t.Log("res body:", w.Body.String())
		return
	}

	// ADD A SECOND ADDRESS
	reqBody = bytes.NewReader([]byte(`{
		"department": "cundinamarca",
		"city": "bogota",
		"address": "cr 2 cll 4",
		"receiver_name": "Cristian Alvarez",
		"receiver_phone": "3497775888"
	}`))
	req, _ = http.NewRequest(
		http.MethodPost, "/api/v1/users/"+strconv.Itoa(int(cristianUser.ID))+"/addresses/", reqBody,
	)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("status code created address:", w.Code, "expected:", http.StatusCreated)
		t.Log("body:", w.Body.String())
		return
	}
	if len(userService.ListAddresses(cristianUser.ID)) != 2 {
		t.Error("address must one address")
		t.Log("addresses:", userService.ListAddresses(cristianUser.ID))
		return
	}

	// LIST ADDRESS (2 address on repo)
	req, _ = http.NewRequest(
		http.MethodGet, "/api/v1/users/"+strconv.Itoa(int(cristianUser.ID))+"/addresses", nil,
	)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("list addresses status code:", w.Code, "expected:", http.StatusOK)
		t.Log("response:", w.Body.String())
		return
	}

	resBody := make([]interface{}, 0)
	if err := json.Unmarshal(w.Body.Bytes(), &resBody); err != nil {
		t.Error("error to parsing address json:", err, "body:")
		t.Log("body", w.Body.String())
		return
	}

	if len(resBody) != 2 {
		t.Error("incorrect amount of address:", len(resBody), "expected 2")
		return
	}
}
