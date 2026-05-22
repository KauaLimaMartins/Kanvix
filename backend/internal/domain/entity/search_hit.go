package entity

type ProjectSearchHit struct {
	ID   string
	Name string
}

type TaskSearchHit struct {
	ID        string
	Title     string
	ProjectID string
}

