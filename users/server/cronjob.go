package server

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func (s *Server) StartCronJob(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if err := s.service.KafkaWriter.Close(); err != nil {
				s.logger.Error("Failed to close Kafka writer", zap.Error(err))
			}
			if err := s.logger.Sync(); err != nil {
				s.logger.Error("Failed to sync logger", zap.Error(err))
			}
			return
		case <-ticker.C:
			err := s.service.UpdateAccountBalanceJob(ctx)
			if err != nil {
				s.logger.Error(err.Error())
			} else {
				s.logger.Info("Cron job completed", zap.String("time", fmt.Sprintf("%v", time.Now())))
			}
		}
	}
}
