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

type Workspaces struct {
	DB *gorm.DB
}

func (r Workspaces) ListByIDs(ctx context.Context, ids []string) ([]entity.Workspace, error) {
	var rows []model.Workspace
	if len(ids) == 0 {
		return nil, nil
	}
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&rows, "id IN ?", ids).Error; err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	out := make([]entity.Workspace, 0, len(rows))
	for _, w := range rows {
		out = append(out, workspaceToEntity(w))
	}
	return out, nil
}

func (r Workspaces) ListByOwner(ctx context.Context, ownerID string) ([]entity.Workspace, error) {
	var rows []model.Workspace
	if err := r.DB.WithContext(ctx).Order("created_at asc").Find(&rows, "owner_id = ?", ownerID).Error; err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	out := make([]entity.Workspace, 0, len(rows))
	for _, w := range rows {
		out = append(out, workspaceToEntity(w))
	}
	return out, nil
}

func (r Workspaces) GetByID(ctx context.Context, id string) (entity.Workspace, error) {
	var w model.Workspace
	err := r.DB.WithContext(ctx).First(&w, "id = ?", id).Error
	if err == nil {
		return workspaceToEntity(w), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Workspace{}, domain.ErrNotFound
	}
	return entity.Workspace{}, fmt.Errorf("get workspace: %w", err)
}

func (r Workspaces) Create(ctx context.Context, w entity.Workspace) (entity.Workspace, error) {
	row := workspaceFromEntity(w)
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return entity.Workspace{}, fmt.Errorf("create workspace: %w", err)
	}
	return workspaceToEntity(row), nil
}

func (r Workspaces) Update(ctx context.Context, id string, patch map[string]any) (entity.Workspace, error) {
	res := r.DB.WithContext(ctx).Model(&model.Workspace{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return entity.Workspace{}, fmt.Errorf("update workspace: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.Workspace{}, domain.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r Workspaces) Delete(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&model.Workspace{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete workspace: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
