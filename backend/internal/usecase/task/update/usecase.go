package update

import (
	"context"
	"fmt"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID string
	TaskID string

	Patch  map[string]any
	Labels *[]string
}

type Out struct {
	Task   entity.Task
	Labels []string
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type taskRepo interface {
	GetByID(ctx context.Context, id string) (entity.Task, error)
	Update(ctx context.Context, id string, patch map[string]any) (entity.Task, error)
}

type taskLabelsRepo interface {
	Replace(ctx context.Context, taskID string, labelIDs []string) error
	ListLabelIDsByTaskID(ctx context.Context, taskID string) ([]string, error)
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
	TaskLabels taskLabelsRepo
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

	patch := in.Patch
	if patch == nil {
		patch = map[string]any{}
	}
	patch["updated_at"] = time.Now().UTC()
	updated, err := uc.Tasks.Update(ctx, in.TaskID, patch)
	if err != nil {
		return Out{}, err
	}

	if in.Labels != nil {
		if err := uc.TaskLabels.Replace(ctx, in.TaskID, *in.Labels); err != nil {
			return Out{}, fmt.Errorf("replace labels: %w", err)
		}
	}

	lbls, err := uc.TaskLabels.ListLabelIDsByTaskID(ctx, in.TaskID)
	if err != nil {
		return Out{}, err
	}
	return Out{Task: updated, Labels: lbls}, nil
}
