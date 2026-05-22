package workspace

import (
	"context"
	"encoding/json"
	"time"

	"kanvix/backend/internal/usecase/workspace/role"
)

type Out struct {
	WorkspaceID    string         `json:"workspaceId"`
	ProjectCount   int            `json:"projectCount"`
	TaskCount      int            `json:"taskCount"`
	TasksByProject map[string]int `json:"tasksByProject"`
}

type statsRepo interface {
	CountProjects(ctx context.Context, workspaceID string) (int64, error)
	CountTasksByProject(ctx context.Context, workspaceID string) (map[string]int, int, error)
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
	Stats      statsRepo
	Cache      cache
	CacheTTL   time.Duration
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, userID, workspaceID string) (Out, error) {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, userID, workspaceID, "admin", "member"); err != nil {
		return Out{}, err
	}

	key := "stats:" + workspaceID
	if uc.Cache != nil {
		if raw, ok, err := uc.Cache.GetBytes(ctx, key); err == nil && ok {
			var cached Out
			if json.Unmarshal(raw, &cached) == nil {
				return cached, nil
			}
		}
	}

	projectCount, err := uc.Stats.CountProjects(ctx, workspaceID)
	if err != nil {
		return Out{}, err
	}
	byProject, total, err := uc.Stats.CountTasksByProject(ctx, workspaceID)
	if err != nil {
		return Out{}, err
	}

	out := Out{
		WorkspaceID:    workspaceID,
		ProjectCount:   int(projectCount),
		TaskCount:      total,
		TasksByProject: byProject,
	}

	if uc.Cache != nil {
		if b, err := json.Marshal(out); err == nil {
			_ = uc.Cache.SetBytes(ctx, key, b, uc.CacheTTL)
		}
	}

	return out, nil
}

