package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/passionde/user-segmentation-service/internal/repo/repoerrs"
	"github.com/passionde/user-segmentation-service/pkg/postgres"
)

type SegmentRepo struct {
	*postgres.Postgres
}

func NewSegmentRepo(pg *postgres.Postgres) *SegmentRepo {
	return &SegmentRepo{pg}
}

func (s *SegmentRepo) CreateSegment(ctx context.Context, slug string) error {
	sql, args, _ := s.Builder.
		Insert("segments").
		Columns("slug").
		Values(slug).
		ToSql()

	err := s.Pool.QueryRow(ctx, sql, args...).Scan()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return repoerrs.ErrAlreadyExists
			}
		}

		return fmt.Errorf("SegmentRepo.CreateSegment - s.Pool.QueryRow: %v", err)
	}
	return nil
}

func (s *SegmentRepo) DeleteSegment(ctx context.Context, slug string) error {
	sql, args, _ := s.Builder.
		Delete("segments").
		Where("slug = ?", slug).
		Suffix("RETURNING slug").
		ToSql()

	err := s.Pool.QueryRow(ctx, sql, args...).Scan(&slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerrs.ErrNotFound
		}
		return fmt.Errorf("SegmentRepo.DeleteSegment - s.Pool.QueryRow: %v", err)
	}
	return nil
}

func (s *SegmentRepo) GetUsersInSegment(ctx context.Context, slug string) ([]string, error) {
	sql, args, _ := s.Builder.
		Select("user_id").
		From("user_segments").
		Where("segment_slug = ?", slug).
		ToSql()

	rows, err := s.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("SegmentRepo.GetUsersInSegment - u.Pool.Query: %v", err)
	}
	defer rows.Close()

	usersID := make([]string, 0, 2)
	for rows.Next() {
		var id string
		_ = rows.Scan(&id)
		usersID = append(usersID, id)
	}
	return usersID, nil
}
