package delete

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID    string
	ProjectID string
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
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
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) error {
	p, err := uc.Projects.GetByID(ctx, in.ProjectID)
	if err != nil {
		return err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin"); err != nil {
		return err
	}
	return uc.Projects.Delete(ctx, in.ProjectID)
}

