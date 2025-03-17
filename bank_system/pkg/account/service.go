package account

import (
	"bank_system/pkg/transaction"
	"bank_system/postgres/sqlc"
	"bank_system/utils"
	"context"
)

type AccountService struct {
	actRepo AccountRepository
	txRepo  transaction.TxRepository
}

func NewAccountService(actRepo AccountRepository, txRepo transaction.TxRepository) *AccountService {
	return &AccountService{
		actRepo: actRepo,
		txRepo:  txRepo,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID int64) (*sqlc.BKAccount, error) {
	account, err := s.actRepo.CreateAccount(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *AccountService) GetAccountByIDNumber(ctx context.Context, idNumber string) (*sqlc.GetAccountByIDNumberRow, error) {
	account, err := s.actRepo.GetAccountByIDNumber(ctx, idNumber)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *AccountService) GetAccountBalance(ctx context.Context, idNumber string) (float64, error) {
	account, err := s.actRepo.GetAccountByIDNumber(ctx, idNumber)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

// func (s *AccountService) GetAccountTransactions

func (s *AccountService) GetAllAccounts(ctx context.Context) (*[]sqlc.GetAllAccountsRow, error) {
	accounts, err := s.actRepo.GetAllAccounts(ctx)
	if err != nil {
		return nil, err
	}
	return &accounts, nil
}

func (s *AccountService) Withdraw(ctx context.Context, idNumber string, amount float64, detail string) (float64, error) {
	account, err := s.actRepo.GetAccountByIDNumber(ctx, idNumber)
	if err != nil {
		return 0, err
	}
	if account.Balance < amount {
		return 0, utils.NewBankSystemError(utils.ErrInsufficientBalance)
	}
	return s.actRepo.WithdrawFromAccount(ctx, account.ID, amount, detail, s.txRepo)
}

func (s *AccountService) Deposit(ctx context.Context, accountID int64, amount float64, detail string) (float64, error) {
	return s.actRepo.DepositToAccount(ctx, accountID, amount, detail, s.txRepo)
}
