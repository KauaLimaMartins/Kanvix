package list

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID string
	TaskID string
}

type Out struct {
	Subtasks []entity.Subtask
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type taskRepo interface {
	GetByID(ctx context.Context, id string) (entity.Task, error)
}

type subtaskRepo interface {
	ListByTask(ctx context.Context, taskID string) ([]entity.Subtask, error)
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
	ss, err := uc.Subtasks.ListByTask(ctx, in.TaskID)
	if err != nil {
		return Out{}, err
	}
	return Out{Subtasks: ss}, nil
}

