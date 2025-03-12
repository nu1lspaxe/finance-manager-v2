package pkg

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"users/postgres/sqlc"
	"users/utils"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type UserService struct {
	repo        UserRepository
	httpClient  *http.Client
	kafkaWriter *kafka.Writer
}

func NewUserService(repo UserRepository, tlsConfig *tls.Config, kafkaWriter *kafka.Writer) *UserService {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   tlsConfig,
			ForceAttemptHTTP2: true,
		},
		Timeout: utils.TIMEOUT,
	}

	return &UserService{
		repo:        repo,
		httpClient:  client,
		kafkaWriter: kafkaWriter,
	}
}

func (u *UserService) CreateUser(ctx context.Context, username, email, password string) (*sqlc.FMUser, error) {
	if username == "" || email == "" {
		return nil, utils.NewUserError(utils.ErrUserInvalid, "username or email is empty")
	}
	exists, err := u.repo.CheckUserEmailExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.NewUserError(utils.ErrUserEmailExists)
	}

	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := u.repo.CreateUser(ctx, username, email, hashPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) GetUser(ctx context.Context, userId int64) (*sqlc.FMUser, error) {
	user, err := u.repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) GetAllUsers(ctx context.Context) ([]int64, error) {
	userIds, err := u.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return userIds, nil
}

func (u *UserService) UpdateUser(ctx context.Context, userId int64, username, email, password string) error {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return u.repo.UpdateUser(ctx, userId, username, email, hashPassword)
}

func (u *UserService) DeleteUser(ctx context.Context, userId int64) error {
	return u.repo.DeleteUser(ctx, userId)
}

func (u *UserService) AddAccount(ctx context.Context, userId int64, idNumber string) error {
	exists, err := u.repo.CheckAccountExists(ctx, userId, idNumber)
	if err != nil {
		return err
	}
	if exists {
		return utils.NewUserError(utils.ErrAccountExists)
	}

	balance, err := getAccountBalance(ctx, u.httpClient, idNumber)
	if err != nil {
		return err
	}

	return u.repo.AddAccount(ctx, userId, idNumber, balance)
}

func getAccountBalance(ctx context.Context, client *http.Client, idNumber string) (float64, error) {
	baseURL := viper.GetString("fintech.url")
	url := fmt.Sprintf("%s/accounts/%s/balance", baseURL, idNumber)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "users-service/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, utils.NewUserError(utils.ErrStatusCode, fmt.Sprint(resp.StatusCode))
	}

	type BalanceResponse struct {
		Balance float64 `json:"balance"`
	}

	var balanceResp BalanceResponse
	err = json.NewDecoder(resp.Body).Decode(&balanceResp)
	if err != nil {
		return 0, err
	}
	return balanceResp.Balance, nil
}

func (u *UserService) GetUserAccounts(ctx context.Context, userId int64) (*[]sqlc.FMAccount, error) {
	accounts, err := u.repo.GetUserAccounts(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &accounts, nil
}

func (u *UserService) DeleteAccount(ctx context.Context, idNumber string) error {
	return u.repo.DeleteAccount(ctx, idNumber)
}
