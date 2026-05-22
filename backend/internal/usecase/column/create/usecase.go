package create

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID    string
	ProjectID string
	Name      string
}

type Out struct {
	Column entity.Column
}

type projectRepo interface {
	GetByID(ctx context.Context, id string) (entity.Project, error)
}

type columnRepo interface {
	CountByProject(ctx context.Context, projectID string) (int64, error)
	Create(ctx context.Context, c entity.Column) (entity.Column, error)
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
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, p.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}

	n, err := uc.Columns.CountByProject(ctx, in.ProjectID)
	if err != nil {
		return Out{}, err
	}

	now := time.Now().UTC()
	c, err := uc.Columns.Create(ctx, entity.Column{
		ID:        uuid.NewString(),
		ProjectID: in.ProjectID,
		Name:      in.Name,
		Order:     int(n),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return Out{}, err
	}
	return Out{Column: c}, nil
}

