package dto

type SearchResult struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Name        string `json:"name,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	WorkspaceID string `json:"workspaceId,omitempty"`
}

type SearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}
