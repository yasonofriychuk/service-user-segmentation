package service

import (
	"context"
	"errors"
	"github.com/passionde/user-segmentation-service/internal/entity"
	"github.com/passionde/user-segmentation-service/internal/repo"
	"github.com/passionde/user-segmentation-service/internal/repo/repoerrs"
)

type UserService struct {
	userRepo    repo.User
	taskDelete  repo.TaskDelete
	historyRepo repo.History
}

func NewUserService(userRepo repo.User, historyRepo repo.History, taskDelete repo.TaskDelete) *UserService {
	return &UserService{
		userRepo:    userRepo,
		taskDelete:  taskDelete,
		historyRepo: historyRepo,
	}
}

func (u *UserService) SetSegments(ctx context.Context, input SetSegmentsUserInput) error {
	activeSegments, err := u.GetSegments(ctx, GetSegmentsUserInput{UserID: input.UserID})
	if err != nil {
		if !errors.Is(err, ErrUserNotFound) {
			return err
		}
		activeSegments = make([]string, 0)
	}

	err = u.userRepo.SetSegments(
		ctx,
		input.UserID,
		input.SegmentsAdd,
		input.SegmentsDel,
	)
	if err != nil {
		if errors.Is(err, repoerrs.ErrSegmentsNotExist) {
			return ErrSegmentNotFound
		}
		return err
	}
	if input.TTL > 0 {
		if err := u.taskDelete.CreateTasks(ctx, cookTasks(input, activeSegments), input.TTL); err != nil {
			return err
		}
	}
	return u.historyRepo.AddNotes(ctx, cookNotesUser(input, activeSegments))

}

func (u *UserService) GetSegments(ctx context.Context, input GetSegmentsUserInput) ([]string, error) {
	segments, err := u.userRepo.GetSegments(ctx, input.UserID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return segments, nil
}

func cookNotesUser(input SetSegmentsUserInput, activeSegments []string) []entity.History {
	segmentsAdd := getSegmentsAdd(input.SegmentsAdd, activeSegments)
	segmentsDel := getSegmentsDel(input.SegmentsDel, activeSegments)
	countItems := len(segmentsAdd) + len(segmentsDel)
	notes := make([]entity.History, 0, countItems)

	for _, segment := range segmentsAdd {
		notes = append(notes, entity.History{
			UserID:      input.UserID,
			SegmentSlug: segment,
			Type:        entity.OperationTypeAdd,
		})
	}
	for _, segment := range segmentsDel {
		notes = append(notes, entity.History{
			UserID:      input.UserID,
			SegmentSlug: segment,
			Type:        entity.OperationTypeDelete,
		})
	}
	return notes
}

func cookTasks(input SetSegmentsUserInput, activeSegments []string) []entity.Task {
	segmentsAdd := getSegmentsAdd(input.SegmentsAdd, activeSegments)
	tasks := make([]entity.Task, 0, len(segmentsAdd))
	for _, segment := range segmentsAdd {
		tasks = append(tasks, entity.Task{
			UserID:      input.UserID,
			SegmentSlug: segment,
		})
	}
	return tasks
}

func getSegmentsAdd(segmentsAdd, activeSegments []string) []string {
	filteredSegments := make([]string, 0, 2)
	for _, segment := range segmentsAdd {
		if !contains(activeSegments, segment) && !contains(filteredSegments, segment) {
			filteredSegments = append(filteredSegments, segment)
		}
	}
	return filteredSegments
}

func getSegmentsDel(segmentsDel, activeSegments []string) []string {
	filteredSegments := make([]string, 0, 2)
	for _, segment := range segmentsDel {
		if contains(activeSegments, segment) && !contains(filteredSegments, segment) {
			filteredSegments = append(filteredSegments, segment)
		}
	}
	return filteredSegments
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
