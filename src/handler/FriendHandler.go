package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
)

type FriendHandler struct {
	userRepository *repository.UserRepository
}

func NewFriendHandler(userRepository *repository.UserRepository) (*FriendHandler, error) {
	return &FriendHandler{
		userRepository: userRepository,
	}, nil
}

func (h FriendHandler) GetFriends(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friends, err := h.userRepository.GetFriends(c.Request().Context(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(friends) == 0 {
		friends = make([]entity.Friend, 0)
	}

	requestFriends, err := h.userRepository.GetRequestFriends(c.Request().Context(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(requestFriends) == 0 {
		requestFriends = make([]entity.Friend, 0)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"friends":         friends,
		"request_friends": requestFriends,
	})
}

func (h FriendHandler) BlockFriend(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friendUserIDStr := c.Param("friend_user_id")
	friendUserID, _ := strconv.Atoi(friendUserIDStr)

	err := h.userRepository.BlockFriend(context.Background(), uint64(id), uint64(friendUserID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (h FriendHandler) ApplyFriend(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friendUserIDStr := c.Param("friend_user_id")
	friendUserID, _ := strconv.Atoi(friendUserIDStr)

	err := h.userRepository.ApplyFriend(context.Background(), uint64(id), uint64(friendUserID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (h FriendHandler) AcceptApplyFriend(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friendUserIDStr := c.Param("friend_user_id")
	friendUserID, _ := strconv.Atoi(friendUserIDStr)

	err := h.userRepository.AcceptApplyFriend(context.Background(), uint64(id), uint64(friendUserID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (h FriendHandler) GetRecommendUsers(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	users, err := h.userRepository.GetRecommendUsers(c.Request().Context(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(users) == 0 {
		users = make([]entity.User, 0)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"users": users,
	})
}
