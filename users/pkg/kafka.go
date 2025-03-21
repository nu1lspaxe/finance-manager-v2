package pkg

import (
	"context"
	"fmt"
	"sync"
	"users/utils"

	"github.com/segmentio/kafka-go"
)

func (s *UserService) publishMessage(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	err := s.KafkaWriter.WriteMessages(ctx, msg)

	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateAccountBalanceJob(ctx context.Context) error {
	userIds, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, len(userIds))

	for _, userId := range userIds {
		wg.Add(1)
		go func(userId int64, ctx context.Context) {
			defer wg.Done()

			accounts, err := s.repo.GetUserAccounts(ctx, userId)
			if err != nil {
				errChan <- err
				return
			}
			for _, account := range accounts {
				err := s.UpdateAccountBalance(ctx, account.ID, account.IDNumber)
				if err != nil {
					errChan <- err
				}

				err = s.publishMessage(ctx, []byte(fmt.Append([]byte{}, userId)), []byte(account.IDNumber))
				if err != nil {
					errChan <- err
				}
			}
		}(userId, ctx)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errors []string
	for err := range errChan {
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return utils.NewUserError(utils.ErrUpdateAccountBalance, errors...)
	}

	return nil
}
