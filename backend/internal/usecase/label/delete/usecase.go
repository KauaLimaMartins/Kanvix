package delete

import (
	"context"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID  string
	LabelID string
}

type labelRepo interface {
	GetByID(ctx context.Context, id string) (entity.Label, error)
	Delete(ctx context.Context, id string) error
}

type memberRepo interface {
	role.MemberRepo
}

type workspaceRepo interface {
	role.WorkspaceRepo
}

type UseCase struct {
	Labels     labelRepo
	Members    memberRepo
	Workspaces workspaceRepo
}

func (uc UseCase) Execute(ctx context.Context, in In) error {
	l, err := uc.Labels.GetByID(ctx, in.LabelID)
	if err != nil {
		return err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, l.WorkspaceID, "admin"); err != nil {
		return err
	}
	return uc.Labels.Delete(ctx, in.LabelID)
}

