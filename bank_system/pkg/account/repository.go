package account

import (
	"bank_system/pkg/transaction"
	"bank_system/postgres/sqlc"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, userID int64) (sqlc.BKAccount, error)
	GetAccountByIDNumber(ctx context.Context, idNumber string) (sqlc.GetAccountByIDNumberRow, error)
	GetAccountTransactions(ctx context.Context, accountID int64) ([]sqlc.GetAccountTransactionsRow, error)
	GetAllAccounts(ctx context.Context) ([]sqlc.GetAllAccountsRow, error)
	WithdrawFromAccount(ctx context.Context, accountID int64, amount float64, detail string, txRepo transaction.TxRepository) (float64, error)
	DepositToAccount(ctx context.Context, accountID int64, amount float64, detail string, txRepo transaction.TxRepository) (float64, error)
}

type accountRepositoryImpl struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewAccountRepository(pool *pgxpool.Pool) AccountRepository {
	return &accountRepositoryImpl{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *accountRepositoryImpl) CreateAccount(ctx context.Context, userID int64) (sqlc.BKAccount, error) {
	tx, err := r.queries.CreateAccount(ctx, userID)
	if err != nil {
		return sqlc.BKAccount{}, err
	}
	return tx, nil
}

func (r *accountRepositoryImpl) GetAccountByIDNumber(ctx context.Context, idNumber string) (sqlc.GetAccountByIDNumberRow, error) {
	return r.queries.GetAccountByIDNumber(ctx, idNumber)
}

func (r *accountRepositoryImpl) GetAccountTransactions(ctx context.Context, accountID int64) ([]sqlc.GetAccountTransactionsRow, error) {
	return r.queries.GetAccountTransactions(ctx, accountID)
}

func (r *accountRepositoryImpl) GetAllAccounts(ctx context.Context) ([]sqlc.GetAllAccountsRow, error) {
	accounts, err := r.queries.GetAllAccounts(ctx)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepositoryImpl) WithdrawFromAccount(
	ctx context.Context, accountID int64, amount float64, detail string, txRepo transaction.TxRepository,
) (float64, error) {
	txOptions := pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}

	tx, err := r.pool.BeginTx(ctx, txOptions)

	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	balanceAfter, err := r.queries.WithTx(tx).WithdrawFromAccount(ctx, sqlc.WithdrawFromAccountParams{
		ID:      accountID,
		Balance: amount})

	if err != nil {
		return 0, err
	}

	_, err = txRepo.CreateTransaction(ctx, accountID, amount, transaction.TxType_WITHDRAW, "")
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return balanceAfter, nil
}

func (r *accountRepositoryImpl) DepositToAccount(
	ctx context.Context, accountID int64, amount float64, detail string, txRepo transaction.TxRepository,
) (float64, error) {
	txOptions := pgx.TxOptions{
		IsoLevel:       pgx.Serializable,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}

	tx, err := r.pool.BeginTx(ctx, txOptions)

	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	balanceAfter, err := r.queries.WithTx(tx).DepositToAccount(ctx, sqlc.DepositToAccountParams{
		ID:      accountID,
		Balance: amount})

	if err != nil {
		return 0, err
	}

	_, err = txRepo.CreateTransaction(ctx, accountID, amount, transaction.TxType_DEPOSIT, detail)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return balanceAfter, nil
}
