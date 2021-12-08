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
	err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
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
	err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) CreateUser(ctx context.Context, user entity.User) (*entity.User, error) {
	query := `INSERT INTO users (username, image_url, twitter_screen_name, twitter_username, twitter_user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	rows, err := u.sqlHandler.QueryContext(ctx, query, user.Username, user.ImageURL, user.TwitterScreenName, user.TwitterUsername, user.TwitterUserID)

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
	query := `UPDATE users SET username = $1, image_url = $2 WHERE id = $3`

	_, err := u.sqlHandler.QueryContext(ctx, query, user.Username, user.ImageURL, user.ID)

	if err != nil {
		return nil, err
	}
	newUser, _ := u.GetUser(ctx, user.ID)

	return newUser, nil
}

func (u UserRepository) GetFriends(ctx context.Context, userID uint64) ([]entity.Friend, error) {
	query := `
	SELECT friends.status, users.* FROM friends
	LEFT JOIN users ON users.id = friends.friend_user_id
	WHERE friends.user_id = $1 and friends.status = 'accepted'`

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var friends []entity.Friend
	for rows.Next() {
		var user entity.User
		var friend entity.Friend
		err = rows.Scan(&friend.Status, &user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
		if err != nil {
			return nil, err
		}

		friend.User = user
		friends = append(friends, friend)
	}

	if len(friends) == 0 {
		return []entity.Friend{}, nil
	}
	return friends, nil
}

func (u UserRepository) GetRequestFriends(ctx context.Context, userID uint64) ([]entity.Friend, error) {
	query := `
	SELECT * FROM friends
	LEFT JOIN users ON users.id = friends.user_id
	WHERE friends.friend_user_id = $1 AND friends.status = 'applying';`

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var friends []entity.Friend
	for rows.Next() {
		var user entity.User
		var friend entity.Friend
		err = rows.Scan(&friend.Status, &user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
		if err != nil {
			return nil, err
		}

		friend.User = user
		friends = append(friends, friend)
	}

	if len(friends) == 0 {
		return []entity.Friend{}, nil
	}
	return friends, nil
}

func (u UserRepository) GetRecommendUsers(ctx context.Context, userID uint64) ([]entity.User, error) {
	query := `
	SELECT users.* FROM friends
		JOIN friends AS recommend_friends ON friends.friend_user_id = recommend_friends.user_id AND NOT recommend_friends.friend_user_id = $1
		JOIN users ON recommend_friends.friend_user_id = users.id
		WHERE friends.user_id = $2 AND friends.status = 'accepted';
	`

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
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

func (u UserRepository) GetJoiningUsers(ctx context.Context, eventID uint64) ([]entity.User, error) {
	query := "SELECT * FROM users WHERE users.id in (select event_users.user_id from event_users where event_users.event_id = $1)"

	rows, err := u.sqlHandler.QueryContext(ctx, query, eventID)

	if err != nil {
		return nil, err
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
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

func (u UserRepository) ApplyFriend(ctx context.Context, userID uint64, friendUserID uint64) error {
	query := "INSERT INTO friends (user_id, friend_user_id, status) VALUES ($1, $2, 'applying')"

	_, err := u.sqlHandler.QueryContext(ctx, query, userID, friendUserID)
	return err
}

func (u UserRepository) AcceptApplyFriend(ctx context.Context, userID uint64, friendUserID uint64) error {
	query := "UPDATE friends SET status = 'accepted' where user_id = $1 AND friend_user_id = $2"

	_, err := u.sqlHandler.QueryContext(ctx, query, userID, friendUserID)
	return err
}

func (u UserRepository) BlockFriend(ctx context.Context, userID uint64, friendUserID uint64) error {
	query := "UPDATE friends SET status = 'blocked' where user_id = $1 AND friend_user_id = $2"

	_, err := u.sqlHandler.QueryContext(ctx, query, userID, friendUserID)
	return err
}
