package repository

import (
	"context"

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
	query := `SELECT * FROM events WHERE id = $1`

	rows, err := r.sqlHandler.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	var event entity.Event
	rows.Next()
	err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.Latitude, &event.Longitude, &event.OrganizeUserID, &event.StateDatetime, &event.EndDatetime)
	if err != nil {
		return nil, err
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

// 参加者のエンドポイント
