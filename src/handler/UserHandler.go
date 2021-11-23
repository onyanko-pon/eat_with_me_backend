package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/image"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
	"github.com/onyanko-pon/eat_with_me_backend/src/service"
)

type UserHandler struct {
	UserRepository  *repository.UserRepository
	EventRepository *repository.EventRepository
	FileService     *service.FileService
}

func NewUserHandler(userRepository *repository.UserRepository, eventRepository *repository.EventRepository) (*UserHandler, error) {
	fileService, _ := service.NewFileService()
	return &UserHandler{
		UserRepository:  userRepository,
		EventRepository: eventRepository,
		FileService:     fileService,
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

	_, err := u.UserRepository.CreateUser(c.Request().Context(), *user)

	if err != nil {
		return err
	}
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

	_, err := u.UserRepository.UpdateUser(c.Request().Context(), *user)

	if err != nil {
		return err
	}

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

	friends, err := u.UserRepository.GetFriends(c.Request().Context(), uint64(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responseGetFriends{
		Friends: friends,
	})
}

type responseGetEvents struct {
	User   *entity.User   `json:"user"`
	Events []entity.Event `json:"events"`
}

func (h UserHandler) GetEvents(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	user, err := h.UserRepository.GetUser(c.Request().Context(), uint64(id))

	if err != nil {
		fmt.Print(err)
		return err
	}

	events, err := h.EventRepository.GetEventsRelatedToUser(c.Request().Context(), *user)

	if err != nil {
		fmt.Print(err)
		return err
	}

	return c.JSON(http.StatusOK, responseGetEvents{
		Events: events,
		User:   user,
	})
}

func (h UserHandler) UploadUserIcon(c echo.Context) error {

	userIDStr := c.Param("id")
	userID, _ := strconv.Atoi(userIDStr)

	file, err := c.FormFile("usericon")
	if err != nil {
		return err
	}

	data, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return err
	}

	extension := filepath.Ext(file.Filename)
	filename := "username" + userIDStr + extension

	d, _ := image.Resize(data, filename)

	url, err := h.FileService.UploadUserIcon(d, filename)
	if err != nil {
		fmt.Println(err)
		return err
	}

	user, err := h.UserRepository.GetUser(c.Request().Context(), uint64(userID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	user.ImageURL = url
	user, err = h.UserRepository.UpdateUser(c.Request().Context(), *user)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"filename": file.Filename,
		"url":      url,
	})
}
