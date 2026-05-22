package delete

import (
	"context"

	"kanvix/backend/internal/domain"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID      string
	WorkspaceID string
}

type workspaceRepo interface {
	role.WorkspaceRepo
	Delete(ctx context.Context, id string) error
}

type memberRepo interface {
	role.MemberRepo
}

type UseCase struct {
	Workspaces workspaceRepo
	Members    memberRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) error {
	w, err := uc.Workspaces.GetByID(ctx, in.WorkspaceID)
	if err != nil {
		return err
	}
	if w.OwnerID != in.UserID {
		return domain.ErrForbidden
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin"); err != nil {
		return err
	}
	return uc.Workspaces.Delete(ctx, in.WorkspaceID)
}
