package service

import (
	"time"

	"github.com/Komilov31/calendar-service/internal/model"
)

type EventStorage interface {
	CreateEvent(model.Event) model.Event
	UpdateEvent(model.UpdateEvent) (model.Event, error)
	DeleteEvent(int, int) error
	GetEventsForDay(int, time.Time) ([]*model.Event, error)
	GetEventsForWeek(int, time.Time) ([]*model.Event, error)
	GetEventsForMonth(int, time.Time) ([]*model.Event, error)
}

type Service struct {
	storage EventStorage
}

func New(storage EventStorage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateEvent(event model.Event) model.Event {
	return s.storage.CreateEvent(event)
}

func (s *Service) UpdateEvent(updateEvent model.UpdateEvent) (model.Event, error) {
	return s.storage.UpdateEvent(updateEvent)
}

func (s *Service) DeleteEvent(userId int, eventId int) error {
	return s.storage.DeleteEvent(userId, eventId)
}

func (s *Service) GetEventsForDay(userId int, date time.Time) ([]*model.Event, error) {
	return s.storage.GetEventsForDay(userId, date)
}

func (s *Service) GetEventsForWeek(userId int, date time.Time) ([]*model.Event, error) {
	return s.storage.GetEventsForWeek(userId, date)
}

func (s *Service) GetEventsForMonth(userId int, date time.Time) ([]*model.Event, error) {
	return s.storage.GetEventsForMonth(userId, date)
}
