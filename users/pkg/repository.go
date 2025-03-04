package pkg

import (
	"context"
	"users/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(username, email, password string) (sqlc.FianaceManagerFMUser, error)
	CheckUserExists(id int64) (bool, error)
	CheckUserEmailExists(email string) (bool, error)
	GetUserById(id int64) (sqlc.FianaceManagerFMUser, error)
	GetUserByEmail(email string) (sqlc.FianaceManagerFMUser, error)
	GetAllUsers(page, pagesize int32) ([]sqlc.FianaceManagerFMUser, error)
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

func (u *userRepositoryImpl) CreateUser(username, email, password string) (sqlc.FianaceManagerFMUser, error) {
	user, err := u.queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return sqlc.FianaceManagerFMUser{}, err
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

func (u *userRepositoryImpl) GetUserById(id int64) (sqlc.FianaceManagerFMUser, error) {
	user, err := u.queries.GetUserById(context.Background(), id)
	if err != nil {
		return sqlc.FianaceManagerFMUser{}, err
	}
	return user, nil
}

func (u *userRepositoryImpl) GetUserByEmail(email string) (sqlc.FianaceManagerFMUser, error) {
	user, err := u.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return sqlc.FianaceManagerFMUser{}, err
	}
	return user, nil
}

func (u *userRepositoryImpl) GetAllUsers(page, pagesize int32) ([]sqlc.FianaceManagerFMUser, error) {
	offset := (page - 1) * pagesize

	users, err := u.queries.GetAllUsers(context.Background(), sqlc.GetAllUsersParams{
		Limit:  pagesize,
		Offset: offset,
	})
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
