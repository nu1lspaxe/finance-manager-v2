package pkg

import (
	"context"
	"users/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username, email, password string) (sqlc.FMUser, error)
	CheckUserExists(ctx context.Context, id int64) (bool, error)
	CheckUserEmailExists(ctx context.Context, email string) (bool, error)
	GetUserById(ctx context.Context, id int64) (sqlc.GetUserByIdRow, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error)
	GetAllUsers(ctx context.Context) ([]int64, error)
	UpdateUser(ctx context.Context, d int64, username, email string, password string) error
	DeleteUser(ctx context.Context, id int64) error

	AddAccount(ctx context.Context, userId int64, idNumber string, balance float64) error
	CheckAccountExists(ctx context.Context, userId int64, idNumber string) (bool, error)
	GetUserAccounts(ctx context.Context, userId int64) ([]sqlc.FMAccount, error)
	UpdateAccountBalance(ctx context.Context, id int64, balance float64) error
	DeleteAccount(ctx context.Context, idNumber string) error
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

func (u *userRepositoryImpl) CreateUser(ctx context.Context, username, email, password string) (sqlc.FMUser, error) {
	user, err := u.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Username: username,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return sqlc.FMUser{}, err
	}
	return user, nil
}

func (u *userRepositoryImpl) CheckUserExists(ctx context.Context, id int64) (bool, error) {
	exists, err := u.queries.CheckUserExists(ctx, id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *userRepositoryImpl) CheckUserEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := u.queries.CheckUserEmailExists(ctx, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *userRepositoryImpl) GetUserById(ctx context.Context, id int64) (sqlc.GetUserByIdRow, error) {
	user, err := u.queries.GetUserById(ctx, id)
	if err != nil {
		return sqlc.GetUserByIdRow{}, err
	}
	return user, nil
}

func (u *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (sqlc.GetUserByEmailRow, error) {
	user, err := u.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return sqlc.GetUserByEmailRow{}, err
	}
	return user, nil
}

func (u *userRepositoryImpl) GetAllUsers(ctx context.Context) ([]int64, error) {
	users, err := u.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepositoryImpl) UpdateUser(ctx context.Context, id int64, username, email string, password string) error {
	err := u.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
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

func (u *userRepositoryImpl) DeleteUser(ctx context.Context, id int64) error {
	err := u.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepositoryImpl) AddAccount(ctx context.Context, userId int64, idNumber string, balance float64) error {
	_, err := u.queries.AddAccount(ctx, sqlc.AddAccountParams{
		UserID:   userId,
		IDNumber: idNumber,
		Balance:  balance,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepositoryImpl) CheckAccountExists(ctx context.Context, userId int64, idNumber string) (bool, error) {
	exists, err := u.queries.CheckAccountExists(ctx, sqlc.CheckAccountExistsParams{
		UserID:   userId,
		IDNumber: idNumber,
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (u *userRepositoryImpl) GetUserAccounts(ctx context.Context, userId int64) ([]sqlc.FMAccount, error) {
	accounts, err := u.queries.GetUserAccounts(ctx, userId)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (u *userRepositoryImpl) UpdateAccountBalance(ctx context.Context, id int64, balance float64) error {
	err := u.queries.UpdateAccountBalance(ctx, sqlc.UpdateAccountBalanceParams{
		ID:      id,
		Balance: balance,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepositoryImpl) DeleteAccount(ctx context.Context, idNumber string) error {
	err := u.queries.DeleteAccount(ctx, idNumber)
	if err != nil {
		return err
	}
	return nil
}
