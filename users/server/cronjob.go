package server

import (
	"context"
	"time"
)

func (s *Server) StartCronJob(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.service.UpdateAccountBalance(ctx)
			if err != nil {
				s.logger.Error(err.Error())
			}
		}
	}
}
