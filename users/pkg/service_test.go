package pkg

import (
	"context"
	"crypto/tls"
	"errors"
	"testing"
	"users/pkg/mocks"
	"users/postgres/sqlc"
	"users/utils"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("CheckUserEmailExists", ctx, "test@example.com").
			Return(false, nil).
			On("CreateUser", ctx, "testuser", "test@example.com", mock.Anything).
			Return(sqlc.FMUser{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: mock.Anything,
			}, nil)

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		_, err := service.SignUp(ctx, "testuser", "test@example.com", "password")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - email already exists", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("CheckUserEmailExists", ctx, "test@example.com").
			Return(true, nil)

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		_, err := service.SignUp(ctx, "testuser", "test@example.com", "password")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("CheckUserEmailExists", ctx, "test@example.com").
			Return(false, nil).
			On("CreateUser", ctx, "testuser", "test@example.com", mock.Anything).
			Return(sqlc.FMUser{}, errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		_, err := service.SignUp(ctx, "testuser", "test@example.com", "password")
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("GetUserById", ctx, int64(1)).
			Return(sqlc.GetUserByIdRow{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: mock.Anything,
			}, nil)

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		_, err := service.GetUser(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("GetUserById", ctx, int64(1)).
			Return(sqlc.GetUserByIdRow{}, errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		_, err := service.GetUser(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllUsers(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("GetAllUsers", ctx).
			Return([]int64{1}, nil)

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		_, err := service.GetAllUsers(ctx)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("GetAllUsers", ctx).
			Return([]int64{}, errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		_, err := service.GetAllUsers(ctx)
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("UpdateUser", ctx, int64(1), "testuser", "test@example.com", mock.Anything).
			Return(nil)

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		err := service.UpdateUser(ctx, 1, "testuser", "test@example.com", "password")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("UpdateUser", ctx, int64(1), "testuser", "test@example.com", mock.Anything).
			Return(errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		err := service.UpdateUser(ctx, 1, "testuser", "test@example.com", "password")
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("DeleteUser", ctx, int64(1)).
			Return(nil)

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})
		err := service.DeleteUser(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("DeleteUser", ctx, int64(1)).
			Return(errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{}, &kafka.Writer{}, &utils.JWTManager{}, &redis.Client{})

		err := service.DeleteUser(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
