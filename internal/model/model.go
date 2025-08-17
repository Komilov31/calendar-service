package model

import (
	"fmt"
	"strings"
	"time"
)

type Event struct {
	EventId int    `json:"event_id"`
	UserId  int    `json:"user_id" validate:"required"`
	Text    string `json:"text" validate:"required"`
	Date    Date   `json:"date" validate:"required,date_after_now"`
}

type UpdateEvent struct {
	EventId *int    `json:"event_id"`
	UserId  *int    `json:"user_id"`
	Text    *string `json:"text"`
	Date    *Date   `json:"date" validate:"date_after_now"`
}

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	layouts := []string{
		"2006-01-02",
		time.RFC3339,
	}
	var err error
	var t time.Time
	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			*d = Date(t)
			return nil
		}
	}
	return err
}

func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	formatted := fmt.Sprintf("\"%s\"", t.Format(time.DateOnly))
	return []byte(formatted), nil
}
