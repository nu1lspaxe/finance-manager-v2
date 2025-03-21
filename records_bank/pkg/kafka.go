package pkg

import (
	"context"
	"strconv"
)

func (s *BankRecordService) ConsumeUpdateAccountMessage(ctx context.Context) error {
	defer s.kafkaReader.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := s.kafkaReader.FetchMessage(ctx)
			if err != nil {
				if err == context.Canceled || err == context.DeadlineExceeded {
					return err
				}
				continue
			}

			userId, err := strconv.ParseInt(string(msg.Key), 10, 64)
			if err != nil {
				continue
			}

			accountNumber := string(msg.Value)
			err = s.CreateBankRecordsBulk(ctx, userId, accountNumber)
			if err != nil {
				continue
			}

			if err := s.kafkaReader.CommitMessages(ctx, msg); err != nil {
				continue
			}
		}
	}
}
