package repository

import "github.com/Komilov31/calendar-service/internal/model"

func (r *Repository) getEventByUserId(userId int, eventId int) (*model.Event, bool) {
	events := r.events[userId]

	for _, event := range events {
		if event.EventId == eventId {
			return event, true
		}
	}

	return nil, false
}
