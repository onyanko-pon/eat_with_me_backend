package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/auth"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
)

type DevHandler struct {
	userRepository repository.UserRepository
}

func NewDevHandler(*userRepository repository.UserRepository) (*DevHandler, error) {
	return &DevHandler{
		userRepository: userRepository,
	}, nil
}

func (h DevHandler) Restricted(c echo.Context) error {
	token := c.Get("token").(*jwt.Token)
	fmt.Println(token)
	claims := token.Claims.(*auth.JWTClaim)
	authUser, _ := claims.GenAuthUser()

	return c.String(http.StatusOK, "Welcome "+authUser.UserID)
}

func (h DevHandler) GenToken(c echo.Context) error {

	if os.Getenv("GO_ENV") == "production" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "good bye",
		})
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	user, _ := h.userRepository.GetUser(c.Request().Context(), uint64(id))
	authUser := &auth.AuthUser{
		UserID: idStr,
	}
	token, _ := authUser.GenToken()

	return c.JSON(http.StatusOK, echo.Map{
		"user":  user,
		"token": token,
	})
}
