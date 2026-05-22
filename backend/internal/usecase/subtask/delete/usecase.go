package delete

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID    string
	SubtaskID string
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type taskRepo interface {
	GetByID(ctx context.Context, id string) (entity.Task, error)
}

type subtaskRepo interface {
	GetByID(ctx context.Context, id string) (entity.Subtask, error)
	Delete(ctx context.Context, id string) error
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

func (uc UseCase) Execute(ctx context.Context, in In) error {
	st, err := uc.Subtasks.GetByID(ctx, in.SubtaskID)
	if err != nil {
		return err
	}
	t, err := uc.Tasks.GetByID(ctx, st.TaskID)
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
	return uc.Subtasks.Delete(ctx, in.SubtaskID)
}

