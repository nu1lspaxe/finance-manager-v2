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

	err := s.kafkaWriter.WriteMessages(ctx, msg)
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

	var mux sync.Mutex
	var wg sync.WaitGroup
	var errors []string

	for _, userId := range userIds {
		wg.Add(1)
		go func(userId int64) {
			accounts, err := s.repo.GetUserAccounts(ctx, userId)
			if err != nil {
				mux.Lock()
				errors = append(errors, err.Error())
				mux.Unlock()
				return
			}
			for _, account := range accounts {
				err := s.UpdateAccountBalance(ctx, userId, account.IDNumber)
				if err != nil {
					mux.Lock()
					errors = append(errors, err.Error())
					mux.Unlock()
				}
				err = s.publishMessage(
					ctx,
					[]byte(fmt.Append([]byte{}, userId)),
					[]byte(account.IDNumber),
				)
				if err != nil {
					mux.Lock()
					errors = append(errors, err.Error())
					mux.Unlock()
				}
			}
		}(userId)
	}

	wg.Wait()

	if len(errors) > 0 {
		return utils.NewUserError(utils.ErrUpdateAccountBalance, errors...)
	}

	return nil
}
