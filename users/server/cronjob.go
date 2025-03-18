package server

import (
	"context"
	"time"
)

func (s *Server) StartCronJob(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.service.UpdateAccountBalanceJob(ctx)
			if err != nil {
				s.logger.Error(err.Error())
			}
		}
	}
}
