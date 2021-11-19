package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
)

type UserHandler struct {
	UserRepository *repository.UserRepository
}

func NewUserHandler(userRepository *repository.UserRepository) (*UserHandler, error) {
	return &UserHandler{
		UserRepository: userRepository,
	}, nil
}

type responseGetUser struct {
	User *entity.User `json:"user"`
}

func (u UserHandler) GetUser(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	user, _ := u.UserRepository.GetUser(c.Request().Context(), uint64(id))

	return c.JSON(http.StatusOK, responseGetUser{
		User: user,
	})
}

type responseCreateUser struct {
	User *entity.User `json:"user"`
}

func (u UserHandler) CreateUser(c echo.Context) error {

	user := new(entity.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	u.UserRepository.CreateUser(c.Request().Context(), *user)

	return c.JSON(http.StatusOK, responseCreateUser{
		User: user,
	})
}

type responseUpdateUser struct {
	User *entity.User `json:"user"`
}

func (u UserHandler) UpdateUser(c echo.Context) error {

	user := new(entity.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	u.UserRepository.UpdateUser(c.Request().Context(), *user)

	return c.JSON(http.StatusOK, responseUpdateUser{
		User: user,
	})
}

type responseGetFriends struct {
	Friends []entity.User `json:"friends"`
}

func (u UserHandler) GetFriends(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friends, _ := u.UserRepository.GetFriends(c.Request().Context(), uint64(id))

	return c.JSON(http.StatusOK, responseGetFriends{
		Friends: friends,
	})
}
