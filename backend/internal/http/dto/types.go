package dto

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AvatarColor string `json:"avatarColor"`
}

type MeUser struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	AvatarColor string `json:"avatarColor"`
}

type Workspace struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon,omitempty"`
	Color string `json:"color,omitempty"`
}

type Project struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Column struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
}

type Task struct {
	ID          string   `json:"id"`
	ProjectID   string   `json:"projectId"`
	ColumnID    string   `json:"columnId"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	DueDate     *string  `json:"dueDate,omitempty"`
	AssigneeID  *string  `json:"assigneeId,omitempty"`
	Order       int      `json:"order"`
	CreatedAt   string   `json:"createdAt"`
}

type Label struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceId"`
	Name        string `json:"name"`
	Color       string `json:"color"`
}
