package pkg

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"users/postgres/sqlc"
	"users/utils"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type UserService struct {
	repo        UserRepository
	httpClient  *http.Client
	KafkaWriter *kafka.Writer
	jwtManager  *utils.JWTManager
	redisClient *redis.Client
}

func NewUserService(
	repo UserRepository, tlsConfig *tls.Config, kafkaWriter *kafka.Writer, jwtManager *utils.JWTManager, redisClient *redis.Client,
) *UserService {
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
		KafkaWriter: kafkaWriter,
		jwtManager:  jwtManager,
		redisClient: redisClient,
	}
}

func (u *UserService) SignIn(ctx context.Context, email, password string) (int64, string, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return utils.INVALID, "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return utils.INVALID, "", utils.NewUserError(utils.ErrPasswdInvalid)
	}

	issueTime := time.Now()
	tokenString, err := u.jwtManager.Generate(user.ID, issueTime)
	if err != nil {
		return utils.INVALID, "", err
	}

	key := fmt.Sprintf("user:%d", user.ID)
	expiryTime := time.Until(issueTime.Add(utils.TOKEN_EXPIRY))

	err = u.redisClient.HSet(ctx, key, map[string]interface{}{
		"token": tokenString,
	}).Err()
	if err != nil {
		return utils.INVALID, "", err
	}

	err = u.redisClient.Expire(ctx, key, expiryTime).Err()
	if err != nil {
		return utils.INVALID, "", err
	}

	return user.ID, tokenString, nil
}

func (u *UserService) Logout(ctx context.Context, userId int64, token string) error {
	claims, err := u.jwtManager.Verify(token)
	if err != nil {
		return err
	}
	if claims.UserId != userId {
		return utils.NewUserError(utils.ErrTokenInvalid)
	}

	key := fmt.Sprintf("user:%d", userId)
	err = u.redisClient.HDel(ctx, key, token).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) SignUp(ctx context.Context, username, email, password string) (*sqlc.FMUser, error) {
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

func (u *UserService) GetUser(ctx context.Context, userId int64) (*sqlc.GetUserByIdRow, error) {
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

func (u *UserService) AddAccount(ctx context.Context, userId int64, idNumber string) (float64, error) {
	exists, err := u.repo.CheckAccountExists(ctx, userId, idNumber)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, utils.NewUserError(utils.ErrAccountExists)
	}

	balance, err := getAccountBalance(ctx, u.httpClient, idNumber)
	if err != nil {
		return 0, err
	}

	return balance, u.repo.AddAccount(ctx, userId, idNumber, balance)
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing response body: %v\n", err)
		}
	}()

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

func (u *UserService) UpdateAccountBalance(ctx context.Context, accountId int64, idNumber string) error {
	balance, err := getAccountBalance(ctx, u.httpClient, idNumber)
	if err != nil {
		return err
	}
	return u.repo.UpdateAccountBalance(ctx, accountId, balance)
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
