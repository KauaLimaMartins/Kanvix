package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/infra/database/postgres/model"
)

type Projects struct {
	DB *gorm.DB
}

func (r Projects) ListByWorkspaceIDs(ctx context.Context, workspaceIDs []string) ([]entity.Project, error) {
	var rows []model.Project
	if len(workspaceIDs) == 0 {
		return nil, nil
	}
	if err := r.DB.WithContext(ctx).Find(&rows, "workspace_id IN ?", workspaceIDs).Error; err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	out := make([]entity.Project, 0, len(rows))
	for _, p := range rows {
		out = append(out, projectToEntity(p))
	}
	return out, nil
}

func (r Projects) ListByWorkspace(ctx context.Context, workspaceID string) ([]entity.Project, error) {
	var rows []model.Project
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&rows, "workspace_id = ?", workspaceID).Error; err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	out := make([]entity.Project, 0, len(rows))
	for _, p := range rows {
		out = append(out, projectToEntity(p))
	}
	return out, nil
}

func (r Projects) GetByID(ctx context.Context, id string) (entity.Project, error) {
	var p model.Project
	err := r.DB.WithContext(ctx).First(&p, "id = ?", id).Error
	if err == nil {
		return projectToEntity(p), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Project{}, domain.ErrNotFound
	}
	return entity.Project{}, fmt.Errorf("get project: %w", err)
}

func (r Projects) Create(ctx context.Context, p entity.Project) (entity.Project, error) {
	row := projectFromEntity(p)
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return entity.Project{}, fmt.Errorf("create project: %w", err)
	}
	return projectToEntity(row), nil
}

func (r Projects) Update(ctx context.Context, id string, patch map[string]any) (entity.Project, error) {
	res := r.DB.WithContext(ctx).Model(&model.Project{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return entity.Project{}, fmt.Errorf("update project: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.Project{}, domain.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r Projects) Delete(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&model.Project{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete project: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
