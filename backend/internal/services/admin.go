package services

import (
	"context"

	"kanvix/backend/internal/database"
)

func (s AppService) ResetDemo(ctx context.Context) error {
	return database.ResetDemo(ctx, s.Repo.DB)
}
