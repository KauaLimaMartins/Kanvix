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

type Columns struct {
	DB *gorm.DB
}

func (r Columns) ListByProjectIDs(ctx context.Context, projectIDs []string) ([]entity.Column, error) {
	var rows []model.Column
	if len(projectIDs) == 0 {
		return nil, nil
	}
	if err := r.DB.WithContext(ctx).Find(&rows, "project_id IN ?", projectIDs).Error; err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}
	out := make([]entity.Column, 0, len(rows))
	for _, c := range rows {
		out = append(out, columnToEntity(c))
	}
	return out, nil
}

func (r Columns) ListByProject(ctx context.Context, projectID string) ([]entity.Column, error) {
	var rows []model.Column
	if err := r.DB.WithContext(ctx).Order("`order` asc").Find(&rows, "project_id = ?", projectID).Error; err != nil {
		return nil, fmt.Errorf("list columns: %w", err)
	}
	out := make([]entity.Column, 0, len(rows))
	for _, c := range rows {
		out = append(out, columnToEntity(c))
	}
	return out, nil
}

func (r Columns) CountByProject(ctx context.Context, projectID string) (int64, error) {
	var n int64
	if err := r.DB.WithContext(ctx).Model(&model.Column{}).Where("project_id = ?", projectID).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count columns: %w", err)
	}
	return n, nil
}

func (r Columns) GetByID(ctx context.Context, id string) (entity.Column, error) {
	var c model.Column
	err := r.DB.WithContext(ctx).First(&c, "id = ?", id).Error
	if err == nil {
		return columnToEntity(c), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Column{}, domain.ErrNotFound
	}
	return entity.Column{}, fmt.Errorf("get column: %w", err)
}

func (r Columns) Create(ctx context.Context, c entity.Column) (entity.Column, error) {
	row := columnFromEntity(c)
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return entity.Column{}, fmt.Errorf("create column: %w", err)
	}
	return columnToEntity(row), nil
}

func (r Columns) Update(ctx context.Context, id string, patch map[string]any) (entity.Column, error) {
	res := r.DB.WithContext(ctx).Model(&model.Column{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return entity.Column{}, fmt.Errorf("update column: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.Column{}, domain.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r Columns) Delete(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&model.Column{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete column: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
