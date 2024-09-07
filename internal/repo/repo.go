package repo

import (
	"context"
	"github.com/passionde/user-segmentation-service/internal/entity"
	"github.com/passionde/user-segmentation-service/internal/repo/pgdb"
	"github.com/passionde/user-segmentation-service/pkg/postgres"
)

type User interface {
	SetSegments(ctx context.Context, userID string, segmentsAdd, segmentsDel []string) error
	GetSegments(ctx context.Context, userID string) ([]string, error)
	GetRandomUsers(ctx context.Context, percent int) ([]string, error)
}

type Segment interface {
	CreateSegment(ctx context.Context, slug string) error
	DeleteSegment(ctx context.Context, slug string) error
	GetUsersInSegment(ctx context.Context, slug string) ([]string, error)
}

type History interface {
	AddNotes(ctx context.Context, notes []entity.History) error
	GetNotes(ctx context.Context, userID string, month, year int) ([]entity.History, error)
}

type TaskDelete interface {
	GetExpiredTasks(ctx context.Context) ([]entity.Task, error)
	ChangeStatusTasks(ctx context.Context, tasks []entity.Task) error
	CreateTasks(ctx context.Context, tasks []entity.Task, ttl uint64) error
}

type Auth interface {
	WriteToken(ctx context.Context, token string) (int, error)
	TokenExist(ctx context.Context, token string) (int, error)
}

type Repositories struct {
	User
	Segment
	History
	TaskDelete
	Auth
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:       pgdb.NewUserRepo(pg),
		Segment:    pgdb.NewSegmentRepo(pg),
		History:    pgdb.NewHistoryRepo(pg),
		TaskDelete: pgdb.NewTasksDeleteRepo(pg),
		Auth:       pgdb.NewAuthRepo(pg),
	}
}
