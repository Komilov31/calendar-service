package repository

import (
	"testing"
	"time"

	"github.com/Komilov31/calendar-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	repo := New()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.mu)
	assert.NotNil(t, repo.events)
	assert.Empty(t, repo.events)
}

func TestCreateEvent(t *testing.T) {
	repo := New()

	t.Run("Create first event for user", func(t *testing.T) {
		event := model.Event{
			UserId: 1,
			Text:   "Test Event 1",
			Date:   model.Date(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)),
		}

		result := repo.CreateEvent(event)
		assert.Equal(t, 1, result.EventId)
		assert.Equal(t, "Test Event 1", result.Text)
		assert.Equal(t, 1, result.UserId)
	})

	t.Run("Create second event for same user", func(t *testing.T) {
		event := model.Event{
			UserId: 1,
			Text:   "Test Event 2",
			Date:   model.Date(time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)),
		}

		result := repo.CreateEvent(event)
		assert.Equal(t, 2, result.EventId)
		assert.Equal(t, "Test Event 2", result.Text)
	})

	t.Run("Create event for different user", func(t *testing.T) {
		event := model.Event{
			UserId: 2,
			Text:   "Test Event for User 2",
			Date:   model.Date(time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC)),
		}

		result := repo.CreateEvent(event)
		assert.Equal(t, 1, result.EventId)
		assert.Equal(t, "Test Event for User 2", result.Text)
		assert.Equal(t, 2, result.UserId)
	})
}

func TestUpdateEvent(t *testing.T) {
	repo := New()

	event1 := model.Event{
		UserId: 1,
		Text:   "Original Text",
		Date:   model.Date(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event1)

	event2 := model.Event{
		UserId: 1,
		Text:   "Another Event",
		Date:   model.Date(time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event2)

	t.Run("Update event text only", func(t *testing.T) {
		newText := "Updated Text"
		updateEvent := model.UpdateEvent{
			EventId: intPtr(1),
			UserId:  intPtr(1),
			Text:    &newText,
		}

		result, err := repo.UpdateEvent(updateEvent)
		assert.NoError(t, err)
		assert.Equal(t, 1, result.EventId)
		assert.Equal(t, "Updated Text", result.Text)
		assert.Equal(t, model.Date(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)), result.Date)
	})

	t.Run("Update event date only", func(t *testing.T) {
		newDate := model.Date(time.Date(2024, 1, 20, 10, 0, 0, 0, time.UTC))
		updateEvent := model.UpdateEvent{
			EventId: intPtr(2),
			UserId:  intPtr(1),
			Date:    &newDate,
		}

		result, err := repo.UpdateEvent(updateEvent)
		assert.NoError(t, err)
		assert.Equal(t, 2, result.EventId)
		assert.Equal(t, "Another Event", result.Text)
		assert.Equal(t, newDate, result.Date)
	})

	t.Run("Update both text and date", func(t *testing.T) {
		newText := "Completely Updated"
		newDate := model.Date(time.Date(2024, 1, 25, 10, 0, 0, 0, time.UTC))
		updateEvent := model.UpdateEvent{
			EventId: intPtr(1),
			UserId:  intPtr(1),
			Text:    &newText,
			Date:    &newDate,
		}

		result, err := repo.UpdateEvent(updateEvent)
		assert.NoError(t, err)
		assert.Equal(t, 1, result.EventId)
		assert.Equal(t, newText, result.Text)
		assert.Equal(t, newDate, result.Date)
	})

	t.Run("Update non-existent event", func(t *testing.T) {
		updateEvent := model.UpdateEvent{
			EventId: intPtr(999),
			UserId:  intPtr(1),
			Text:    stringPtr("Should not work"),
		}

		result, err := repo.UpdateEvent(updateEvent)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSuchEvent, err)
		assert.Equal(t, model.Event{}, result)
	})

	t.Run("Update event for non-existent user", func(t *testing.T) {
		updateEvent := model.UpdateEvent{
			EventId: intPtr(1),
			UserId:  intPtr(999),
			Text:    stringPtr("Should not work"),
		}

		result, err := repo.UpdateEvent(updateEvent)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSuchEvent, err)
		assert.Equal(t, model.Event{}, result)
	})
}

func TestDeleteEvent(t *testing.T) {
	repo := New()

	event1 := model.Event{
		UserId: 1,
		Text:   "Event 1",
		Date:   model.Date(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event1)

	event2 := model.Event{
		UserId: 1,
		Text:   "Event 2",
		Date:   model.Date(time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event2)

	event3 := model.Event{
		UserId: 2,
		Text:   "Event for User 2",
		Date:   model.Date(time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event3)

	t.Run("Delete existing event", func(t *testing.T) {
		err := repo.DeleteEvent(1, 1)
		assert.NoError(t, err)

		events, _ := repo.GetEventsForDay(1, time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC))
		assert.Len(t, events, 1)
		assert.Equal(t, 2, events[0].EventId)
	})

	t.Run("Delete non-existent event", func(t *testing.T) {
		err := repo.DeleteEvent(1, 999)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSuchEvent, err)
	})

	t.Run("Delete event from non-existent user", func(t *testing.T) {
		err := repo.DeleteEvent(999, 1)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSuchUser, err)
	})

}

