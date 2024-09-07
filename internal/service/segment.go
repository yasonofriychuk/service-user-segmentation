package service

import (
	"context"
	"errors"
	"github.com/passionde/user-segmentation-service/internal/entity"
	"github.com/passionde/user-segmentation-service/internal/repo"
	"github.com/passionde/user-segmentation-service/internal/repo/repoerrs"
)

type SegmentService struct {
	segmentRepo repo.Segment
	historyRepo repo.History
	userRepo    repo.User
}

func NewSegmentService(segmentRepo repo.Segment, historyRepo repo.History, userRepo repo.User) *SegmentService {
	return &SegmentService{
		segmentRepo: segmentRepo,
		historyRepo: historyRepo,
		userRepo:    userRepo,
	}
}

func (s *SegmentService) CreateSegment(ctx context.Context, input CreateSegmentInput) error {
	err := s.segmentRepo.CreateSegment(ctx, input.Slug)
	if err != nil {
		if errors.Is(err, repoerrs.ErrAlreadyExists) {
			return ErrSegmentAlreadyExists
		}
		return ErrCannotCreateSegment
	}
	if input.PercentageUsers <= 0 {
		return nil
	}

	// todo вынести в фоновый процесс с использование RabbitMQ
	usersID, err := s.userRepo.GetRandomUsers(ctx, input.PercentageUsers)
	if err != nil {
		return err
	}
	for _, userID := range usersID {
		err := s.userRepo.SetSegments(ctx, userID, []string{input.Slug}, []string{})
		if err != nil {
			return err
		}
	}
	return s.historyRepo.AddNotes(ctx, cookNotesSegmentAdd(usersID, input.Slug))
}

func (s *SegmentService) DeleteSegment(ctx context.Context, input SegmentInput) error {
	usersID, err := s.segmentRepo.GetUsersInSegment(ctx, input.Slug)
	if err != nil {
		return err
	}

	err = s.segmentRepo.DeleteSegment(ctx, input.Slug)
	if err != nil {
		if errors.Is(err, repoerrs.ErrNotFound) {
			return ErrSegmentNotFound
		}
		return err
	}
	return s.historyRepo.AddNotes(ctx, cookNotesSegmentDel(usersID, input.Slug))
}

func cookNotesSegmentDel(usersID []string, segment string) []entity.History {
	notes := make([]entity.History, 0, len(usersID))
	for _, userID := range usersID {
		notes = append(notes, entity.History{
			UserID:      userID,
			SegmentSlug: segment,
			Type:        entity.OperationTypeSegmentDelete,
		})
	}
	return notes
}

func cookNotesSegmentAdd(usersID []string, segment string) []entity.History {
	notes := make([]entity.History, 0, len(usersID))
	for _, userID := range usersID {
		notes = append(notes, entity.History{
			UserID:      userID,
			SegmentSlug: segment,
			Type:        entity.OperationTypeAutoAdd,
		})
	}
	return notes
}
