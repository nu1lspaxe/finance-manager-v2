package pkg

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

func (s *BankRecordService) ConsumeUpdateAccountMessage(ctx context.Context) error {
	defer func() {
		if err := s.kafkaReader.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing Kafka reader: %v\n", err)
		}
	}()

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
