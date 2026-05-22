package entity

import "time"

type Subtask struct {
	ID        string
	TaskID    string
	Title     string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Label struct {
	ID          string
	WorkspaceID string
	Name        string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskLabel struct {
	TaskID    string
	LabelID   string
	CreatedAt time.Time
}

