package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (s AppService) ListLabels(ctx context.Context, ownerID, workspaceID string) ([]dto.Label, error) {
	ws, err := s.Repo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if ws.OwnerID != ownerID {
		return nil, repositories.ErrForbidden
	}
	ls, err := s.Repo.ListLabelsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.Label, 0, len(ls))
	for _, l := range ls {
		out = append(out, dto.Label{ID: l.ID, WorkspaceID: l.WorkspaceID, Name: l.Name, Color: l.Color})
	}
	return out, nil
}

func (s AppService) CreateLabel(ctx context.Context, ownerID, workspaceID, name, color string) (dto.Label, error) {
	ws, err := s.Repo.GetWorkspaceByID(ctx, workspaceID)
	if err != nil {
		return dto.Label{}, err
	}
	if ws.OwnerID != ownerID {
		return dto.Label{}, repositories.ErrForbidden
	}
	now := time.Now().UTC()
	l := models.Label{
		ID:          uuid.NewString(),
		WorkspaceID: workspaceID,
		Name:        name,
		Color:       color,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	created, err := s.Repo.CreateLabel(ctx, l)
	if err != nil {
		return dto.Label{}, err
	}
	return dto.Label{ID: created.ID, WorkspaceID: created.WorkspaceID, Name: created.Name, Color: created.Color}, nil
}

func (s AppService) UpdateLabel(ctx context.Context, ownerID, labelID string, patch map[string]any) (dto.Label, error) {
	l, err := s.Repo.GetLabelByID(ctx, labelID)
	if err != nil {
		return dto.Label{}, err
	}
	ws, err := s.Repo.GetWorkspaceByID(ctx, l.WorkspaceID)
	if err != nil {
		return dto.Label{}, err
	}
	if ws.OwnerID != ownerID {
		return dto.Label{}, repositories.ErrForbidden
	}
	patch["updated_at"] = time.Now().UTC()
	updated, err := s.Repo.UpdateLabel(ctx, labelID, patch)
	if err != nil {
		return dto.Label{}, err
	}
	return dto.Label{ID: updated.ID, WorkspaceID: updated.WorkspaceID, Name: updated.Name, Color: updated.Color}, nil
}

func (s AppService) DeleteLabel(ctx context.Context, ownerID, labelID string) error {
	l, err := s.Repo.GetLabelByID(ctx, labelID)
	if err != nil {
		return err
	}
	ws, err := s.Repo.GetWorkspaceByID(ctx, l.WorkspaceID)
	if err != nil {
		return err
	}
	if ws.OwnerID != ownerID {
		return repositories.ErrForbidden
	}
	return s.Repo.DeleteLabel(ctx, labelID)
}
