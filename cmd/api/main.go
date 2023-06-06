package main

import (
	"fmt"
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/services"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/CrissAlvarezH/fundart-api/internal/users/infrastructure"
	"github.com/CrissAlvarezH/fundart-api/internal/users/infrastructure/handler"
	"github.com/CrissAlvarezH/fundart-api/internal/users/infrastructure/memoryrepo"
	"github.com/CrissAlvarezH/fundart-api/internal/users/infrastructure/notifications"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)
import "net/http"

func main() {
	app := gin.Default()

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	app.POST("/upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["images"]
		dst := "./images/"

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			err := c.SaveUploadedFile(file, dst)
			if err != nil {
				log.Println("error on save image ", file.Filename)
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	userData := []users.User{
		{
			ID:        1,
			Name:      "Cristian",
			Email:     "cristian@email.com",
			Phone:     "320684398",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        2,
			Name:      "Yulisa",
			Email:     "yuli@email.com",
			Phone:     "442546536",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        3,
			Name:      "Andrea",
			Email:     "andrea@email.com",
			Phone:     "320684398",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        4,
			Name:      "Gabriel",
			Email:     "gabriel@email.com",
			Phone:     "53462345",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        5,
			Name:      "Matias",
			Email:     "matias@email.com",
			Phone:     "542534345",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        6,
			Name:      "Manuel",
			Email:     "manuel@email.com",
			Phone:     "623452345",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        7,
			Name:      "Camilo",
			Email:     "camilo@email.com",
			Phone:     "320684398",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        8,
			Name:      "Laura",
			Email:     "laura@email.com",
			Phone:     "5345636456",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        9,
			Name:      "Carolina",
			Email:     "carolina@email.com",
			Phone:     "320684398",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        10,
			Name:      "Martin",
			Email:     "martin@email.com",
			Phone:     "42512345234",
			IsActive:  false,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
	}
	userMemoRepo := memoryrepo.NewMemoryUserRepository(userData)

	addressMemoRepo := memoryrepo.NewMemoryAddressRepository(make([]users.Address, 0))
	mockPassManager := infrastructure.NewMockPasswordManager()
	mockVerifyCode := notifications.NewMockVerificationCodeManager()

	userService := services.NewUserService(userMemoRepo, addressMemoRepo, mockVerifyCode, mockPassManager)
	userHandler := handler.NewUserHandler(userService)

	userHandler.AddRoutes(app)

	err := app.Run(":8000")
	if err != nil {
		log.Fatal("Server error", err)
	}
}
