package get

import (
	"context"
	"fmt"
	"kanvix/backend/internal/domain/entity"
)

type Out struct {
	User       entity.User
	Workspaces []entity.Workspace
	RoleByWS   map[string]string
	Projects   []entity.Project
	Columns    []entity.Column
	Tasks      []entity.Task
	TaskLabels map[string][]string
	Labels     []entity.Label
	Users      []entity.User
}

type userRepo interface {
	GetByID(ctx context.Context, id string) (entity.User, error)
	ListDistinctByWorkspaceIDs(ctx context.Context, workspaceIDs []string) ([]entity.User, error)
}

type membershipRepo interface {
	ListByUser(ctx context.Context, userID string) ([]entity.WorkspaceMember, error)
}

type workspaceRepo interface {
	ListByIDs(ctx context.Context, ids []string) ([]entity.Workspace, error)
}

type projectRepo interface {
	ListByWorkspaceIDs(ctx context.Context, workspaceIDs []string) ([]entity.Project, error)
}

type columnRepo interface {
	ListByProjectIDs(ctx context.Context, projectIDs []string) ([]entity.Column, error)
}

type taskRepo interface {
	ListByProjectIDs(ctx context.Context, projectIDs []string) ([]entity.Task, error)
}

type labelRepo interface {
	ListByWorkspaceIDs(ctx context.Context, workspaceIDs []string) ([]entity.Label, error)
}

type taskLabelsRepo interface {
	ListLabelIDsByTaskIDs(ctx context.Context, taskIDs []string) (map[string][]string, error)
}

type UseCase struct {
	Users      userRepo
	Members    membershipRepo
	Workspaces workspaceRepo
	Projects   projectRepo
	Columns    columnRepo
	Tasks      taskRepo
	Labels     labelRepo
	TaskLabels taskLabelsRepo
}

func (uc UseCase) Execute(ctx context.Context, userID string) (Out, error) {
	u, err := uc.Users.GetByID(ctx, userID)
	if err != nil {
		return Out{}, err
	}

	memberships, err := uc.Members.ListByUser(ctx, userID)
	if err != nil {
		return Out{}, err
	}

	roleByWS := map[string]string{}
	wsIDs := make([]string, 0, len(memberships))
	for _, m := range memberships {
		wsIDs = append(wsIDs, m.WorkspaceID)
		roleByWS[m.WorkspaceID] = m.Role
	}

	workspaces, err := uc.Workspaces.ListByIDs(ctx, wsIDs)
	if err != nil {
		return Out{}, err
	}

	projects, err := uc.Projects.ListByWorkspaceIDs(ctx, wsIDs)
	if err != nil {
		return Out{}, err
	}
	labels, err := uc.Labels.ListByWorkspaceIDs(ctx, wsIDs)
	if err != nil {
		return Out{}, err
	}

	pIDs := make([]string, 0, len(projects))
	for _, p := range projects {
		pIDs = append(pIDs, p.ID)
	}

	columns, err := uc.Columns.ListByProjectIDs(ctx, pIDs)
	if err != nil {
		return Out{}, err
	}
	tasks, err := uc.Tasks.ListByProjectIDs(ctx, pIDs)
	if err != nil {
		return Out{}, err
	}

	taskIDs := make([]string, 0, len(tasks))
	for _, t := range tasks {
		taskIDs = append(taskIDs, t.ID)
	}

	taskToLabels, err := uc.TaskLabels.ListLabelIDsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return Out{}, fmt.Errorf("list task labels: %w", err)
	}

	users, err := uc.Users.ListDistinctByWorkspaceIDs(ctx, wsIDs)
	if err != nil {
		return Out{}, err
	}

	return Out{
		User:       u,
		Workspaces: workspaces,
		RoleByWS:   roleByWS,
		Projects:   projects,
		Columns:    columns,
		Tasks:      tasks,
		TaskLabels: taskToLabels,
		Labels:     labels,
		Users:      users,
	}, nil
}
