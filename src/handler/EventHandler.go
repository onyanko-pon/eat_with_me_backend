package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/onyanko-pon/eat_with_me_backend/src/entity"
	"github.com/onyanko-pon/eat_with_me_backend/src/repository"
)

type EventHandler struct {
	EventRepository *repository.EventRepository
}

func NewEventHandler(eventRepository *repository.EventRepository) (*EventHandler, error) {
	return &EventHandler{
		EventRepository: eventRepository,
	}, nil
}

type responseGetEvent struct {
	Event *entity.Event `json:"event"`
}

func (h EventHandler) GetEvent(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	event, _ := h.EventRepository.GetEvent(c.Request().Context(), uint64(id))

	return c.JSON(http.StatusOK, responseGetEvent{
		Event: event,
	})
}

type responseCreateEvent struct {
	Event *entity.Event `json:"event"`
}

func (h EventHandler) CreateEvent(c echo.Context) error {

	event := new(entity.Event)
	if err := c.Bind(event); err != nil {
		return err
	}

	h.EventRepository.CreateEvent(c.Request().Context(), *event)

	return c.JSON(http.StatusOK, responseCreateEvent{
		Event: event,
	})
}

type responseUpdateEvent struct {
	Event *entity.Event `json:"event"`
}

func (h EventHandler) UpdateEvent(c echo.Context) error {

	event := new(entity.Event)
	if err := c.Bind(event); err != nil {
		return err
	}

	h.EventRepository.UpdateEvent(c.Request().Context(), *event)

	return c.JSON(http.StatusOK, responseUpdateEvent{
		Event: event,
	})
}
