package model

import "time"

type User struct {
	ID           string    `gorm:"primaryKey;type:text" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	Name         string    `gorm:"not null" json:"name"`
	AvatarColor  string    `gorm:"not null" json:"avatarColor"`
	Role         string    `gorm:"not null;default:'member'" json:"role"`
	Disabled     bool      `gorm:"not null;default:false" json:"disabled"`
	PasswordHash string    `gorm:"not null;default:''" json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Workspace struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	OwnerID   string    `gorm:"index;not null" json:"ownerId"`
	Owner     User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name      string    `gorm:"index;not null" json:"name"`
	Icon      string    `gorm:"not null" json:"icon"`
	Color     string    `gorm:"not null" json:"color"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WorkspaceMember struct {
	WorkspaceID string    `gorm:"primaryKey;type:text" json:"workspaceId"`
	UserID      string    `gorm:"primaryKey;type:text" json:"userId"`
	Workspace   Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Role        string    `gorm:"not null;default:'member'" json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Project struct {
	ID          string    `gorm:"primaryKey;type:text" json:"id"`
	WorkspaceID string    `gorm:"index;not null" json:"workspaceId"`
	Workspace   Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name        string    `gorm:"index;not null" json:"name"`
	Description string    `gorm:"type:text;not null;default:''" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Column struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	ProjectID string    `gorm:"index;not null" json:"projectId"`
	Project   Project   `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name      string    `gorm:"not null" json:"name"`
	Order     int       `gorm:"index;not null" json:"order"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Task struct {
	ID          string    `gorm:"primaryKey;type:text" json:"id"`
	ProjectID   string    `gorm:"index;not null" json:"projectId"`
	ColumnID    string    `gorm:"index;not null" json:"columnId"`
	Project     Project   `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Column      Column    `gorm:"foreignKey:ColumnID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text;not null;default:''" json:"description"`
	DueDate     *string   `gorm:"type:text" json:"dueDate"`
	AssigneeID  *string   `gorm:"index;type:text" json:"assigneeId"`
	Assignee    *User     `gorm:"foreignKey:AssigneeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Order       int       `gorm:"index;not null" json:"order"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Subtask struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	TaskID    string    `gorm:"index;not null" json:"taskId"`
	Task      Task      `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Title     string    `gorm:"not null" json:"title"`
	Done      bool      `gorm:"not null;default:false" json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Label struct {
	ID          string    `gorm:"primaryKey;type:text" json:"id"`
	WorkspaceID string    `gorm:"index;not null" json:"workspaceId"`
	Workspace   Workspace `gorm:"foreignKey:WorkspaceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name        string    `gorm:"index;not null" json:"name"`
	Color       string    `gorm:"not null" json:"color"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TaskLabel struct {
	TaskID    string    `gorm:"primaryKey;type:text" json:"taskId"`
	LabelID   string    `gorm:"primaryKey;type:text" json:"labelId"`
	Task      Task      `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Label     Label     `gorm:"foreignKey:LabelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

