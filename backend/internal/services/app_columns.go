package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
)

func (s AppService) ListColumns(ctx context.Context, ownerID, projectID string) ([]dto.Column, error) {
	p, err := s.Repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin", "member"); err != nil {
		return nil, err
	}
	cols, err := s.Repo.ListColumnsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.Column, 0, len(cols))
	for _, c := range cols {
		out = append(out, dto.Column{ID: c.ID, ProjectID: c.ProjectID, Name: c.Name, Order: c.Order})
	}
	return out, nil
}

func (s AppService) CreateColumn(ctx context.Context, ownerID, projectID, name string) (dto.Column, error) {
	p, err := s.Repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return dto.Column{}, err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin"); err != nil {
		return dto.Column{}, err
	}
	n, err := s.Repo.CountColumnsInProject(ctx, projectID)
	if err != nil {
		return dto.Column{}, err
	}
	now := time.Now().UTC()
	col := models.Column{
		ID:        uuid.NewString(),
		ProjectID: projectID,
		Name:      name,
		Order:     int(n),
		CreatedAt: now,
		UpdatedAt: now,
	}
	created, err := s.Repo.CreateColumn(ctx, col)
	if err != nil {
		return dto.Column{}, err
	}
	return dto.Column{ID: created.ID, ProjectID: created.ProjectID, Name: created.Name, Order: created.Order}, nil
}

func (s AppService) UpdateColumn(ctx context.Context, ownerID, columnID string, patch map[string]any) (dto.Column, error) {
	col, err := s.Repo.GetColumnByID(ctx, columnID)
	if err != nil {
		return dto.Column{}, err
	}
	p, err := s.Repo.GetProjectByID(ctx, col.ProjectID)
	if err != nil {
		return dto.Column{}, err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin"); err != nil {
		return dto.Column{}, err
	}
	patch["updated_at"] = time.Now().UTC()
	updated, err := s.Repo.UpdateColumn(ctx, columnID, patch)
	if err != nil {
		return dto.Column{}, err
	}
	return dto.Column{ID: updated.ID, ProjectID: updated.ProjectID, Name: updated.Name, Order: updated.Order}, nil
}

func (s AppService) DeleteColumn(ctx context.Context, ownerID, columnID string) error {
	col, err := s.Repo.GetColumnByID(ctx, columnID)
	if err != nil {
		return err
	}
	p, err := s.Repo.GetProjectByID(ctx, col.ProjectID)
	if err != nil {
		return err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin"); err != nil {
		return err
	}
	return s.Repo.DeleteColumn(ctx, columnID)
}
