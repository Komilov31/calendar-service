package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Komilov31/calendar-service/internal/model"
	"github.com/Komilov31/calendar-service/internal/repository"
	"github.com/Komilov31/calendar-service/internal/validator"
	"github.com/gin-gonic/gin"
)

type EventsService interface {
	CreateEvent(model.Event) model.Event
	UpdateEvent(model.UpdateEvent) (model.Event, error)
	DeleteEvent(int, int) error
	GetEventsForDay(int, time.Time) ([]*model.Event, error)
	GetEventsForWeek(int, time.Time) ([]*model.Event, error)
	GetEventsForMonth(int, time.Time) ([]*model.Event, error)
}

type Handler struct {
	service EventsService
}

func New(eventsService EventsService) *Handler {
	return &Handler{
		service: eventsService,
	}
}

func (h *Handler) CreateEvent(c *gin.Context) {
	var event model.Event

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(event); err != nil {
		errMsg := validator.CreateValidationErrorResponse(err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		return
	}

	event = h.service.CreateEvent(event)
	c.JSON(http.StatusOK, map[string]model.Event{"result": event})
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	var updateEvent model.UpdateEvent

	id := c.Query("user_id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id or was not provided"})
		return
	}

	e := c.Query("event_id")
	event_id, err := strconv.Atoi(e)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid event_id or was not provided"})
		return
	}

	updateEvent.EventId = &event_id
	updateEvent.UserId = &userId

	if err := c.BindJSON(&updateEvent); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(updateEvent); err != nil {
		errMsg := validator.CreateValidationErrorResponse(err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		return
	}

	event, err := h.service.UpdateEvent(updateEvent)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchEvent) || errors.Is(err, repository.ErrNoSuchUser) {
			c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]model.Event{"result": event})
}

func (h *Handler) DeleteEvent(c *gin.Context) {
	id := c.Query("user_id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id or was not provided"})
		return
	}

	e := c.Query("event_id")
	event_id, err := strconv.Atoi(e)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid event_id or was not provided"})
		return
	}

	err = h.service.DeleteEvent(userId, event_id)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchEvent) || errors.Is(err, repository.ErrNoSuchUser) {
			c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"result": "sucessfully deleted event"})
}

func (h *Handler) GetEventsForDay(c *gin.Context) {
	id := c.Query("user_id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id or was not provided"})
		return
	}

	d := c.Query("date")
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format"})
		return
	}

	events, err := h.service.GetEventsForDay(userId, date)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchEvent) || errors.Is(err, repository.ErrNoSuchUser) {
			c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string][]*model.Event{"result": events})
}

func (h *Handler) GetEventsForWeek(c *gin.Context) {
	id := c.Query("user_id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id or was not provided"})
		return
	}

	d := c.Query("date")
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format"})
		return
	}

	events, err := h.service.GetEventsForWeek(userId, date)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchEvent) || errors.Is(err, repository.ErrNoSuchUser) {
			c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string][]*model.Event{"result": events})
}

func (h *Handler) GetEventsForMonth(c *gin.Context) {
	id := c.Query("user_id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id or was not provided"})
		return
	}

	d := c.Query("date")
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format"})
		return
	}

	events, err := h.service.GetEventsForMonth(userId, date)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchEvent) || errors.Is(err, repository.ErrNoSuchUser) {
			c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string][]*model.Event{"result": events})
}
