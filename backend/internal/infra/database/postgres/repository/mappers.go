package repository

import (
	"time"

	"kanvix/backend/internal/domain/entity"
	"kanvix/backend/internal/infra/database/postgres/model"
)

func userToEntity(u model.User) entity.User {
	return entity.User{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		AvatarColor:  u.AvatarColor,
		Role:         u.Role,
		Disabled:     u.Disabled,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func userFromEntity(u entity.User) model.User {
	return model.User{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		AvatarColor:  u.AvatarColor,
		Role:         u.Role,
		Disabled:     u.Disabled,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func workspaceFromEntity(w entity.Workspace) model.Workspace {
	return model.Workspace{
		ID:        w.ID,
		OwnerID:   w.OwnerID,
		Name:      w.Name,
		Icon:      w.Icon,
		Color:     w.Color,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

func workspaceToEntity(w model.Workspace) entity.Workspace {
	return entity.Workspace{
		ID:        w.ID,
		OwnerID:   w.OwnerID,
		Name:      w.Name,
		Icon:      w.Icon,
		Color:     w.Color,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}

func workspaceMemberFromEntity(m entity.WorkspaceMember) model.WorkspaceMember {
	return model.WorkspaceMember{
		WorkspaceID: m.WorkspaceID,
		UserID:      m.UserID,
		Role:        m.Role,
		CreatedAt:   m.CreatedAt,
	}
}

func workspaceMemberToEntity(m model.WorkspaceMember) entity.WorkspaceMember {
	return entity.WorkspaceMember{
		WorkspaceID: m.WorkspaceID,
		UserID:      m.UserID,
		Role:        m.Role,
		CreatedAt:   m.CreatedAt,
	}
}

func projectToEntity(p model.Project) entity.Project {
	return entity.Project{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func projectFromEntity(p entity.Project) model.Project {
	return model.Project{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func columnToEntity(c model.Column) entity.Column {
	return entity.Column{
		ID:        c.ID,
		ProjectID: c.ProjectID,
		Name:      c.Name,
		Order:     c.Order,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func columnFromEntity(c entity.Column) model.Column {
	return model.Column{
		ID:        c.ID,
		ProjectID: c.ProjectID,
		Name:      c.Name,
		Order:     c.Order,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func taskToEntity(t model.Task) entity.Task {
	return entity.Task{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		ColumnID:    t.ColumnID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		AssigneeID:  t.AssigneeID,
		Order:       t.Order,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func taskFromEntity(t entity.Task) model.Task {
	return model.Task{
		ID:          t.ID,
		ProjectID:   t.ProjectID,
		ColumnID:    t.ColumnID,
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		AssigneeID:  t.AssigneeID,
		Order:       t.Order,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func subtaskToEntity(s model.Subtask) entity.Subtask {
	return entity.Subtask{
		ID:        s.ID,
		TaskID:    s.TaskID,
		Title:     s.Title,
		Done:      s.Done,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func subtaskFromEntity(s entity.Subtask) model.Subtask {
	return model.Subtask{
		ID:        s.ID,
		TaskID:    s.TaskID,
		Title:     s.Title,
		Done:      s.Done,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func labelToEntity(l model.Label) entity.Label {
	return entity.Label{
		ID:          l.ID,
		WorkspaceID: l.WorkspaceID,
		Name:        l.Name,
		Color:       l.Color,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

func labelFromEntity(l entity.Label) model.Label {
	return model.Label{
		ID:          l.ID,
		WorkspaceID: l.WorkspaceID,
		Name:        l.Name,
		Color:       l.Color,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}

func taskLabelFrom(taskID, labelID string, createdAt time.Time) model.TaskLabel {
	return model.TaskLabel{
		TaskID:    taskID,
		LabelID:   labelID,
		CreatedAt: createdAt,
	}
}
