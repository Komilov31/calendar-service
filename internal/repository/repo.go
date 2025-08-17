package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/Komilov31/calendar-service/internal/model"
)

var (
	ErrNoSuchEvent = errors.New("no such event in database")
	ErrNoSuchUser  = errors.New("no such user in database")
)

type Repository struct {
	mu     *sync.RWMutex
	events map[int][]*model.Event
}

func New() *Repository {
	return &Repository{
		mu:     &sync.RWMutex{},
		events: make(map[int][]*model.Event),
	}
}

func (r *Repository) CreateEvent(event model.Event) model.Event {
	r.mu.Lock()
	defer r.mu.Unlock()
	events := r.events[event.UserId]

	event.EventId = len(events) + 1
	events = append(events, &event)

	r.events[event.UserId] = events

	return event
}

func (r *Repository) UpdateEvent(updateEvent model.UpdateEvent) (model.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	event, ok := r.getEventByUserId(*updateEvent.UserId, *updateEvent.EventId)
	if !ok {
		return model.Event{}, ErrNoSuchEvent
	}

	if updateEvent.Text != nil {
		event.Text = *updateEvent.Text
	}

	if updateEvent.Date != nil {
		event.Date = *updateEvent.Date
	}

	return *event, nil
}

func (r *Repository) DeleteEvent(userId int, eventId int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	events, ok := r.events[userId]
	if !ok {
		return ErrNoSuchUser
	}

	found := false
	for i, event := range events {
		if event.EventId == eventId {
			found = true
			events[i] = events[len(events)-1]
			events = events[:len(events)-1]
			break
		}
	}

	if !found {
		return ErrNoSuchEvent
	}

	r.events[userId] = events

	return nil
}

func (r *Repository) GetEventsForDay(userId int, date time.Time) ([]*model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events, ok := r.events[userId]
	if !ok {
		return nil, ErrNoSuchUser
	}

	var eventsForDay []*model.Event
	for _, event := range events {
		curDate := time.Time(event.Date)
		eventsDay, eventsYear := curDate.Day(), curDate.Year()
		day, year := date.Day(), date.Year()

		if eventsDay == day && eventsYear == year {
			eventsForDay = append(eventsForDay, event)
		}
	}

	return eventsForDay, nil
}

func (r *Repository) GetEventsForWeek(userId int, date time.Time) ([]*model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events, ok := r.events[userId]
	if !ok {
		return nil, ErrNoSuchUser
	}

	var eventsForWeek []*model.Event
	for _, event := range events {
		curDate := time.Time(event.Date)

		currentYear, currentWeek := date.ISOWeek()
		eventYear, eventWeek := curDate.ISOWeek()

		if currentYear == eventYear && currentWeek == eventWeek {
			eventsForWeek = append(eventsForWeek, event)
		}
	}

	return eventsForWeek, nil
}

func (r *Repository) GetEventsForMonth(userId int, date time.Time) ([]*model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events, ok := r.events[userId]
	if !ok {
		return nil, ErrNoSuchUser
	}

	var eventsForMonth []*model.Event
	for _, event := range events {
		curDate := time.Time(event.Date)
		eventsMonth, eventsYear := curDate.Month(), curDate.Year()
		month, year := date.Month(), date.Year()

		if eventsMonth == month && eventsYear == year {
			eventsForMonth = append(eventsForMonth, event)
		}
	}

	return eventsForMonth, nil
}
