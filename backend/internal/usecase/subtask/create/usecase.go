package create

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID string
	TaskID string
	Title  string
}

type Out struct {
	Subtask entity.Subtask
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type taskRepo interface {
	GetByID(ctx context.Context, id string) (entity.Task, error)
}

type subtaskRepo interface {
	Create(ctx context.Context, s entity.Subtask) (entity.Subtask, error)
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Projects   projectRepo
	Tasks      taskRepo
	Subtasks   subtaskRepo
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	t, err := uc.Tasks.GetByID(ctx, in.TaskID)
	if err != nil {
		return Out{}, err
	}
	p, err := uc.Projects.GetByID(ctx, t.ProjectID)
	if err != nil {
		return Out{}, err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin", "member"); err != nil {
		return Out{}, err
	}

	now := time.Now().UTC()
	s, err := uc.Subtasks.Create(ctx, entity.Subtask{
		ID:        uuid.NewString(),
		TaskID:    in.TaskID,
		Title:     in.Title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return Out{}, err
	}
	return Out{Subtask: s}, nil
}

