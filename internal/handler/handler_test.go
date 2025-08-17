package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Komilov31/calendar-service/internal/model"
	"github.com/Komilov31/calendar-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventsService is a mock implementation of EventsService
type MockEventsService struct {
	mock.Mock
}

func (m *MockEventsService) CreateEvent(event model.Event) model.Event {
	args := m.Called(event)
	return args.Get(0).(model.Event)
}

func (m *MockEventsService) UpdateEvent(updateEvent model.UpdateEvent) (model.Event, error) {
	args := m.Called(updateEvent)
	return args.Get(0).(model.Event), args.Error(1)
}

func (m *MockEventsService) DeleteEvent(userId int, eventId int) error {
	args := m.Called(userId, eventId)
	return args.Error(0)
}

func (m *MockEventsService) GetEventsForDay(userId int, date time.Time) ([]*model.Event, error) {
	args := m.Called(userId, date)
	return args.Get(0).([]*model.Event), args.Error(1)
}

func (m *MockEventsService) GetEventsForWeek(userId int, date time.Time) ([]*model.Event, error) {
	args := m.Called(userId, date)
	return args.Get(0).([]*model.Event), args.Error(1)
}

func (m *MockEventsService) GetEventsForMonth(userId int, date time.Time) ([]*model.Event, error) {
	args := m.Called(userId, date)
	return args.Get(0).([]*model.Event), args.Error(1)
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/events", h.CreateEvent)
	router.PUT("/events", h.UpdateEvent)
	router.DELETE("/events", h.DeleteEvent)
	router.GET("/events/day", h.GetEventsForDay)
	router.GET("/events/week", h.GetEventsForWeek)
	router.GET("/events/month", h.GetEventsForMonth)
	return router
}

func TestCreateEvent_Success(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	futureDate := time.Now().Add(48 * time.Hour)
	event := model.Event{
		UserId: 1,
		Text:   "Test Event",
		Date:   model.Date(futureDate),
	}

	mockService.On("CreateEvent", mock.MatchedBy(func(e model.Event) bool {
		return e.UserId == 1 && e.Text == "Test Event"
	})).Return(event)

	body, _ := json.Marshal(event)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateEvent_InvalidJSON(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("POST", "/events", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateEvent_ValidationError(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	pastDate := time.Now().Add(-48 * time.Hour)
	event := model.Event{
		UserId: 1,
		Text:   "Test Event",
		Date:   model.Date(pastDate),
	}

	body, _ := json.Marshal(event)
	req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateEvent_Success(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	futureDate := time.Now().Add(48 * time.Hour)
	updateEvent := model.UpdateEvent{
		EventId: intPtr(1),
		UserId:  intPtr(1),
		Text:    stringPtr("Updated Event"),
		Date:    (*model.Date)(&futureDate),
	}

	updatedEvent := model.Event{
		EventId: 1,
		UserId:  1,
		Text:    "Updated Event",
		Date:    model.Date(futureDate),
	}

	mockService.On("UpdateEvent", mock.MatchedBy(func(e model.UpdateEvent) bool {
		return *e.EventId == 1 && *e.UserId == 1 && *e.Text == "Updated Event"
	})).Return(updatedEvent, nil)

	body, _ := json.Marshal(updateEvent)
	req, _ := http.NewRequest("PUT", "/events?user_id=1&event_id=1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateEvent_InvalidUserID(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("PUT", "/events?user_id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateEvent_InvalidEventID(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("PUT", "/events?user_id=1&event_id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateEvent_ServiceError(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	futureDate := time.Now().Add(48 * time.Hour)
	updateEvent := model.UpdateEvent{
		EventId: intPtr(1),
		UserId:  intPtr(1),
		Text:    stringPtr("Updated Event"),
		Date:    (*model.Date)(&futureDate),
	}

	mockService.On("UpdateEvent", mock.MatchedBy(func(e model.UpdateEvent) bool {
		return *e.EventId == 1 && *e.UserId == 1
	})).Return(model.Event{}, repository.ErrNoSuchEvent)

	body, _ := json.Marshal(updateEvent)
	req, _ := http.NewRequest("PUT", "/events?user_id=1&event_id=1", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestDeleteEvent_Success(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	mockService.On("DeleteEvent", 1, 1).Return(nil)

	req, _ := http.NewRequest("DELETE", "/events?user_id=1&event_id=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteEvent_InvalidUserID(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("DELETE", "/events?user_id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteEvent_InvalidEventID(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("DELETE", "/events?user_id=1&event_id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteEvent_ServiceError(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	mockService.On("DeleteEvent", 1, 1).Return(repository.ErrNoSuchEvent)

	req, _ := http.NewRequest("DELETE", "/events?user_id=1&event_id=1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestGetEventsForDay_Success(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	events := []*model.Event{
		{
			EventId: 1,
			UserId:  1,
			Text:    "Test Event",
			Date:    model.Date(date),
		},
	}

	mockService.On("GetEventsForDay", 1, date).Return(events, nil)

	req, _ := http.NewRequest("GET", "/events/day?user_id=1&date=2024-01-15", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetEventsForDay_InvalidDate(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	req, _ := http.NewRequest("GET", "/events/day?user_id=1&date=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetEventsForWeek_Success(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	events := []*model.Event{
		{
			EventId: 1,
			UserId:  1,
			Text:    "Test Event",
			Date:    model.Date(date),
		},
	}

	mockService.On("GetEventsForWeek", 1, date).Return(events, nil)

	req, _ := http.NewRequest("GET", "/events/week?user_id=1&date=2024-01-15", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetEventsForMonth_Success(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	events := []*model.Event{
		{
			EventId: 1,
			UserId:  1,
			Text:    "Test Event",
			Date:    model.Date(date),
		},
	}

	mockService.On("GetEventsForMonth", 1, date).Return(events, nil)

	req, _ := http.NewRequest("GET", "/events/month?user_id=1&date=2024-01-15", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetEventsForDay_ServiceError(t *testing.T) {
	mockService := new(MockEventsService)
	handler := New(mockService)
	router := setupRouter(handler)

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	mockService.On("GetEventsForDay", 1, date).Return([]*model.Event(nil), repository.ErrNoSuchUser)

	req, _ := http.NewRequest("GET", "/events/day?user_id=1&date=2024-01-15", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
