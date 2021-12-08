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
	"github.com/onyanko-pon/eat_with_me_backend/src/service"
	"github.com/onyanko-pon/eat_with_me_backend/src/sql_handler"
	"github.com/onyanko-pon/eat_with_me_backend/src/usecase"
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

	twitterAuthService := &service.TwitterAuthService{}
	userService, _ := service.NewUserService(*userRepository)
	createUserUsecase, _ := usecase.NewCreatUserUsercase(twitterAuthService, userService, userRepository)

	userHandler, _ := handler.NewUserHandler(userRepository, eventRepository, createUserUsecase)
	friendHandler, _ := handler.NewFriendHandler(userRepository)
	eventHandler, _ := handler.NewEventHandler(eventRepository)
	twitterHandler, _ := handler.NewTwitterHandler()
	devHandler, _ := handler.NewDevHandler(*userRepository)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	userAPI := e.Group("/api/users")
	{
		// TODO 多分これidが必要
		userAPI.PUT("", userHandler.UpdateUser, jwtMiddleware)
		userAPI.GET("/:id", userHandler.GetUser, jwtMiddleware)
		userAPI.GET("/:username/by_username", userHandler.FetchUserByUsername, jwtMiddleware)
		userAPI.POST("/:id/usericons", userHandler.UploadUserIcon, jwtMiddleware)
		userAPI.POST("/twitter_verify", userHandler.CreateUserWithTwitterVerify)
	}

	eventAPI := e.Group("/api/events")
	{
		eventAPI.POST("", eventHandler.CreateEvent, jwtMiddleware)
		eventAPI.PUT("", eventHandler.UpdateEvent, jwtMiddleware)
		eventAPI.GET("/:id", eventHandler.GetEvent, jwtMiddleware)
	}

	userEventAPI := e.Group("/api/users/:id/events")
	{
		userEventAPI.GET("", userHandler.GetEvents, jwtMiddleware)
		userEventAPI.GET("/joining", eventHandler.GetJoiningEvents, jwtMiddleware)
		userEventAPI.POST("/:event_id/join", eventHandler.JoinEvent, jwtMiddleware)
	}

	friendAPI := e.Group("/api/users/:id/friends")
	{
		friendAPI.GET("", friendHandler.GetFriends, jwtMiddleware)
		friendAPI.GET("/recommended", friendHandler.GetRecommendUsers, jwtMiddleware)
		friendAPI.POST("/:friend_user_id/apply", friendHandler.ApplyFriend, jwtMiddleware)
		friendAPI.POST("/:friend_user_id/accept", friendHandler.AcceptApplyFriend, jwtMiddleware)
		friendAPI.POST("/:friend_user_id/block", friendHandler.BlockFriend, jwtMiddleware)
	}

	e.GET("/api/twitter/request_token", twitterHandler.FetchRequestToken)

	// TODO Devlop用
	e.GET("/api/restricted", devHandler.Restricted, jwtMiddleware)
	e.GET("/api/users/:id/token", devHandler.GenToken)

	e.Logger.Fatal(
		e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))),
	)
}
