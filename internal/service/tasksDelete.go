package service

import (
	"context"
	"github.com/passionde/user-segmentation-service/internal/entity"
	"github.com/passionde/user-segmentation-service/internal/repo"
)

type TasksDeleteService struct {
	tasksDeleteRepo repo.TaskDelete
}

func NewTasksDeleteService(tasksDeleteRepo repo.TaskDelete) *TasksDeleteService {
	return &TasksDeleteService{
		tasksDeleteRepo: tasksDeleteRepo,
	}
}

func (t *TasksDeleteService) GetExpiredTasks(ctx context.Context) ([]entity.Task, error) {
	return t.tasksDeleteRepo.GetExpiredTasks(ctx)
}

func (t *TasksDeleteService) CompleteTasks(ctx context.Context, tasks []entity.Task, callback func([]entity.Task) error) error {
	if err := callback(tasks); err != nil {
		return err
	}
	return t.tasksDeleteRepo.ChangeStatusTasks(ctx, tasks)
}

func (t *TasksDeleteService) CreateTasks(ctx context.Context, tasks []entity.Task, ttl uint64) error {
	return t.tasksDeleteRepo.CreateTasks(ctx, tasks, ttl)
}
