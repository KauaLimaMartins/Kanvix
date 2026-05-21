package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (s AppService) ListWorkspaces(ctx context.Context, ownerID string) ([]dto.Workspace, error) {
	ws, err := s.Repo.ListWorkspacesForUser(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.Workspace, 0, len(ws))
	for _, w := range ws {
		out = append(out, dto.Workspace{ID: w.ID, Name: w.Name, Icon: w.Icon, Color: w.Color})
	}
	return out, nil
}

func (s AppService) CreateWorkspace(ctx context.Context, ownerID string, name string, icon *string, color *string) (dto.Workspace, error) {
	if err := s.requireGlobalAdmin(ctx, ownerID); err != nil {
		return dto.Workspace{}, err
	}
	now := time.Now().UTC()
	w := models.Workspace{
		ID:        uuid.NewString(),
		OwnerID:   ownerID,
		Name:      name,
		Icon:      "LayoutGrid",
		Color:     "#6366f1",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if icon != nil && *icon != "" {
		w.Icon = *icon
	}
	if color != nil && *color != "" {
		w.Color = *color
	}
	created, err := s.Repo.CreateWorkspace(ctx, w)
	if err != nil {
		return dto.Workspace{}, err
	}
	if err := s.Repo.UpsertWorkspaceMember(ctx, models.WorkspaceMember{
		WorkspaceID: created.ID,
		UserID:      ownerID,
		Role:        "admin",
		CreatedAt:   now,
	}); err != nil {
		return dto.Workspace{}, err
	}
	return dto.Workspace{ID: created.ID, Name: created.Name, Icon: created.Icon, Color: created.Color}, nil
}

func (s AppService) UpdateWorkspace(ctx context.Context, ownerID, workspaceID string, patch map[string]any) (dto.Workspace, error) {
	if _, err := s.requireWorkspaceRole(ctx, ownerID, workspaceID, "admin"); err != nil {
		return dto.Workspace{}, err
	}
	patch["updated_at"] = time.Now().UTC()
	updated, err := s.Repo.UpdateWorkspace(ctx, workspaceID, patch)
	if err != nil {
		return dto.Workspace{}, err
	}
	return dto.Workspace{ID: updated.ID, Name: updated.Name, Icon: updated.Icon, Color: updated.Color}, nil
}

func (s AppService) DeleteWorkspace(ctx context.Context, ownerID, workspaceID string) error {
	w, err := s.Repo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return err
	}
	if w.OwnerID != ownerID {
		return repositories.ErrForbidden
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, workspaceID, "admin"); err != nil {
		return err
	}
	return s.Repo.DeleteWorkspace(ctx, workspaceID)
}
