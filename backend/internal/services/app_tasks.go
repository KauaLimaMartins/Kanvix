package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
	"kanvix/backend/internal/repositories"
)

func (s AppService) ListTasks(ctx context.Context, ownerID, projectID string) ([]dto.Task, error) {
	p, err := s.Repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	ws, err := s.Repo.GetWorkspaceByID(ctx, p.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if ws.OwnerID != ownerID {
		return nil, repositories.ErrForbidden
	}

	tasks, err := s.Repo.ListTasksByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	taskIDs := make([]string, 0, len(tasks))
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}

	taskToLabels := map[string][]string{}
	if len(taskIDs) > 0 {
		var joins []models.TaskLabel
		if err := s.Repo.DB.WithContext(ctx).Find(&joins, "task_id IN ?", taskIDs).Error; err != nil {
			return nil, fmt.Errorf("list task labels: %w", err)
		}
		for _, j := range joins {
			taskToLabels[j.TaskID] = append(taskToLabels[j.TaskID], j.LabelID)
		}
	}

	out := make([]dto.Task, 0, len(tasks))
	for _, t := range tasks {
		out = append(out, dto.Task{
			ID:          t.ID,
			ProjectID:   t.ProjectID,
			ColumnID:    t.ColumnID,
			Title:       t.Title,
			Description: t.Description,
			Labels:      taskToLabels[t.ID],
			DueDate:     t.DueDate,
			AssigneeID:  t.AssigneeID,
			Order:       t.Order,
			CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}
	return out, nil
}

func (s AppService) GetTask(ctx context.Context, ownerID, taskID string) (dto.Task, error) {
	t, err := s.Repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return dto.Task{}, err
	}
	p, err := s.Repo.GetProjectByID(ctx, t.ProjectID)
	if err != nil {
		return dto.Task{}, err
	}
	ws, err := s.Repo.GetWorkspaceByID(ctx, p.WorkspaceID)
	if err != nil {
		return dto.Task{}, err
	}
	if ws.OwnerID != ownerID {
		return dto.Task{}, repositories.ErrForbidden
	}
	var joins []models.TaskLabel
	if err := s.Repo.DB.WithContext(ctx).Find(&joins, "task_id = ?", t.ID).Error; err != nil {
		return dto.Task{}, fmt.Errorf("list task labels: %w", err)
	}
	labels := make([]string, 0, len(joins))
	for _, j := range joins {
		labels = append(labels, j.LabelID)
	}
	return dto.Task{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		ColumnID:    t.ColumnID,
		Title:       t.Title,
		Description: t.Description,
		Labels:      labels,
		DueDate:     t.DueDate,
		AssigneeID:  t.AssigneeID,
		Order:       t.Order,
		CreatedAt:   t.CreatedAt.UTC().Format(time.RFC3339Nano),
	}, nil
}

func (s AppService) CreateTask(ctx context.Context, ownerID, projectID, columnID, title string) (dto.Task, error) {
	p, err := s.Repo.GetProjectByID(ctx, projectID)
	if err != nil {
		return dto.Task{}, err
	}
	ws, err := s.Repo.GetWorkspaceByID(ctx, p.WorkspaceID)
	if err != nil {
		return dto.Task{}, err
	}
	if ws.OwnerID != ownerID {
		return dto.Task{}, repositories.ErrForbidden
	}

	n, err := s.Repo.CountTasksInColumn(ctx, columnID)
	if err != nil {
		return dto.Task{}, err
	}

	now := time.Now().UTC()
	t := models.Task{
		ID:          uuid.NewString(),
		ProjectID:   projectID,
		ColumnID:    columnID,
		Title:       title,
		Description: "",
		Order:       int(n),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	created, err := s.Repo.CreateTask(ctx, t)
	if err != nil {
		return dto.Task{}, err
	}
	return dto.Task{
		ID:          created.ID,
		ProjectID:   created.ProjectID,
		ColumnID:    created.ColumnID,
		Title:       created.Title,
		Description: created.Description,
		Labels:      []string{},
		DueDate:     created.DueDate,
		AssigneeID:  created.AssigneeID,
		Order:       created.Order,
		CreatedAt:   created.CreatedAt.UTC().Format(time.RFC3339Nano),
	}, nil
}

func (s AppService) UpdateTask(ctx context.Context, ownerID, taskID string, patch map[string]any, labels *[]string) (dto.Task, error) {
	t, err := s.Repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return dto.Task{}, err
	}
	p, err := s.Repo.GetProjectByID(ctx, t.ProjectID)
	if err != nil {
		return dto.Task{}, err
	}
	ws, err := s.Repo.GetWorkspaceByID(ctx, p.WorkspaceID)
	if err != nil {
		return dto.Task{}, err
	}
	if ws.OwnerID != ownerID {
		return dto.Task{}, repositories.ErrForbidden
	}

	patch["updated_at"] = time.Now().UTC()
	updated, err := s.Repo.UpdateTask(ctx, taskID, patch)
	if err != nil {
		return dto.Task{}, err
	}

	if labels != nil {
		if err := s.Repo.ReplaceTaskLabels(ctx, taskID, *labels); err != nil {
			return dto.Task{}, fmt.Errorf("replace labels: %w", err)
		}
	}

	return s.GetTask(ctx, ownerID, updated.ID)
}

func (s AppService) DeleteTask(ctx context.Context, ownerID, taskID string) error {
	t, err := s.Repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}
	p, err := s.Repo.GetProjectByID(ctx, t.ProjectID)
	if err != nil {
		return err
	}
	ws, err := s.Repo.GetWorkspaceByID(ctx, p.WorkspaceID)
	if err != nil {
		return err
	}
	if ws.OwnerID != ownerID {
		return repositories.ErrForbidden
	}
	return s.Repo.DeleteTask(ctx, taskID)
}

