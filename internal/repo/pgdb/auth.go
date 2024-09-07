package pgdb

import (
	"context"
	"fmt"
	"github.com/passionde/user-segmentation-service/pkg/postgres"
)

type AuthRepo struct {
	*postgres.Postgres
}

func NewAuthRepo(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (a *AuthRepo) WriteToken(ctx context.Context, token string) (int, error) {
	sql, args, _ := a.Builder.
		Insert("api_keys").
		Columns("hash_key").
		Values(token).
		Suffix("RETURNING id").
		ToSql()

	var tokenID int
	err := a.Pool.QueryRow(ctx, sql, args...).Scan(&tokenID)
	if err != nil {
		return 0, fmt.Errorf("AuthRepo.WriteToken - u.Pool.QueryRow: %v", err)
	}
	return tokenID, nil
}

func (a *AuthRepo) TokenExist(ctx context.Context, token string) (int, error) {
	sql, args, _ := a.Builder.
		Select("id").
		From("api_keys").
		Where("hash_key = ?", token).
		ToSql()

	var tokenID int
	err := a.Pool.QueryRow(ctx, sql, args...).Scan(&tokenID)
	if err != nil {
		return 0, fmt.Errorf("AuthRepo.TokenExist - u.Pool.QueryRow: %v", err)
	}
	return tokenID, nil
}
