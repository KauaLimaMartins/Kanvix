package dto

type WorkspaceStats struct {
	WorkspaceID    string         `json:"workspaceId"`
	ProjectCount   int            `json:"projectCount"`
	TaskCount      int            `json:"taskCount"`
	TasksByProject map[string]int `json:"tasksByProject"`
}
