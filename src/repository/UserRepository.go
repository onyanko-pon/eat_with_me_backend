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
	err = rows.Scan(&user.ID, &user.Username, &user.ImageURL)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	query := `INSERT INTO users (username, image_url) VALUES ($1, $2)`

	_, err := u.sqlHandler.QueryContext(ctx, query, user.Username, user.ImageURL)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u UserRepository) UpdateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	query := `UPDATE users SET username = $1, image_url = $2 $1 WHERE id = $3`

	_, err := u.sqlHandler.QueryContext(ctx, query, user.Username, user.ImageURL, user.ID)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (u UserRepository) GetFriends(ctx context.Context, userID uint64) ([]entity.User, error) {
	query := "SELECT * FROM users WHERE id IN (SELECT friend_user_id FROM friends WHERE user_id = $1)"

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Username, &user.ImageURL)
		users = append(users, user)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}