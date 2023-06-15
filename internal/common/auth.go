package common

import (
	"github.com/CrissAlvarezH/fundart-api/internal/users/application/ports"
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth(jwt ports.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}
		authHeaderSplit := strings.Split(authHeader, " ")
		if len(authHeaderSplit) < 2 {
			c.Next()
			return
		}

		token := authHeaderSplit[1]

		user, err := jwt.Verify(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user", user)
	}
}

type ScopeValidatorFunc func(user users.User, isAnonymous bool, c *gin.Context) (bool, string)

func extractUser(c *gin.Context) (users.User, bool) {
	user, ok := c.Get("user")
	if ok {
		return user.(users.User), true
	} else {
		return users.User{}, false
	}
}

func ValidOr(functions ...ScopeValidatorFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := extractUser(c)
		validationMsg := ""
		for _, f := range functions {
			result, msg := f(user, !ok, c)
			validationMsg = msg
			if result {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": validationMsg})
		c.Abort()
	}
}

func Valid(function ScopeValidatorFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := extractUser(c)
		result, msg := function(user, !ok, c)

		if !result {
			c.JSON(http.StatusForbidden, gin.H{"error": msg})
			c.Abort()
		}

		c.Next()
	}
}
