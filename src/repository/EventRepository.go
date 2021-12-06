package repository

import (
	"context"
	"time"

	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/sql_handler"
)

type EventRepository struct {
	sqlHandler *sql_handler.SQLHandler
}

func NewEventRepository(sqlHandler *sql_handler.SQLHandler) *EventRepository {
	return &EventRepository{
		sqlHandler: sqlHandler,
	}
}

func (r EventRepository) GetEvent(ctx context.Context, eventID uint64) (*entity.Event, error) {
	query := `SELECT * FROM events LEFT
							JOIN users as organize_user ON organize_user.id = events.organize_user_id
							WHERE events.id = $1`

	rows, err := r.sqlHandler.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	var event entity.Event
	var organize_user entity.User
	rows.Next()
	err = rows.Scan(
		&event.ID, &event.Title, &event.Description, &event.Latitude, &event.Longitude, &event.OrganizeUserID, &event.StateDatetime, &event.EndDatetime,
		&organize_user.ID, &organize_user.Username, &organize_user.ImageURL, &organize_user.TwitterScreenName, &organize_user.TwitterUsername, &organize_user.TwitterUserID,
	)
	event.OrganizeUser = organize_user

	if err != nil {
		return nil, err
	}
	var userRepository = NewUserRepository(r.sqlHandler)
	event.JoinUsers, _ = userRepository.GetJoiningUsers(ctx, eventID)

	if len(event.JoinUsers) == 0 {
		event.JoinUsers = make([]entity.User, 0)
	}

	return &event, nil
}

func (r EventRepository) CreateEvent(ctx context.Context, event entity.Event) (*entity.Event, error) {
	query := `INSERT INTO events (title, description, latitude, longitude, organize_user_id, start_datetime, end_datetime) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.sqlHandler.QueryContext(ctx, query, event.Title, event.Description, event.Latitude, event.Longitude, event.OrganizeUserID, event.StateDatetime, event.EndDatetime)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r EventRepository) UpdateEvent(ctx context.Context, event entity.Event) (*entity.Event, error) {
	query := `UPDATE events SET title = $1, description = $2, latitude = $3, longitude = $4, organize_user_id = $5, state_date_time = $6, end_date_time = $7 WHERE id = $8`

	_, err := r.sqlHandler.QueryContext(ctx, query, event.Title, event.Description, event.Latitude, event.Longitude, event.OrganizeUserID, event.StateDatetime, event.EndDatetime, event.ID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r EventRepository) GetJoiningEvents(ctx context.Context, userID uint64) ([]entity.Event, error) {
	query := `SELECT * FROM events LEFT
							JOIN users as organize_user ON organize_user.id = events.organize_user_id
							WHERE events.id IN (SELECT event_id FROM event_users WHERE user_id = $1) OR events.organize_user_id = $2`

	rows, err := r.sqlHandler.QueryContext(ctx, query, userID, userID)

	if err != nil {
		return nil, err
	}

	var events []entity.Event

	for rows.Next() {
		var event entity.Event
		var organize_user entity.User
		err = rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.Latitude, &event.Longitude, &event.OrganizeUserID, &event.StateDatetime, &event.EndDatetime,
			&organize_user.ID, &organize_user.Username, &organize_user.ImageURL, &organize_user.TwitterScreenName, &organize_user.TwitterUsername, &organize_user.TwitterUserID,
		)
		event.OrganizeUser = organize_user
		event.JoinUsers = []entity.User{}
		events = append(events, event)
		if err != nil {
			return nil, err
		}
	}

	if len(events) == 0 {
		return make([]entity.Event, 0), nil
	}
	return events, nil
}

func (r EventRepository) JoinEvent(ctx context.Context, eventID uint64, userID uint64) (*entity.Event, error) {
	query := "INSERT INTO event_users (event_id, user_id) VALUES ($1, $2)"

	_, err := r.sqlHandler.QueryContext(ctx, query, eventID, userID)

	if err != nil {
		return nil, err
	}

	event, err := r.GetEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r EventRepository) GetEventsRelatedToUser(ctx context.Context, user entity.User) ([]entity.Event, error) {
	query := `
	SELECT * FROM events LEFT
	JOIN users as organize_user ON organize_user.id = events.organize_user_id
	WHERE
		events.start_datetime > $2 and (
			events.organize_user_id IN (SELECT friend_user_id FROM friends WHERE friends.user_id = $1 and friends.status = 'accepted')
			OR events.organize_user_id = $3
		)
	`

	location := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(location)
	nowStr := now.Format(time.RFC3339)

	rows, err := r.sqlHandler.QueryContext(ctx, query, user.ID, nowStr, user.ID)
	if err != nil {
		return nil, err
	}

	var events = []entity.Event{}
	for rows.Next() {
		var event entity.Event
		var organize_user entity.User
		err = rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.Latitude, &event.Longitude, &event.OrganizeUserID, &event.StateDatetime, &event.EndDatetime,
			&organize_user.ID, &organize_user.Username, &organize_user.ImageURL, &organize_user.TwitterScreenName, &organize_user.TwitterUsername, &organize_user.TwitterUserID,
		)
		event.OrganizeUser = organize_user

		if err != nil {
			return nil, err
		}

		// TODO 現状 N+1になってしまうので取得していない
		event.JoinUsers = make([]entity.User, 0)
		events = append(events, event)
	}

	if len(events) == 0 {
		return make([]entity.Event, 0), nil
	}

	return events, nil
}
