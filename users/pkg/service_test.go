package pkg

import (
	"context"
	"crypto/tls"
	"errors"
	"testing"
	"users/pkg/mocks"
	"users/postgres/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("CheckUserEmailExists", "test@example.com").
			Return(false, nil).
			On("CreateUser", "testuser", "test@example.com", mock.Anything).
			Return(sqlc.FMUser{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: mock.Anything,
			}, nil)

		service := NewUserService(mockRepo, &tls.Config{})

		_, err := service.CreateUser(ctx, "testuser", "test@example.com", "password")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - email already exists", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("CheckUserEmailExists", "test@example.com").
			Return(true, nil)

		service := NewUserService(mockRepo, &tls.Config{})

		_, err := service.CreateUser(ctx, "testuser", "test@example.com", "password")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("CheckUserEmailExists", "test@example.com").
			Return(false, nil).
			On("CreateUser", "testuser", "test@example.com", mock.Anything).
			Return(sqlc.FMUser{}, errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{})

		_, err := service.CreateUser(ctx, "testuser", "test@example.com", "password")
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
			On("GetUser", int64(1)).
			Return(sqlc.FMUser{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: mock.Anything,
			}, nil)

		service := NewUserService(mockRepo, &tls.Config{})

		_, err := service.GetUser(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("GetUser", int64(1)).
			Return(sqlc.FMUser{}, errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{})

		_, err := service.GetUser(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestListUsers(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("ListUsers").
			Return([]sqlc.FMUser{
				{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Password: mock.Anything,
				},
			}, nil)

		service := NewUserService(mockRepo, &tls.Config{})

		_, err := service.GetAllUsers(ctx, 1, 10)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("ListUsers").
			Return([]sqlc.FMUser{}, errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{})

		_, err := service.GetAllUsers(ctx, 1, 10)
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
			On("UpdateUser", int64(1), "testuser", "test@example.com", mock.Anything).
			Return(nil)

		service := NewUserService(mockRepo, &tls.Config{})

		err := service.UpdateUser(ctx, 1, "testuser", "test@example.com", "password")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("UpdateUser", int64(1), "testuser", "test@example.com", mock.Anything).
			Return(errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{})

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
			On("DeleteUser", int64(1)).
			Return(nil)

		service := NewUserService(mockRepo, &tls.Config{})

		err := service.DeleteUser(ctx, 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case - repository error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo := mocks.NewUserRepository(t)

		mockRepo.
			On("DeleteUser", int64(1)).
			Return(errors.New("repository error"))

		service := NewUserService(mockRepo, &tls.Config{})

		err := service.DeleteUser(ctx, 1)
		assert.Error(t, err)
		assert.Equal(t, "repository error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
