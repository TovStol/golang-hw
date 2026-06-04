package models

import "time"

type Event struct {
	ID           int64         `db:"id"            json:"id"`
	Title        string        `db:"title"         json:"title"`
	DateTime     time.Time     `db:"date_time"     json:"date_time"`
	EndDateTime  time.Time     `db:"end_date_time" json:"end_date_time"`
	Description  string        `db:"description"   json:"description,omitempty"`
	UserID       int64         `db:"user_id"       json:"user_id"`
	NotifyBefore time.Duration `db:"notify_before" json:"notify_before,omitempty"`
}
