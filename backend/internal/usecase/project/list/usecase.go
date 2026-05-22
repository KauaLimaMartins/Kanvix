package list

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID      string
	WorkspaceID string
}

type Out struct {
	Projects []entity.Project
}

type projectRepo interface {
	ListByWorkspace(ctx context.Context, workspaceID string) ([]entity.Project, error)
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

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin", "member"); err != nil {
		return Out{}, err
	}
	ps, err := uc.Projects.ListByWorkspace(ctx, in.WorkspaceID)
	if err != nil {
		return Out{}, err
	}
	return Out{Projects: ps}, nil
}