func TestGetEventsForDay(t *testing.T) {
	repo := New()

	event1 := model.Event{
		UserId: 1,
		Text:   "Event Jan 15",
		Date:   model.Date(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event1)

	event2 := model.Event{
		UserId: 1,
		Text:   "Event Jan 15 Afternoon",
		Date:   model.Date(time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event2)

	event3 := model.Event{
		UserId: 1,
		Text:   "Event Jan 16",
		Date:   model.Date(time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event3)

	event4 := model.Event{
		UserId: 2,
		Text:   "Event for User 2",
		Date:   model.Date(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)),
	}
	repo.CreateEvent(event4)

	t.Run("Get events for specific day", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForDay(1, date)
		assert.NoError(t, err)
		assert.Len(t, events, 2)
		assert.Equal(t, "Event Jan 15", events[0].Text)
		assert.Equal(t, "Event Jan 15 Afternoon", events[1].Text)
	})

	t.Run("Get events for different day", func(t *testing.T) {
		date := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForDay(1, date)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "Event Jan 16", events[0].Text)
	})

	t.Run("Get events for day with no events", func(t *testing.T) {
		date := time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForDay(1, date)
		assert.NoError(t, err)
		assert.Empty(t, events)
	})

	t.Run("Get events for non-existent user", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForDay(999, date)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSuchUser, err)
		assert.Nil(t, events)
	})

	t.Run("User isolation - user 2 events", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForDay(2, date)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "Event for User 2", events[0].Text)
	})
}

func TestGetEventsForWeek(t *testing.T) {
	repo := New()

	events := []model.Event{
		{UserId: 1, Text: "Monday Event", Date: model.Date(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC))},
		{UserId: 1, Text: "Wednesday Event", Date: model.Date(time.Date(2024, 1, 17, 10, 0, 0, 0, time.UTC))},
		{UserId: 1, Text: "Sunday Event", Date: model.Date(time.Date(2024, 1, 21, 10, 0, 0, 0, time.UTC))},
		{UserId: 1, Text: "Next Monday Event", Date: model.Date(time.Date(2024, 1, 22, 10, 0, 0, 0, time.UTC))},
		{UserId: 2, Text: "User 2 Monday Event", Date: model.Date(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC))},
	}

	for _, event := range events {
		repo.CreateEvent(event)
	}

	t.Run("Get events for week starting Monday", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForWeek(1, date)
		assert.NoError(t, err)
		assert.Len(t, events, 3)
	})

	t.Run("Get events for week starting Wednesday", func(t *testing.T) {
		date := time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForWeek(1, date)
		assert.NoError(t, err)
		assert.Len(t, events, 3)
	})

	t.Run("Get events for different week", func(t *testing.T) {
		date := time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForWeek(1, date)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "Next Monday Event", events[0].Text)
	})

	t.Run("Get events for non-existent user", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForWeek(999, date)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSuchUser, err)
		assert.Nil(t, events)
	})

	t.Run("User isolation", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForWeek(2, date)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "User 2 Monday Event", events[0].Text)
	})
}

func TestGetEventsForMonth(t *testing.T) {
	repo := New()

	events := []model.Event{
		{UserId: 1, Text: "Jan 1 Event", Date: model.Date(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))},
		{UserId: 1, Text: "Jan 15 Event", Date: model.Date(time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC))},
		{UserId: 1, Text: "Jan 31 Event", Date: model.Date(time.Date(2024, 1, 31, 10, 0, 0, 0, time.UTC))},
		{UserId: 1, Text: "Feb 1 Event", Date: model.Date(time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC))},
		{UserId: 2, Text: "User 2 Jan Event", Date: model.Date(time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC))},
	}

	for _, event := range events {
		repo.CreateEvent(event)
	}

	t.Run("Get events for January", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForMonth(1, date)
		assert.NoError(t, err)
		assert.Len(t, events, 3)
	})

	t.Run("Get events for February", func(t *testing.T) {
		date := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForMonth(1, date)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "Feb 1 Event", events[0].Text)
	})

	t.Run("Get events for month with no events", func(t *testing.T) {
		date := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForMonth(1, date)
		assert.NoError(t, err)
		assert.Empty(t, events)
	})

	t.Run("Get events for non-existent user", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForMonth(999, date)
		assert.Error(t, err)
		assert.Equal(t, ErrNoSuchUser, err)
		assert.Nil(t, events)
	})

	t.Run("User isolation", func(t *testing.T) {
		date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		events, err := repo.GetEventsForMonth(2, date)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "User 2 Jan Event", events[0].Text)
	})
}

func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
