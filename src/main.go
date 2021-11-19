package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/handler"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	userHandler, _ := handler.NewUserHandler()

	e.GET("/api/users/:id/friends", userHandler.GetFriends)

	fmt.Println(os.Environ())

	e.Logger.Fatal(
		e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))),
	)
}
