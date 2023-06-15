package handler

import (
	users "github.com/CrissAlvarezH/fundart-api/internal/users/domain"
	"github.com/gin-gonic/gin"
	"strconv"
)

func IsSameUser(user users.User, isAnonymous bool, c *gin.Context) (bool, string) {
	if isAnonymous {
		return false, "user is anonymous"
	}

	// get user id from path param
	id := c.Param("id")

	if strconv.Itoa(int(user.ID)) != id {
		return false, "is not the same user"
	}

	return true, ""
}

func ScopeUserRead(user users.User, isAnonymous bool, c *gin.Context) (bool, string) {
	if isAnonymous {
		return false, "user is anonymous"
	}

	has := user.HasScope(users.USERS_READ)
	if !has {
		return false, "doesn't have permissions"
	}

	return true, ""
}