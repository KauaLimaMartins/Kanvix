package entity

import "time"

type TaskOrderUpdate struct {
	ID        string
	ColumnID  string
	Order     int
	UpdatedAt time.Time
}

