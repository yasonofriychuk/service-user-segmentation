package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/passionde/user-segmentation-service/internal/entity"
	"github.com/passionde/user-segmentation-service/pkg/postgres"
	"time"
)

type TasksDeleteRepo struct {
	*postgres.Postgres
}

func NewTasksDeleteRepo(pg *postgres.Postgres) *TasksDeleteRepo {
	return &TasksDeleteRepo{pg}
}

func (t *TasksDeleteRepo) GetExpiredTasks(ctx context.Context) ([]entity.Task, error) {
	sql, args, _ := t.Builder.
		Select("task_id", "user_id", "segment_slug").
		From("tasks_delete").
		Where("deadline < now()").
		Where(squirrel.Eq{"done": false}).
		ToSql()

	rows, err := t.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TasksDeleteRepo.GetExpiredTasks - t.Pool.Query: %v", err)
	}
	defer rows.Close()

	expiredTasks := make([]entity.Task, 0, 1)
	for rows.Next() {
		task := entity.Task{}
		_ = rows.Scan(&task.TaskID, &task.UserID, &task.SegmentSlug)
		expiredTasks = append(expiredTasks, task)
	}
	return expiredTasks, nil
}

func (t *TasksDeleteRepo) ChangeStatusTasks(ctx context.Context, tasks []entity.Task) error {
	tasksID := make([]int, 0, len(tasks))
	for _, task := range tasks {
		tasksID = append(tasksID, task.TaskID)
	}

	sql, args, _ := t.Builder.
		Update("tasks_delete").
		Set("done", true).
		Where(squirrel.Eq{"task_id": tasksID}).
		ToSql()

	err := t.Pool.QueryRow(ctx, sql, args...).Scan()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("TasksDeleteRepo.ChangeStatusTasks - u.Pool.QueryRow: %v", err)
	}
	return nil
}

func (t *TasksDeleteRepo) CreateTasks(ctx context.Context, tasks []entity.Task, ttl uint64) error {
	var serverTime time.Time
	err := t.Pool.QueryRow(ctx, "SELECT now()").Scan(&serverTime)
	if err != nil {
		return fmt.Errorf("TasksDeleteRepo.CreateTasks - u.Pool.QueryRow (serverTime): %v", err)
	}
	deadline := serverTime.Add(time.Duration(ttl) * time.Minute)

	b := t.Builder.Insert("tasks_delete").Columns("user_id", "segment_slug", "deadline")
	for _, task := range tasks {
		b = b.Values(task.UserID, task.SegmentSlug, deadline)
	}
	sql, args, _ := b.ToSql()

	err = t.Pool.QueryRow(ctx, sql, args...).Scan()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("TasksDeleteRepo.CreateTasks - u.Pool.QueryRow: %v", err)
	}
	return nil
}
