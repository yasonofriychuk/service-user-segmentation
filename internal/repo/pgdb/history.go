package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/passionde/user-segmentation-service/internal/entity"
	"github.com/passionde/user-segmentation-service/pkg/postgres"
)

type HistoryRepo struct {
	*postgres.Postgres
}

func NewHistoryRepo(pg *postgres.Postgres) *HistoryRepo {
	return &HistoryRepo{pg}
}

func (h *HistoryRepo) AddNotes(ctx context.Context, notes []entity.History) error {
	b := h.Builder.Insert("history").Columns("user_id", "segment_slug", "type")
	for _, note := range notes {
		b = b.Values(note.UserID, note.SegmentSlug, note.Type)
	}
	sql, args, _ := b.ToSql()

	err := h.Pool.QueryRow(ctx, sql, args...).Scan()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("HistoryRepo.AddNote - s.Pool.QueryRow: %v", err)
	}
	return nil
}

func (h *HistoryRepo) GetNotes(ctx context.Context, userID string, month, year int) ([]entity.History, error) {
	sql, args, _ := h.Builder.
		Select("user_id", "segment_slug", "type", "created_at").
		From("history").
		Where("user_id = ? and extract(month from created_at) = ? and extract(year from created_at) = ?", userID, month, year).
		ToSql()

	rows, err := h.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("HistoryRepo.GetNotes - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	notes := make([]entity.History, 0, 1)
	for rows.Next() {
		note := entity.History{}
		err = rows.Scan(&note.UserID, &note.SegmentSlug, &note.Type, &note.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("HistoryRepo.GetNotes - rows.Scan: %v", err)
		}
		notes = append(notes, note)
	}
	return notes, nil
}
