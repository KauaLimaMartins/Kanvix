package services

import (
	"context"
	"errors"
	"time"

	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func hasRole(role string, allowed ...string) bool {
	for _, a := range allowed {
		if role == a {
			return true
		}
	}
	return false
}

func (s AppService) requireGlobalAdmin(ctx context.Context, userID string) error {
	u, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if u.Role != "admin" {
		return repositories.ErrForbidden
	}
	return nil
}

func (s AppService) requireWorkspaceRole(ctx context.Context, userID, workspaceID string, allowedRoles ...string) (models.WorkspaceMember, error) {
	m, err := s.Repo.GetWorkspaceMember(ctx, workspaceID, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			ws, werr := s.Repo.GetWorkspaceByID(ctx, workspaceID)
			if werr != nil {
				return models.WorkspaceMember{}, werr
			}
			if ws.OwnerID != userID {
				return models.WorkspaceMember{}, repositories.ErrForbidden
			}
			m = models.WorkspaceMember{
				WorkspaceID: workspaceID,
				UserID:      userID,
				Role:        "admin",
				CreatedAt:   time.Now().UTC(),
			}
			if uerr := s.Repo.UpsertWorkspaceMember(ctx, m); uerr != nil {
				return models.WorkspaceMember{}, uerr
			}
		} else {
			return models.WorkspaceMember{}, err
		}
	}
	if !hasRole(m.Role, allowedRoles...) {
		return models.WorkspaceMember{}, repositories.ErrForbidden
	}
	return m, nil
}

