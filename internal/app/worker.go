package app

import (
	"context"
	"github.com/passionde/user-segmentation-service/internal/entity"
	"github.com/passionde/user-segmentation-service/internal/service"
	log "github.com/sirupsen/logrus"
	"time"
)

func RunWorker(services *service.Services) {
	// TODO: Реализовать graceful shutdown
	for {
		time.Sleep(45 * time.Second)

		ctx := context.Background()
		tasks, err := services.TaskDelete.GetExpiredTasks(ctx)
		if err != nil {
			log.Errorf("App - RunWorker - services.TaskDelete.GetExpiredTasks: %v", err)
		}

		err = services.TaskDelete.CompleteTasks(ctx, tasks, handler(services.User))
		if err != nil {
			log.Errorf("App - RunWorker - services.TaskDelete.CompleteTasks: %v", err)
		}
	}
}

func handler(userService service.User) func([]entity.Task) error {
	return func(tasks []entity.Task) error {
		ctx := context.Background()
		for _, setSegmentsInput := range getSegmentsInput(tasks) {
			log.Debug(setSegmentsInput) // todo
			if err := userService.SetSegments(ctx, setSegmentsInput); err != nil {
				return err
			}
		}
		return nil
	}
}

func getSegmentsInput(tasks []entity.Task) []service.SetSegmentsUserInput {
	usersMap := make(map[string][]string)
	for _, task := range tasks {
		usersMap[task.UserID] = append(usersMap[task.UserID], task.SegmentSlug)
	}

	result := make([]service.SetSegmentsUserInput, 0, len(usersMap))
	for userID, segmentsDel := range usersMap {
		result = append(result, service.SetSegmentsUserInput{
			UserID:      userID,
			SegmentsAdd: make([]string, 0, 1),
			SegmentsDel: segmentsDel,
			TTL:         0,
		})
	}
	return result
}
