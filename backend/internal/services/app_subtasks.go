package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/http/dto"
	"kanvix/backend/internal/models"
)

func (s AppService) ListSubtasks(ctx context.Context, ownerID, taskID string) ([]dto.Subtask, error) {
	t, err := s.Repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	p, err := s.Repo.GetProjectByID(ctx, t.ProjectID)
	if err != nil {
		return nil, err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin", "member"); err != nil {
		return nil, err
	}

	ss, err := s.Repo.ListSubtasksByTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.Subtask, 0, len(ss))
	for _, st := range ss {
		out = append(out, dto.Subtask{
			ID:        st.ID,
			TaskID:    st.TaskID,
			Title:     st.Title,
			Done:      st.Done,
			CreatedAt: st.CreatedAt.UTC().Format(time.RFC3339Nano),
		})
	}
	return out, nil
}

func (s AppService) CreateSubtask(ctx context.Context, ownerID, taskID, title string) (dto.Subtask, error) {
	t, err := s.Repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return dto.Subtask{}, err
	}
	p, err := s.Repo.GetProjectByID(ctx, t.ProjectID)
	if err != nil {
		return dto.Subtask{}, err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin", "member"); err != nil {
		return dto.Subtask{}, err
	}

	now := time.Now().UTC()
	st := models.Subtask{
		ID:        uuid.NewString(),
		TaskID:    taskID,
		Title:     title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	created, err := s.Repo.CreateSubtask(ctx, st)
	if err != nil {
		return dto.Subtask{}, err
	}
	return dto.Subtask{
		ID:        created.ID,
		TaskID:    created.TaskID,
		Title:     created.Title,
		Done:      created.Done,
		CreatedAt: created.CreatedAt.UTC().Format(time.RFC3339Nano),
	}, nil
}

func (s AppService) SetSubtaskDone(ctx context.Context, ownerID, subtaskID string, done bool) (dto.Subtask, error) {
	st, err := s.Repo.GetSubtaskByID(ctx, subtaskID)
	if err != nil {
		return dto.Subtask{}, err
	}
	t, err := s.Repo.GetTaskByID(ctx, st.TaskID)
	if err != nil {
		return dto.Subtask{}, err
	}
	p, err := s.Repo.GetProjectByID(ctx, t.ProjectID)
	if err != nil {
		return dto.Subtask{}, err
	}
	if _, err := s.requireWorkspaceRole(ctx, ownerID, p.WorkspaceID, "admin", "member"); err != nil {
		return dto.Subtask{}, err
	}

	updated, err := s.Repo.UpdateSubtaskDone(ctx, subtaskID, done, time.Now().UTC())
	if err != nil {
		return dto.Subtask{}, err
	}
	return dto.Subtask{
		ID:        updated.ID,
		TaskID:    updated.TaskID,
		Title:     updated.Title,
		Done:      updated.Done,
		CreatedAt: updated.CreatedAt.UTC().Format(time.RFC3339Nano),
	}, nil
}
