package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
)

type UserHandler struct{}

func NewUserHandler() (*UserHandler, error) {
	return &UserHandler{}, nil
}

type responseFriends struct {
	Friends []entity.User `json:"friends"`
}

func (u UserHandler) GetFriends(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var friends []entity.User
	var friend entity.User
	friend.ID = uint64(id)
	friend.Username = "username"
	friend.ImageURL = "https://avatars.githubusercontent.com/u/54364185?v=4"

	friends = append(friends, friend)

	return c.JSON(http.StatusOK, responseFriends{
		Friends: friends,
	})
}
