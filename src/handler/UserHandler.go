package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/auth"
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

func (u UserHandler) FetchUserByUsername(c echo.Context) error {

	username := c.Param("username")

	user, _ := u.UserRepository.FetchUserByUsername(c.Request().Context(), username)

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

type responseCreateUser struct {
	User  *entity.User `json:"user"`
	Token string       `json:"token"`
}

type requestBodyCreateUser struct {
	User *entity.User `json:"user"`
}

func (u UserHandler) CreateUser(c echo.Context) error {

	requestBody := new(requestBodyCreateUser)
	if err := c.Bind(requestBody); err != nil {
		return err
	}

	user, err := u.UserRepository.CreateUser(c.Request().Context(), *requestBody.User)

	if err != nil {
		return err
	}

	authUser := &auth.AuthUser{
		UserID: strconv.Itoa(int(user.ID)),
	}
	token, _ := authUser.GenToken()

	return c.JSON(http.StatusOK, responseCreateUser{
		User:  user,
		Token: token,
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

	params, _ := c.FormParams()
	fmt.Println("form params", params)

	if err != nil {
		fmt.Println(err)
		return err
	}

	data, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return err
	}

	extension := filepath.Ext(file.Filename)
	filename := "username" + userIDStr + extension

	d, err := image.Resize(data, filename)
	if err != nil {
		fmt.Println(err)
		return err
	}

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

func (u UserHandler) Restricted(c echo.Context) error {
	token := c.Get("token").(*jwt.Token)
	fmt.Println(token)
	claims := token.Claims.(*auth.JWTClaim)
	authUser, _ := claims.GenAuthUser()

	return c.String(http.StatusOK, "Welcome "+authUser.UserID)
}

func (u UserHandler) GenToken(c echo.Context) error {

	if os.Getenv("GO_ENV") == "production" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "good bye",
		})
	}

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	user, _ := u.UserRepository.GetUser(c.Request().Context(), uint64(id))
	authUser := &auth.AuthUser{
		UserID: idStr,
	}
	token, _ := authUser.GenToken()

	return c.JSON(http.StatusOK, echo.Map{
		"user":  user,
		"token": token,
	})
}

func (u UserHandler) ApplyFriend(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	friendUserIDStr := c.Param("friend_user_id")
	friendUserID, _ := strconv.Atoi(friendUserIDStr)

	err := u.UserRepository.ApplyFriend(context.Background(), uint64(id), uint64(friendUserID))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.NoContent(http.StatusCreated)
}
