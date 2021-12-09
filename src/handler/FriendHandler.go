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
	userRepository   *repository.UserRepository
	friendRepository *repository.FriendRepository
}

func NewFriendHandler(
	userRepository *repository.UserRepository,
	friendRepository *repository.FriendRepository,
) (*FriendHandler, error) {
	return &FriendHandler{
		userRepository:   userRepository,
		friendRepository: friendRepository,
	}, nil
}

func (h FriendHandler) GetFriends(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friends, err := h.friendRepository.GetFriends(c.Request().Context(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return err
	}

	applyings, err := h.friendRepository.GetApplyings(c.Request().Context(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return err
	}

	applyeds, err := h.friendRepository.GetApplieds(c.Request().Context(), uint64(id))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"friends":   friends,
		"applyings": applyings,
		"applyeds":  applyeds,
	})
}

func (h FriendHandler) Blind(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friendUserIDStr := c.Param("friend_user_id")
	friendUserID, _ := strconv.Atoi(friendUserIDStr)

	err := h.friendRepository.Blind(context.Background(), uint64(id), uint64(friendUserID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (h FriendHandler) Apply(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friendUserIDStr := c.Param("friend_user_id")
	friendUserID, _ := strconv.Atoi(friendUserIDStr)

	err := h.friendRepository.Apply(context.Background(), uint64(id), uint64(friendUserID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (h FriendHandler) AcceptApply(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friendUserIDStr := c.Param("friend_user_id")
	friendUserID, _ := strconv.Atoi(friendUserIDStr)

	err := h.friendRepository.AcceptApply(context.Background(), uint64(id), uint64(friendUserID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (h FriendHandler) GetRecommendUsers(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	users, err := h.friendRepository.GetRecommendUsers(c.Request().Context(), uint64(id))
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
