package pkg

import (
	"context"
	"fmt"
	"sync"
	"users/utils"

	"github.com/segmentio/kafka-go"
)

func (u *UserService) publishBankAccount(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	err := u.kafkaWriter.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateAccountBalance(ctx context.Context) error {
	userIds, err := u.repo.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	var mux sync.Mutex
	var wg sync.WaitGroup
	var errors []string

	for _, userId := range userIds {
		wg.Add(1)
		go func(userId int64) {
			accounts, err := u.repo.GetUserAccounts(ctx, userId)
			if err != nil {
				mux.Lock()
				errors = append(errors, err.Error())
				mux.Unlock()
				return
			}
			for _, account := range accounts {
				err := u.publishBankAccount(ctx, []byte(fmt.Append(nil, userId)), []byte(account.IDNumber))
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
