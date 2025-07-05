package service

import (
	"context"
	"jsonjunk/internal/helper"
	"jsonjunk/internal/scheduler"
	logger "jsonjunk/pkg/logging"
	"time"

	"go.uber.org/zap"
)

func (s *pasteService) internal(ctx context.Context) {
	s.schedule(ctx)
}

func (s *pasteService) schedule(ctx context.Context) {
	s.registerSchedule(ctx, "paste_soft_remover", helper.GetNextTime(6, 0), 6*time.Hour, s.removeSoftPaste)
	s.registerSchedule(ctx, "paste_hard_remover", helper.GetNextTime(12, 0), 12*time.Hour, s.removeHardPaste)
}

func (s *pasteService) registerSchedule(ctx context.Context, name string, executeAt time.Time, interval time.Duration, callback func(ctx context.Context)) {
	scheduler.Register(&scheduler.Task{
		Value:     name,
		ExecuteAt: executeAt,
		Interval:  interval,
		Ctx:       ctx,
		Action:    callback,
	})
}

func (s *pasteService) removeSoftPaste(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()
	select {
	case <-ctx.Done():
		return
	default:
		matchCount, modifiedCount, err := s.repo.DeleteSoftPaste(ctx)
		if err != nil {
			logger.Log.Error("remove soft paste", zap.Error(err))
			return
		}

		logger.Log.Info("remove soft paste",
			zap.Int("match count", matchCount),
			zap.Int("modified count", modifiedCount),
		)
		return
	}
}

func (s *pasteService) removeHardPaste(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()
	select {
	case <-ctx.Done():
		return
	default:
		removeCount, err := s.repo.DeletHardPaste(ctx)
		if err != nil {
			logger.Log.Error("remove hard paste", zap.Error(err))
			return
		}
		logger.Log.Info("remove hard paste", zap.Int("remove count", removeCount))
		return
	}
}