func (s AppService) MoveTask(ctx context.Context, ownerID, taskID, toColumnID string, toIndex int) error {
	t, err := s.Repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return err
	}
	p, err := s.Repo.GetProjectByID(ctx, t.ProjectID)
	if err != nil {
		return err
	}
	ws, err := s.Repo.GetWorkspaceByID(ctx, p.WorkspaceID)
	if err != nil {
		return err
	}
	if ws.OwnerID != ownerID {
		return repositories.ErrForbidden
	}

	fromColumnID := t.ColumnID
	return s.Repo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var source []models.Task
		if err := tx.Order("`order` asc").Find(&source, "column_id = ? AND id <> ?", fromColumnID, taskID).Error; err != nil {
			return err
		}

		target := source
		if fromColumnID != toColumnID {
			if err := tx.Order("`order` asc").Find(&target, "column_id = ? AND id <> ?", toColumnID, taskID).Error; err != nil {
				return err
			}
		}

		clamped := toIndex
		if clamped < 0 {
			clamped = 0
		}
		if clamped > len(target) {
			clamped = len(target)
		}

		inserted := make([]models.Task, 0, len(target)+1)
		inserted = append(inserted, target[:clamped]...)
		t.ColumnID = toColumnID
		inserted = append(inserted, t)
		inserted = append(inserted, target[clamped:]...)

		now := time.Now().UTC()
		for i := range inserted {
			if err := tx.Model(&models.Task{}).
				Where("id = ?", inserted[i].ID).
				Updates(map[string]any{
					"column_id":  inserted[i].ColumnID,
					"order":      i,
					"updated_at": now,
				}).Error; err != nil {
				return err
			}
		}

		if fromColumnID != toColumnID {
			for i := range source {
				if err := tx.Model(&models.Task{}).
					Where("id = ?", source[i].ID).
					Updates(map[string]any{
						"order":      i,
						"updated_at": now,
					}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}
