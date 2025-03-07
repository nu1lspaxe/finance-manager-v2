package pkg

import (
	"context"
	"users/postgres/sqlc"
	"users/utils"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) CreateUser(ctx context.Context, username, email, password string) (*sqlc.FMUser, error) {
	if username == "" || email == "" {
		return nil, utils.NewUserError(utils.ErrUserInvalid, "Username and email are required")
	}
	exists, err := u.repo.CheckUserEmailExists(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.NewUserError(utils.ErrUserEmailExists, "Email already exists")
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

func (u *UserService) GetAllUsers(ctx context.Context, page, pagesize int32) (*[]sqlc.FMUser, error) {
	users, err := u.repo.GetAllUsers(ctx, page, pagesize)
	if err != nil {
		return nil, err
	}
	return &users, nil
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
