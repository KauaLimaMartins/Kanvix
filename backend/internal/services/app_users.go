package services

import (
	"context"

	"kanvix/backend/internal/http/dto"
)

func (s AppService) ListUsers(ctx context.Context) ([]dto.User, error) {
	users, err := s.Repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.User, 0, len(users))
	for _, u := range users {
		out = append(out, dto.User{ID: u.ID, Name: u.Name, AvatarColor: u.AvatarColor})
	}
	return out, nil
}
