package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/infra/database/postgres/model"
)

type Tasks struct {
	DB *gorm.DB
}

func (r Tasks) WithTx(ctx context.Context, fn func(r Tasks) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(Tasks{DB: tx})
	})
}

func (r Tasks) ListByProjectIDs(ctx context.Context, projectIDs []string) ([]entity.Task, error) {
	var rows []model.Task
	if len(projectIDs) == 0 {
		return nil, nil
	}
	if err := r.DB.WithContext(ctx).Find(&rows, "project_id IN ?", projectIDs).Error; err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	out := make([]entity.Task, 0, len(rows))
	for _, t := range rows {
		out = append(out, taskToEntity(t))
	}
	return out, nil
}

func (r Tasks) ListByProject(ctx context.Context, projectID string) ([]entity.Task, error) {
	var rows []model.Task
	if err := r.DB.WithContext(ctx).Order("column_id asc, `order` asc").Find(&rows, "project_id = ?", projectID).Error; err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	out := make([]entity.Task, 0, len(rows))
	for _, t := range rows {
		out = append(out, taskToEntity(t))
	}
	return out, nil
}

func (r Tasks) ListByColumnExcluding(ctx context.Context, columnID string, excludeTaskID string) ([]entity.Task, error) {
	var rows []model.Task
	if err := r.DB.WithContext(ctx).
		Order("`order` asc").
		Find(&rows, "column_id = ? AND id <> ?", columnID, excludeTaskID).Error; err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	out := make([]entity.Task, 0, len(rows))
	for _, t := range rows {
		out = append(out, taskToEntity(t))
	}
	return out, nil
}

func (r Tasks) CountInColumn(ctx context.Context, columnID string) (int64, error) {
	var n int64
	if err := r.DB.WithContext(ctx).Model(&model.Task{}).Where("column_id = ?", columnID).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count tasks: %w", err)
	}
	return n, nil
}

func (r Tasks) GetByID(ctx context.Context, id string) (entity.Task, error) {
	var t model.Task
	err := r.DB.WithContext(ctx).First(&t, "id = ?", id).Error
	if err == nil {
		return taskToEntity(t), nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Task{}, domain.ErrNotFound
	}
	return entity.Task{}, fmt.Errorf("get task: %w", err)
}

func (r Tasks) Create(ctx context.Context, t entity.Task) (entity.Task, error) {
	row := taskFromEntity(t)
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return entity.Task{}, fmt.Errorf("create task: %w", err)
	}
	return taskToEntity(row), nil
}

func (r Tasks) Update(ctx context.Context, id string, patch map[string]any) (entity.Task, error) {
	res := r.DB.WithContext(ctx).Model(&model.Task{}).Where("id = ?", id).Updates(patch)
	if res.Error != nil {
		return entity.Task{}, fmt.Errorf("update task: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return entity.Task{}, domain.ErrNotFound
	}
	return r.GetByID(ctx, id)
}

func (r Tasks) Delete(ctx context.Context, id string) error {
	res := r.DB.WithContext(ctx).Delete(&model.Task{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete task: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r Tasks) UpdateOrders(ctx context.Context, updates []entity.TaskOrderUpdate) error {
	for _, u := range updates {
		if err := r.DB.WithContext(ctx).
			Model(&model.Task{}).
			Where("id = ?", u.ID).
			Updates(map[string]any{
				"column_id":   u.ColumnID,
				"order":       u.Order,
				"updated_at":  u.UpdatedAt,
			}).Error; err != nil {
			return fmt.Errorf("update task order: %w", err)
		}
	}
	return nil
}

func (r Tasks) Move(ctx context.Context, taskID, toColumnID string, toIndex int) error {
	now := time.Now().UTC()
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txTasks := Tasks{DB: tx}

		var moving model.Task
		if err := tx.First(&moving, "id = ?", taskID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrNotFound
			}
			return fmt.Errorf("get task: %w", err)
		}

		fromColumnID := moving.ColumnID

		source, err := txTasks.ListByColumnExcluding(ctx, fromColumnID, taskID)
		if err != nil {
			return err
		}

		target := source
		if fromColumnID != toColumnID {
			target, err = txTasks.ListByColumnExcluding(ctx, toColumnID, taskID)
			if err != nil {
				return err
			}
		}

		if toIndex < 0 {
			toIndex = 0
		}
		if toIndex > len(target) {
			toIndex = len(target)
		}

		moved := taskToEntity(moving)
		moved.ColumnID = toColumnID

		inserted := make([]entity.Task, 0, len(target)+1)
		inserted = append(inserted, target[:toIndex]...)
		inserted = append(inserted, moved)
		inserted = append(inserted, target[toIndex:]...)

		updates := make([]entity.TaskOrderUpdate, 0, len(inserted))
		for i := range inserted {
			updates = append(updates, entity.TaskOrderUpdate{
				ID:        inserted[i].ID,
				ColumnID:  inserted[i].ColumnID,
				Order:     i,
				UpdatedAt: now,
			})
		}
		if err := txTasks.UpdateOrders(ctx, updates); err != nil {
			return err
		}

		if fromColumnID != toColumnID {
			srcUpdates := make([]entity.TaskOrderUpdate, 0, len(source))
			for i := range source {
				srcUpdates = append(srcUpdates, entity.TaskOrderUpdate{
					ID:        source[i].ID,
					ColumnID:  fromColumnID,
					Order:     i,
					UpdatedAt: now,
				})
			}
			if err := txTasks.UpdateOrders(ctx, srcUpdates); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r Tasks) UnassignForWorkspace(ctx context.Context, workspaceID, userID string, updatedAt time.Time) error {
	if err := r.DB.WithContext(ctx).Exec(`
		UPDATE tasks
		SET assignee_id = NULL, updated_at = ?
		WHERE assignee_id = ?
		AND project_id IN (SELECT id FROM projects WHERE workspace_id = ?)
	`, updatedAt, userID, workspaceID).Error; err != nil {
		return fmt.Errorf("unassign tasks: %w", err)
	}
	return nil
}

func (r Tasks) ReassignForWorkspace(ctx context.Context, workspaceID, fromUserID, toUserID string, updatedAt time.Time) error {
	if err := r.DB.WithContext(ctx).Exec(`
		UPDATE tasks
		SET assignee_id = ?, updated_at = ?
		WHERE assignee_id = ?
		AND project_id IN (SELECT id FROM projects WHERE workspace_id = ?)
	`, toUserID, updatedAt, fromUserID, workspaceID).Error; err != nil {
		return fmt.Errorf("reassign tasks: %w", err)
	}
	return nil
}
