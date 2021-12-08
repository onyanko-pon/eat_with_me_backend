package usecase

import (
	"context"

	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
	"github.com/onyanko-pon/eat_with_me_backend/src/service"
)

type CreateUserUsecase struct {
	twitterAuthService *service.TwitterAuthService
	userService        *service.UserService
	userRepository     *repository.UserRepository
}

func NewCreatUserUsercase(twitterAuthService *service.TwitterAuthService, userService *service.UserService, userRepository *repository.UserRepository) (*CreateUserUsecase, error) {
	return &CreateUserUsecase{
		twitterAuthService: twitterAuthService,
		userService:        userService,
		userRepository:     userRepository,
	}, nil
}

func (usecase CreateUserUsecase) CreateUserWithTwitterVerify(ctx context.Context, oauthToken string, oauthSecret string, oauthVerifier string) (*entity.User, error) {
	accessToken, accessSecret, err := usecase.twitterAuthService.GenAccessToken(oauthToken, oauthSecret, oauthVerifier)

	if err != nil {
		return nil, err
	}

	twitterUser, err := usecase.twitterAuthService.VerifyUser(accessToken, accessSecret)
	if err != nil {
		return nil, err
	}

	exists, user := usecase.userService.ExistsByTwitterUserID(ctx, int(twitterUser.ID))
	if exists {
		return user, nil
	}

	username, err := usecase.userService.GenUniqueUsername(ctx, twitterUser.ScreenName)
	if err != nil {
		return nil, err
	}

	user = &entity.User{
		ID:                0,
		Username:          username,
		ImageURL:          twitterUser.ProfileImageUrlHttps,
		TwitterScreenName: twitterUser.ScreenName,
		TwitterUsername:   twitterUser.Name,
		TwitterUserID:     twitterUser.ID,
	}

	user, err = usecase.userRepository.CreateUser(ctx, *user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
