package list

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID    string
	ProjectID string
}

type Out struct {
	Tasks  []entity.Task
	Labels map[string][]string
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type taskRepo interface {
	ListByProject(ctx context.Context, projectID string) ([]entity.Task, error)
}

type taskLabelsRepo interface {
	ListLabelIDsByTaskIDs(ctx context.Context, taskIDs []string) (map[string][]string, error)
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
	p, err := uc.Projects.GetByID(ctx, in.ProjectID)
	if err != nil {
		return Out{}, err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin", "member"); err != nil {
		return Out{}, err
	}

	tasks, err := uc.Tasks.ListByProject(ctx, in.ProjectID)
	if err != nil {
		return Out{}, err
	}

	taskIDs := make([]string, 0, len(tasks))
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}
	taskToLabels, err := uc.TaskLabels.ListLabelIDsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return Out{}, err
	}

	return Out{Tasks: tasks, Labels: taskToLabels}, nil
}
