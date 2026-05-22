package create

import (
	"context"
	"time"

	"github.com/google/uuid"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/usecase/workspace/role"
)

type In struct {
	UserID      string
	WorkspaceID string
	Name        string
}

type Out struct {
	Project entity.Project
	Columns []entity.Column
}

type projectRepo interface {
	Create(ctx context.Context, p entity.Project) (entity.Project, error)
}

type columnRepo interface {
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
	if _, err := role.Require(ctx, uc.Members, uc.Workspaces, in.UserID, in.WorkspaceID, "admin"); err != nil {
		return Out{}, err
	}

	now := time.Now().UTC()
	p, err := uc.Projects.Create(ctx, entity.Project{
		ID:          uuid.NewString(),
		WorkspaceID: in.WorkspaceID,
		Name:        in.Name,
		Description: "",
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return Out{}, err
	}

	defaultCols := []struct {
		Name  string
		Order int
	}{
		{Name: "To do", Order: 0},
		{Name: "In progress", Order: 1},
		{Name: "Done", Order: 2},
	}
	cols := make([]entity.Column, 0, len(defaultCols))
	for _, dc := range defaultCols {
		c, err := uc.Columns.Create(ctx, entity.Column{
			ID:        uuid.NewString(),
			ProjectID: p.ID,
			Name:      dc.Name,
			Order:     dc.Order,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return Out{}, err
		}
		cols = append(cols, c)
	}

	return Out{Project: p, Columns: cols}, nil
}

