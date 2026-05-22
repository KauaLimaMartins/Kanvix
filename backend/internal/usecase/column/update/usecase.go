package update

import (
	"context"
	"fmt"
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID   string
	ColumnID string
	Patch    map[string]any
}

type Out struct {
	Column entity.Column
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type columnRepo interface {
	GetByID(ctx context.Context, id string) (entity.Column, error)
	Update(ctx context.Context, id string, patch map[string]any) (entity.Column, error)
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
	if len(in.Patch) == 0 {
		return Out{}, fmt.Errorf("empty patch")
	}
	col, err := uc.Columns.GetByID(ctx, in.ColumnID)
	if err != nil {
		return Out{}, err
	}
	p, err := uc.Projects.GetByID(ctx, col.ProjectID)
	if err != nil {
		return Out{}, err
	}
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}
	in.Patch["updated_at"] = time.Now().UTC()
	updated, err := uc.Columns.Update(ctx, in.ColumnID, in.Patch)
	if err != nil {
		return Out{}, err
	}
	return Out{Column: updated}, nil
}

