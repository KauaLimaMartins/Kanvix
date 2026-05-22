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
	Columns []entity.Column
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type columnRepo interface {
	ListByProject(ctx context.Context, projectID string) ([]entity.Column, error)
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Projects   projectRepo
	Columns    columnRepo
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
	cols, err := uc.Columns.ListByProject(ctx, in.ProjectID)
	if err != nil {
		return Out{}, err
	}
	return Out{Columns: cols}, nil
}

