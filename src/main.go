package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/onyanko-pon/eat_with_me_backend/src/auth"
	"github.com/onyanko-pon/eat_with_me_backend/src/handler"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
	"github.com/onyanko-pon/eat_with_me_backend/src/sql_handler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	e := echo.New()

	var dataSource string
	if os.Getenv("GO_ENV") == "production" {
		dataSource = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
		fmt.Println(dataSource)
	} else {
		dataSource = "host=127.0.0.1 port=5432 user=admin password=password dbname=mydb sslmode=disable"
	}
	sqlHandler, err := sql_handler.NewHandler(dataSource)

	if err != nil {
		fmt.Printf("connect error: %s\n", err.Error())
		panic(1)
	}

	config := middleware.JWTConfig{
		ContextKey:    "token",
		SigningMethod: "HS256",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		Claims:        &auth.JWTClaim{},
		SigningKey:    []byte(os.Getenv("JWT_SIGNINGKEY")),
	}

	jwtMiddleware := middleware.JWTWithConfig(config)

	userRepository := repository.NewUserRepository(sqlHandler)
	eventRepository := repository.NewEventRepository(sqlHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	userHandler, _ := handler.NewUserHandler(userRepository, eventRepository)
	eventHandler, _ := handler.NewEventHandler(eventRepository)

	e.GET("/api/users/:id/friends", userHandler.GetFriends, jwtMiddleware)
	e.POST("/api/users", userHandler.CreateUser)
	e.PUT("/api/users", userHandler.UpdateUser, jwtMiddleware)
	e.GET("/api/users/:id", userHandler.GetUser, jwtMiddleware)

	e.POST("/api/events", eventHandler.CreateEvent, jwtMiddleware)
	e.PUT("/api/events", eventHandler.UpdateEvent, jwtMiddleware)
	e.GET("/api/events/:id", eventHandler.GetEvent, jwtMiddleware)

	e.GET("/api/users/:id/events", userHandler.GetEvents, jwtMiddleware)
	e.GET("/api/users/:id/events/joining", eventHandler.GetJoiningEvents, jwtMiddleware)
	e.POST("/api/events/:id/join", eventHandler.JoinEvent, jwtMiddleware)

	e.POST("/api/users/:id/usericons", userHandler.UploadUserIcon, jwtMiddleware)

	e.GET("/api/restricted", userHandler.Restricted, jwtMiddleware)
	e.GET("/api/users/:id/token", userHandler.GenToken)

	e.POST("/api/users/:id/apply/:friend_user_id", userHandler.ApplyFriend, jwtMiddleware)

	e.Logger.Fatal(
		e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))),
	)
}
