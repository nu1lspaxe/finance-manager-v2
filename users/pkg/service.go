package pkg

import (
	"users/postgres/sqlc"
	"users/utils"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) CreateUser(username, email string) (*sqlc.User, error) {
	if username == "" || email == "" {
		return nil, utils.NewUserError(utils.USER_INVALID, "Username and email are required")
	}
	exists, err := u.repo.CheckUserEmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, utils.NewUserError(utils.USER_EXISTS, "User already exists")
	}

	user, err := u.repo.CreateUser(username, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) GetUser(userId int64) (*sqlc.User, error) {
	user, err := u.repo.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserService) ListUsers(page, pagesize uint32) (*[]sqlc.User, error) {
	users, err := u.repo.ListUsers()
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (u *UserService) UpdateUser(userId int64, username, email, password string) error {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return u.repo.UpdateUser(userId, username, email, hashPassword)
}

func (u *UserService) DeleteUser(userId int64) error {
	return u.repo.DeleteUser(userId)
}
