package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/handler"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
	"github.com/onyanko-pon/eat_with_me_backend/src/sql_handler"
)

func main() {
	e := echo.New()

	dataSource := "host=127.0.0.1 port=5432 user=admin password=password dbname=mydb sslmode=disable"
	sqlHandler, err := sql_handler.NewHandler(dataSource)

	if err != nil {
		fmt.Printf("connect error: %s\n", err.Error())
		panic(1)
	}

	userRepository := repository.NewUserRepository(sqlHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	userHandler, _ := handler.NewUserHandler(userRepository)

	e.GET("/api/users/:id/friends", userHandler.GetFriends)
	e.POST("/api/users", userHandler.CreateUser)
	e.PUT("/api/users", userHandler.UpdateUser)
	e.GET("/api/users/:id", userHandler.GetUser)

	eventRepository := repository.NewEventRepository(sqlHandler)
	eventHandler, _ := handler.NewEventHandler(eventRepository)

	e.POST("/api/events", eventHandler.CreateEvent)
	e.PUT("/api/events", eventHandler.UpdateEvent)
	e.GET("/api/events/:id", eventHandler.GetEvent)

	e.Logger.Fatal(
		e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))),
	)
}
