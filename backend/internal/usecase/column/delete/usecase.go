package delete

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID   string
	ColumnID string
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type columnRepo interface {
	GetByID(ctx context.Context, id string) (entity.Column, error)
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
	Columns    columnRepo
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) error {
	col, err := uc.Columns.GetByID(ctx, in.ColumnID)
	if err != nil {
		return err
	}
	p, err := uc.Projects.GetByID(ctx, col.ProjectID)
	if err != nil {
		return err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin"); err != nil {
		return err
	}
	return uc.Columns.Delete(ctx, in.ColumnID)
}

