package repository

import (
	"context"

	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/sql_handler"
)

type UserRepository struct {
	sqlHandler *sql_handler.SQLHandler
}

func NewUserRepository(sqlHandler *sql_handler.SQLHandler) *UserRepository {
	return &UserRepository{
		sqlHandler: sqlHandler,
	}
}

func (u UserRepository) GetUser(ctx context.Context, userID uint64) (*entity.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	var user entity.User
	rows.Next()
	err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID, &user.AppleUserIdentifier)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) FetchUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE username = $1`

	rows, err := u.sqlHandler.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	var user entity.User
	rows.Next()
	err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID, &user.AppleUserIdentifier)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) FetchUserByTwitterUserID(ctx context.Context, twitterUserID int) (*entity.User, error) {
	query := `SELECT * FROM users WHERE twitter_user_id = $1`

	rows, err := u.sqlHandler.QueryContext(ctx, query, twitterUserID)
	if err != nil {
		return nil, err
	}
	var user entity.User
	rows.Next()
	err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID, &user.AppleUserIdentifier)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) FetchUserByAppleUserIdentifier(ctx context.Context, appleUserIdentifier string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE apple_user_identifier = $1`

	rows, err := u.sqlHandler.QueryContext(ctx, query, appleUserIdentifier)
	if err != nil {
		return nil, err
	}
	var user entity.User
	rows.Next()
	err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID, &user.AppleUserIdentifier)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	query := `INSERT INTO users (username, image_url, twitter_screen_name, twitter_username, twitter_user_id, apple_user_identifier) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	rows, err := u.sqlHandler.QueryContext(ctx, query, user.Username, user.ImageURL, user.TwitterScreenName, user.TwitterUsername, user.TwitterUserID, user.AppleUserIdentifier)

	var id uint64
	rows.Next()
	rows.Scan(&id)

	newUser, _ := u.GetUser(ctx, id)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u UserRepository) UpdateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	query := `UPDATE users SET
		username = $1, image_url = $2
		twitter_screen_name = $3, twitter_user_id = $4
		twitter_username = $5, AppleUserIdentifier = $6
		WHERE id = $7`

	_, err := u.sqlHandler.QueryContext(ctx, query,
		user.Username, user.ImageURL, user.TwitterScreenName,
		user.TwitterUserID, user.TwitterUsername, user.AppleUserIdentifier,
		user.ID,
	)

	if err != nil {
		return nil, err
	}
	newUser, _ := u.GetUser(ctx, user.ID)

	return newUser, nil
}

func (u UserRepository) GetJoiningUsers(ctx context.Context, eventID uint64) ([]entity.User, error) {
	query := "SELECT * FROM users WHERE users.id in (select event_users.user_id from event_users where event_users.event_id = $1)"

	rows, err := u.sqlHandler.QueryContext(ctx, query, eventID)

	if err != nil {
		return nil, err
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID, &user.AppleUserIdentifier)
		users = append(users, user)
		if err != nil {
			return nil, err
		}
	}

	if len(users) == 0 {
		return []entity.User{}, nil
	}
	return users, nil
}
