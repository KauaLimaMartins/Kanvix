package move

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID     string
	TaskID     string
	ToColumnID string
	ToIndex    int
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type taskRepo interface {
	GetByID(ctx context.Context, id string) (entity.Task, error)
	Move(ctx context.Context, taskID, toColumnID string, toIndex int) error
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
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) error {
	t, err := uc.Tasks.GetByID(ctx, in.TaskID)
	if err != nil {
		return err
	}
	p, err := uc.Projects.GetByID(ctx, t.ProjectID)
	if err != nil {
		return err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin", "member"); err != nil {
		return err
	}

	return uc.Tasks.Move(ctx, t.ID, in.ToColumnID, in.ToIndex)
}
