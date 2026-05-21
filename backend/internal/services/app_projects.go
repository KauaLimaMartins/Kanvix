package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
)

func (s AppService) ListProjects(ctx context.Context, ownerID, workspaceID string) ([]dto.Project, error) {
	if _, err := s.requireWorkspaceRole(ctx, ownerID, workspaceID, "admin", "member"); err != nil {
		return nil, err
	}
	ps, err := s.Repo.ListProjectsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.Project, 0, len(ps))
	for _, p := range ps {
		out = append(out, dto.Project{ID: p.ID, WorkspaceID: p.WorkspaceID, Name: p.Name, Description: p.Description})
	}
	return out, nil
}

func (s AppService) CreateProject(ctx context.Context, ownerID, workspaceID, name string, description *string) (dto.Project, []dto.Column, error) {
	if _, err := s.requireWorkspaceRole(ctx, ownerID, workspaceID, "admin"); err != nil {
		return dto.Project{}, nil, err
	}

	now := time.Now().UTC()
	p := models.Project{
		ID:          uuid.NewString(),
		WorkspaceID: workspaceID,
		Name:        name,
		Description: "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if description != nil {
		p.Description = *description
	}
	created, err := s.Repo.CreateProject(ctx, p)
	if err != nil {
		return dto.Project{}, nil, err
	}

	colNames := []string{"To do", "In progress", "Done"}
	colsOut := make([]dto.Column, 0, len(colNames))
	for i, n := range colNames {
		col := models.Column{
			ID:        uuid.NewString(),
			ProjectID: created.ID,
			Name:      n,
			Order:     i,
			CreatedAt: now,
			UpdatedAt: now,
		}
		col, err = s.Repo.CreateColumn(ctx, col)
		if err != nil {
			return dto.Project{}, nil, err
		}
		colsOut = append(colsOut, dto.Column{ID: col.ID, ProjectID: col.ProjectID, Name: col.Name, Order: col.Order})
	}

	return dto.Project{ID: created.ID, WorkspaceID: created.WorkspaceID, Name: created.Name, Description: created.Description}, colsOut, nil
}

func (s AppService) UpdateProject(ctx context.Context, ownerID, projectID string, patch map[string]any) (dto.Project, error) {
	p, err := s.Repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return dto.Project{}, err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin"); err != nil {
		return dto.Project{}, err
	}
	patch["updated_at"] = time.Now().UTC()
	updated, err := s.Repo.UpdateProject(ctx, projectID, patch)
	if err != nil {
		return dto.Project{}, err
	}
	return dto.Project{ID: updated.ID, WorkspaceID: updated.WorkspaceID, Name: updated.Name, Description: updated.Description}, nil
}

func (s AppService) DeleteProject(ctx context.Context, ownerID, projectID string) error {
	p, err := s.Repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin"); err != nil {
		return err
	}
	return s.Repo.DeleteProject(ctx, projectID)
}
