package entity

import "time"

type Workspace struct {
	ID        string
	OwnerID   string
	Name      string
	Icon      string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WorkspaceMember struct {
	WorkspaceID string
	UserID      string
	Role        string
	CreatedAt   time.Time
}

