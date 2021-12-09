package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/sql_handler"
)

type FriendRepository struct {
	sqlHandler *sql_handler.SQLHandler
}

func NewFriendRepository(sqlHandler *sql_handler.SQLHandler) *FriendRepository {
	return &FriendRepository{
		sqlHandler: sqlHandler,
	}
}

func (u FriendRepository) GetFriends(ctx context.Context, userID uint64) ([]entity.Friend, error) {
	query := `
	SELECT user_relations.blinding, users.* FROM user_relations
	LEFT JOIN users ON users.id = user_relations.friend_user_id
	WHERE user_relations.user_id = $1;`

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var friends []entity.Friend
	for rows.Next() {
		var user entity.User
		var friend entity.Friend
		err = rows.Scan(&friend.Blinding, &user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
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

func (u FriendRepository) GetApplyings(ctx context.Context, userID uint64) ([]entity.FriendApply, error) {
	query := `
	SELECT users.* FROM friend_applys
		LEFT JOIN users ON users.id = friend_applys.friend_user_id -- 相手のユーザー
		WHERE friend_applys.user_id = $1 AND friend_applys.accepted_at IS NULL
	`

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var applys []entity.FriendApply
	for rows.Next() {
		var user entity.User
		var apply entity.FriendApply
		err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
		if err != nil {
			return nil, err
		}

		apply.User = user
		applys = append(applys, apply)
	}

	if len(applys) == 0 {
		return []entity.FriendApply{}, nil
	}
	return applys, nil
}

func (u FriendRepository) GetApplieds(ctx context.Context, userID uint64) ([]entity.FriendApply, error) {
	query := `
	SELECT users.* FROM friend_applys
		LEFT JOIN users ON users.id = friend_applys.user_id -- 相手のユーザー
		WHERE friend_applys.friend_user_id = $1 AND friend_applys.accepted_at IS NULL;
	`

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var applys []entity.FriendApply
	for rows.Next() {
		var user entity.User
		var apply entity.FriendApply
		err = rows.Scan(&user.ID, &user.Username, &user.ImageURL, &user.TwitterScreenName, &user.TwitterUsername, &user.TwitterUserID)
		if err != nil {
			return nil, err
		}

		apply.User = user
		applys = append(applys, apply)
	}

	if len(applys) == 0 {
		return []entity.FriendApply{}, nil
	}
	return applys, nil
}

func (u FriendRepository) GetRecommendUsers(ctx context.Context, userID uint64) ([]entity.User, error) {
	friendUserIDStrList := []string{"0"}
	friends, _ := u.GetFriends(ctx, userID)

	for _, friend := range friends {
		friendUserIDStrList = append(friendUserIDStrList, strconv.Itoa(int(friend.User.ID)))
	}

	query := fmt.Sprintf(`
		SELECT users.* FROM user_relations
			LEFT JOIN user_relations AS friend_user_relations
				ON user_relations.friend_user_id = friend_user_relations.user_id -- 友達のユーザーIDと繋げて、友達の友達を探す
			LEFT JOIN users ON users.id = friend_user_relations.friend_user_id
			WHERE NOT users.id = $1 AND user_relations.user_id = $2 AND NOT users.id IN (%s)
			GROUP by users.id;
		`,
		strings.Join(friendUserIDStrList, ","),
	)

	rows, err := u.sqlHandler.QueryContext(ctx, query, userID, userID)

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

func (u FriendRepository) Apply(ctx context.Context, userID uint64, friendUserID uint64) error {
	query := "INSERT INTO friend_applys (user_id, friend_user_id) VALUES ($1, $2);"
	_, err := u.sqlHandler.QueryContext(ctx, query, userID, friendUserID)

	return err
}

func (u FriendRepository) AcceptApply(ctx context.Context, userID uint64, friendUserID uint64) error {
	location := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(location)
	nowStr := now.Format(time.RFC3339)

	query := "UPDATE friend_applys SET accepted_at = $1 WHERE user_id = $2 AND friend_user_id = $3;"
	_, err := u.sqlHandler.QueryContext(ctx, query, nowStr, userID, friendUserID)
	if err != nil {
		return err
	}

	query = "INSERT INTO user_relations (user_id, friend_user_id) VALUES ($1, $2), ($3, $4);"
	_, err = u.sqlHandler.QueryContext(ctx, query, userID, friendUserID, friendUserID, userID)
	return err
}

// TODO トランザクション
func (u FriendRepository) Blind(ctx context.Context, userID uint64, friendUserID uint64) error {
	query := "UPDATE user_relations SET brainding = TRUE where user_id = $1 AND friend_user_id = $2;"
	_, err := u.sqlHandler.QueryContext(ctx, query, userID, friendUserID)
	if err != nil {
		return err
	}

	query = "UPDATE user_relations SET brainded = TRUE where friend_user_id = $1 AND user_id = $2;"
	_, err = u.sqlHandler.QueryContext(ctx, query, userID, friendUserID)
	return err
}
