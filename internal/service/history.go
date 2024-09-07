package service

import (
	"context"
	"fmt"
	"github.com/passionde/user-segmentation-service/internal/entity"
	"github.com/passionde/user-segmentation-service/internal/repo"
	"github.com/passionde/user-segmentation-service/pkg/csvwriter"
	"path"
	"time"
)

type HistoryService struct {
	historyRepo repo.History
	csvWriter   csvwriter.CSVWriter
}

func NewHistoryService(historyRepo repo.History, csvWriter csvwriter.CSVWriter) *HistoryService {
	return &HistoryService{
		historyRepo: historyRepo,
		csvWriter:   csvWriter,
	}
}

func (h *HistoryService) AddNotes(ctx context.Context, notes []entity.History) error {
	return h.historyRepo.AddNotes(ctx, notes)
}

func (h *HistoryService) GetNotes(ctx context.Context, input GetHistoryInput) (string, error) {
	notes, err := h.historyRepo.GetNotes(ctx, input.UserID, input.Month, input.Year)
	if err != nil {
		return "", err
	}
	if len(notes) == 0 {
		return "", ErrUserNoData
	}

	fileName := path.Clean(fmt.Sprintf("%s-%d-%d-%d.csv", input.UserID, input.Year, input.Month, time.Now().Unix()))
	return h.csvWriter.CreateCSVFile(fileName, notes)
}
