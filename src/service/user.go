package service

import (
	"context"

	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) (*UserService, error) {
	return &UserService{
		userRepository: userRepository,
	}, nil
}

func (service UserService) ExistsByTwitterUserID(ctx context.Context, twitterUserID int) (bool, *entity.User) {
	user, _ := service.userRepository.FetchUserByTwitterUserID(ctx, twitterUserID)
	return user != nil, user
}

func (service UserService) ExistsByAppleUserIdentifier(ctx context.Context, user_identifier string) (bool, *entity.User) {
	user, _ := service.userRepository.FetchUserByAppleUserIdentifier(ctx, user_identifier)
	return user != nil, user
}

func (service UserService) GenUniqueUsername(ctx context.Context, username string) (string, error) {
	for {
		user, _ := service.userRepository.FetchUserByUsername(ctx, username)
		if user == nil {
			break
		}
		username = "_" + username
	}

	return username, nil
}
