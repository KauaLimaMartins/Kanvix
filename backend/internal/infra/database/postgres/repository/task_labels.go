package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"kanvix/backend/internal/infra/database/postgres/model"
)

type TaskLabels struct {
	DB *gorm.DB
}

func (r TaskLabels) ListByTaskIDs(ctx context.Context, taskIDs []string) ([]model.TaskLabel, error) {
	var rows []model.TaskLabel
	if len(taskIDs) == 0 {
		return nil, nil
	}
	if err := r.DB.WithContext(ctx).Find(&rows, "task_id IN ?", taskIDs).Error; err != nil {
		return nil, fmt.Errorf("list task labels: %w", err)
	}
	return rows, nil
}

func (r TaskLabels) ListLabelIDsByTaskIDs(ctx context.Context, taskIDs []string) (map[string][]string, error) {
	rows, err := r.ListByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}
	out := map[string][]string{}
	for _, r := range rows {
		out[r.TaskID] = append(out[r.TaskID], r.LabelID)
	}
	return out, nil
}

func (r TaskLabels) ListByTaskID(ctx context.Context, taskID string) ([]model.TaskLabel, error) {
	var rows []model.TaskLabel
	if err := r.DB.WithContext(ctx).Find(&rows, "task_id = ?", taskID).Error; err != nil {
		return nil, fmt.Errorf("list task labels: %w", err)
	}
	return rows, nil
}

func (r TaskLabels) ListLabelIDsByTaskID(ctx context.Context, taskID string) ([]string, error) {
	rows, err := r.ListByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.LabelID)
	}
	return out, nil
}

func (r TaskLabels) Replace(ctx context.Context, taskID string, labelIDs []string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.TaskLabel{}, "task_id = ?", taskID).Error; err != nil {
			return fmt.Errorf("delete task labels: %w", err)
		}
		now := time.Now().UTC()
		for _, lid := range labelIDs {
			if lid == "" {
				continue
			}
			row := taskLabelFrom(taskID, lid, now)
			if err := tx.Create(&row).Error; err != nil {
				return fmt.Errorf("create task label: %w", err)
			}
		}
		return nil
	})
}
