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

	userData := []memoryrepo.MemoryUser{
		{
			ID:        1,
			Name:      "Cristian",
			Email:     "cristian@email.com",
			Password:  "23456",
			Phone:     "320684398",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        2,
			Name:      "Yulisa",
			Email:     "yuli@email.com",
			Password:  "dddddd",
			Phone:     "442546536",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        3,
			Name:      "Andrea",
			Email:     "andrea@email.com",
			Password:  "444444",
			Phone:     "320684398",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        4,
			Name:      "Gabriel",
			Email:     "gabriel@email.com",
			Password:  "oooooo",
			Phone:     "53462345",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        5,
			Name:      "Matias",
			Email:     "matias@email.com",
			Password:  "tttttt",
			Phone:     "542534345",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        6,
			Name:      "Manuel",
			Email:     "manuel@email.com",
			Password:  "555555",
			Phone:     "623452345",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        7,
			Name:      "Camilo",
			Email:     "camilo@email.com",
			Password:  "rrrrrr",
			Phone:     "320684398",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        8,
			Name:      "Laura",
			Email:     "laura@email.com",
			Password:  "mmmmmm",
			Phone:     "320684398",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        8,
			Name:      "Laura",
			Email:     "laura@email.com",
			Password:  "22222",
			Phone:     "5345636456",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        9,
			Name:      "Carolina",
			Email:     "carolina@email.com",
			Password:  "22222",
			Phone:     "320684398",
			IsActive:  true,
			CreatedAt: time.Time{},
			Addresses: nil,
			Scopes:    nil,
		},
		{
			ID:        10,
			Name:      "Martin",
			Email:     "martin@email.com",
			Password:  "iriririr",
			Phone:     "42512345234",
			IsActive:  true,
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
