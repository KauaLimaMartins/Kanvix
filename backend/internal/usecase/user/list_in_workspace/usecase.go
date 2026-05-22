package list_in_workspace

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
	Users []entity.WorkspaceUser
}

type userRepo interface {
	ListWorkspaceUsers(ctx context.Context, workspaceID string) ([]entity.WorkspaceUser, error)
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Users      userRepo
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) (Out, error) {
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}
	users, err := uc.Users.ListWorkspaceUsers(ctx, in.WorkspaceID)
	if err != nil {
		return Out{}, err
	}
	return Out{Users: users}, nil
}

