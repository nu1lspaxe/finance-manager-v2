package pkg

import (
	"context"
	"users/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(username, email, password string) (sqlc.User, error)
	CheckUserExists(id int64) (bool, error)
	CheckUserEmailExists(email string) (bool, error)
	GetUser(id int64) (sqlc.User, error)
	ListUsers() ([]sqlc.User, error)
	UpdateUser(id int64, username, email string, password string) error
	DeleteUser(id int64) error
}

type userRepositoryImpl struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepositoryImpl{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (u *userRepositoryImpl) CreateUser(username, email, password string) (sqlc.User, error) {
	user, err := u.queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (u *userRepositoryImpl) CheckUserExists(id int64) (bool, error) {
	exists, err := u.queries.CheckUserExists(context.Background(), id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *userRepositoryImpl) CheckUserEmailExists(email string) (bool, error) {
	exists, err := u.queries.CheckUserEmailExists(context.Background(), email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *userRepositoryImpl) GetUser(id int64) (sqlc.User, error) {
	user, err := u.queries.GetUserById(context.Background(), id)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (u *userRepositoryImpl) ListUsers() ([]sqlc.User, error) {
	users, err := u.queries.ListUsers(context.Background())
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepositoryImpl) UpdateUser(id int64, username, email string, password string) error {
	err := u.queries.UpdateUser(context.Background(), sqlc.UpdateUserParams{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepositoryImpl) DeleteUser(id int64) error {
	err := u.queries.DeleteUser(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}
