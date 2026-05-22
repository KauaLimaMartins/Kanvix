package entity

import "time"

type Project struct {
	ID          string
	WorkspaceID string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Column struct {
	ID        string
	ProjectID string
	Name      string
	Order     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Task struct {
	ID          string
	ProjectID   string
	ColumnID    string
	Title       string
	Description string
	DueDate     *string
	AssigneeID  *string
	Order       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

