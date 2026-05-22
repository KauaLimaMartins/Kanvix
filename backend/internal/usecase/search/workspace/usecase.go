package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type Result struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Name        string `json:"name,omitempty"`
	ProjectID   string `json:"projectId,omitempty"`
	WorkspaceID string `json:"workspaceId,omitempty"`
}

type Out struct {
	Query   string   `json:"query"`
	Results []Result `json:"results"`
}

type searchRepo interface {
	SearchProjects(ctx context.Context, workspaceID string, like string, limit int) ([]entity.ProjectSearchHit, error)
	SearchTasks(ctx context.Context, workspaceID string, like string, limit int) ([]entity.TaskSearchHit, error)
}

type cache interface {
	GetBytes(ctx context.Context, key string) ([]byte, bool, error)
	SetBytes(ctx context.Context, key string, value []byte, ttl time.Duration) error
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Search     searchRepo
	Cache      cache
	CacheTTL   time.Duration
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, userID, workspaceID, q string, limit int) (Out, error) {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, userID, workspaceID, "admin", "member"); err != nil {
		return Out{}, err
	}

	q = strings.TrimSpace(q)
	if q == "" {
		return Out{Query: "", Results: []Result{}}, nil
	}
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	cacheKey := ""
	if uc.Cache != nil {
		cacheKey = "search:" + workspaceID + ":" + strings.ToLower(q) + ":" + fmt.Sprintf("%d", limit)
		if raw, ok, err := uc.Cache.GetBytes(ctx, cacheKey); err == nil && ok {
			var cached Out
			if json.Unmarshal(raw, &cached) == nil {
				return cached, nil
			}
		}
	}

	like := "%" + q + "%"
	projects, err := uc.Search.SearchProjects(ctx, workspaceID, like, limit)
	if err != nil {
		return Out{}, err
	}
	tasks, err := uc.Search.SearchTasks(ctx, workspaceID, like, limit)
	if err != nil {
		return Out{}, err
	}

	results := make([]Result, 0, len(projects)+len(tasks))
	for _, p := range projects {
		results = append(results, Result{Type: "project", ID: p.ID, Name: p.Name, WorkspaceID: workspaceID})
	}
	for _, t := range tasks {
		results = append(results, Result{Type: "task", ID: t.ID, Title: t.Title, ProjectID: t.ProjectID, WorkspaceID: workspaceID})
	}

	out := Out{Query: q, Results: results}
	if uc.Cache != nil && cacheKey != "" {
		if b, err := json.Marshal(out); err == nil {
			_ = uc.Cache.SetBytes(ctx, cacheKey, b, uc.CacheTTL)
		}
	}
	return out, nil
}

