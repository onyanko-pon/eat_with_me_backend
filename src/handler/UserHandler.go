package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/auth"
	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/image"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
	"github.com/onyanko-pon/eat_with_me_backend/src/service"
	"github.com/onyanko-pon/eat_with_me_backend/src/usecase"
)

type UserHandler struct {
	UserRepository    *repository.UserRepository
	EventRepository   *repository.EventRepository
	FileService       *service.FileService
	createUserUsecase *usecase.CreateUserUsecase
}

func NewUserHandler(userRepository *repository.UserRepository, eventRepository *repository.EventRepository, createUserUsecase *usecase.CreateUserUsecase) (*UserHandler, error) {
	fileService, _ := service.NewFileService()
	return &UserHandler{
		UserRepository:    userRepository,
		EventRepository:   eventRepository,
		FileService:       fileService,
		createUserUsecase: createUserUsecase,
	}, nil
}

func (u UserHandler) GetUser(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	user, _ := u.UserRepository.GetUser(c.Request().Context(), uint64(id))

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

func (u UserHandler) FetchUserByUsername(c echo.Context) error {

	username := c.Param("username")

	user, _ := u.UserRepository.FetchUserByUsername(c.Request().Context(), username)

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
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

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
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

	if len(events) == 0 {
		events = make([]entity.Event, 0)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"events": events,
		"user":   user,
	})
}

func (h UserHandler) UploadUserIcon(c echo.Context) error {

	userIDStr := c.Param("id")
	userID, _ := strconv.Atoi(userIDStr)

	file, err := c.FormFile("usericon")

	form, _ := c.MultipartForm()
	fmt.Println("form files", form.File)

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

type reqCreateUserWithTwitterVerify struct {
	OAuthToken    string `json:"oauth_token"`
	OAuthVerifier string `json:"oauth_verifier"`
	OAuthSecret   string `json:"oauth_secret"`
}

func (h UserHandler) CreateUserWithTwitterVerify(c echo.Context) error {

	requestBody := new(reqCreateUserWithTwitterVerify)
	if err := c.Bind(requestBody); err != nil {
		return err
	}

	user, err := h.createUserUsecase.CreateUserWithTwitterVerify(c.Request().Context(), requestBody.OAuthToken, requestBody.OAuthSecret, requestBody.OAuthVerifier)
	if err != nil {
		return err
	}
	authUser := &auth.AuthUser{
		UserID: strconv.Itoa(int(user.ID)),
	}
	jwtToken, _ := authUser.GenToken()

	return c.JSON(http.StatusOK, echo.Map{
		"user":  *user,
		"token": jwtToken,
	})
}

type reqCreateUserWithAppleVerify struct {
	UserIdentifier string `json:"user_identifier"`
}

func (h UserHandler) CreateUserWithAppleVerify(c echo.Context) error {

	requestBody := new(reqCreateUserWithAppleVerify)
	if err := c.Bind(requestBody); err != nil {
		return err
	}

	user, err := h.createUserUsecase.CreateUserWithAppleVerify(c.Request().Context(), requestBody.UserIdentifier)
	if err != nil {
		return err
	}
	authUser := &auth.AuthUser{
		UserID: strconv.Itoa(int(user.ID)),
	}
	jwtToken, _ := authUser.GenToken()

	return c.JSON(http.StatusOK, echo.Map{
		"user":  *user,
		"token": jwtToken,
	})
}
