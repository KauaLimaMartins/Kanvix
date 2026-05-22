package role

import (
	"context"
	"errors"
	"time"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
)

type MemberRepo interface {
	Get(ctx context.Context, workspaceID, userID string) (entity.WorkspaceMember, error)
	Upsert(ctx context.Context, m entity.WorkspaceMember) error
}

type WorkspaceRepo interface {
	GetByID(ctx context.Context, id string) (entity.Workspace, error)
}

func Require(ctx context.Context, members MemberRepo, workspaces WorkspaceRepo, userID, workspaceID string, allowedRoles ...string) (entity.WorkspaceMember, error) {
	m, err := members.Get(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			ws, werr := workspaces.GetByID(ctx, workspaceID)
			if werr != nil {
				return entity.WorkspaceMember{}, werr
			}
			if ws.OwnerID != userID {
				return entity.WorkspaceMember{}, domain.ErrForbidden
			}
			m = entity.WorkspaceMember{
				WorkspaceID: workspaceID,
				UserID:      userID,
				Role:        "admin",
				CreatedAt:   time.Now().UTC(),
			}
			if uerr := members.Upsert(ctx, m); uerr != nil {
				return entity.WorkspaceMember{}, uerr
			}
		} else {
			return entity.WorkspaceMember{}, err
		}
	}
	if !domain.HasRole(m.Role, allowedRoles...) {
		return entity.WorkspaceMember{}, domain.ErrForbidden
	}
	return m, nil
}
