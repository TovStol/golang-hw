package models

import "time"

type Notification struct {
	EventID   int64     `db:"event_id"   json:"event_id"`
	Title     string    `db:"title"      json:"title"`
	EventDate time.Time `db:"event_date" json:"event_date"`
	UserID    int64     `db:"user_id"    json:"user_id"`
}
