package update

import (
	"context"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID    string
	SubtaskID string
	Title     *string
	Done      *bool
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
	GetByID(ctx context.Context, id string) (entity.Subtask, error)
	Update(ctx context.Context, id string, patch map[string]any) (entity.Subtask, error)
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
	st, err := uc.Subtasks.GetByID(ctx, in.SubtaskID)
	if err != nil {
		return Out{}, err
	}
	t, err := uc.Tasks.GetByID(ctx, st.TaskID)
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

	patch := map[string]any{"updated_at": time.Now().UTC()}
	if in.Title != nil {
		patch["title"] = *in.Title
	}
	if in.Done != nil {
		patch["done"] = *in.Done
	}

	updated, err := uc.Subtasks.Update(ctx, in.SubtaskID, patch)
	if err != nil {
		return Out{}, err
	}
	return Out{Subtask: updated}, nil
}

