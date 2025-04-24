package pkg

import (
	"context"
	"records_bank/postgres/sqlc"
	"records_bank/utils"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BankRecordRepository interface {
	CreateBankRecordsBulk(ctx context.Context, userId int64, accountNumber string, amounts []float64, txTypes []string, txIds []int64, txDates []pgtype.Timestamptz, details []string) error
	GetBankRecord(ctx context.Context, id int64) (sqlc.FMRecordBank, error)
	GetUserBankRecords(ctx context.Context, userId int64, accountNumber string) ([]sqlc.FMRecordBank, error)
	GetUserBankRecordsWithPeriod(ctx context.Context, userId int64, accountNumber string, startTime, endTime int64) ([]sqlc.FMRecordBank, error)
	GetUserBankRecordsFromDate(ctx context.Context, userId int64, accountNumber string, startTime int64) ([]sqlc.FMRecordBank, error)
	GetUserBankRecordsToDate(ctx context.Context, userId int64, accountNumber string, endTime int64) ([]sqlc.FMRecordBank, error)
	GetUserBankRecordsByType(ctx context.Context, userId int64, accountNumber, recordType string) ([]sqlc.FMRecordBank, error)
	GetUserBankRecordsByTypeWithPeriod(ctx context.Context, userId int64, accountNumber, recordType string, startTime, endTime int64) ([]sqlc.FMRecordBank, error)
	GetUserBankRecordsByTypeFromDate(ctx context.Context, userId int64, accountNumber, recordType string, startTime int64) ([]sqlc.FMRecordBank, error)
	GetUserBankRecordsByTypeToDate(ctx context.Context, userId int64, accountNumber, recordType string, endTime int64) ([]sqlc.FMRecordBank, error)
	GetExistingBankRecords(ctx context.Context, ids []int64) ([]int64, error)
}

type bankRecordRepositoryImpl struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewBankRecordRepository(pool *pgxpool.Pool) BankRecordRepository {
	return &bankRecordRepositoryImpl{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *bankRecordRepositoryImpl) CreateBankRecordsBulk(
	ctx context.Context, userId int64, accountNumber string, amounts []float64, txTypes []string, txIds []int64, txDates []pgtype.Timestamptz, details []string,
) error {
	_, err := r.queries.CreateBankRecordsBulk(ctx, sqlc.CreateBankRecordsBulkParams{
		Column1: userId,
		Column2: accountNumber,
		Column3: amounts,
		Column4: txTypes,
		Column5: txIds,
		Column6: txDates,
		Column7: details,
	})
	return err
}

func (r *bankRecordRepositoryImpl) GetBankRecord(ctx context.Context, id int64) (sqlc.FMRecordBank, error) {
	return r.queries.GetBankRecord(ctx, id)
}

func (r *bankRecordRepositoryImpl) GetUserBankRecords(ctx context.Context, userId int64, accountNumber string) ([]sqlc.FMRecordBank, error) {
	return r.queries.GetUserBankRecords(ctx, sqlc.GetUserBankRecordsParams{
		UserID:        userId,
		AccountNumber: accountNumber,
	})
}

func (r *bankRecordRepositoryImpl) GetUserBankRecordsWithPeriod(ctx context.Context, userId int64, accountNumber string, startTime, endTime int64) ([]sqlc.FMRecordBank, error) {
	return r.queries.GetUserBankRecordsWithPeriod(ctx, sqlc.GetUserBankRecordsWithPeriodParams{
		UserID:        userId,
		AccountNumber: accountNumber,
		CreatedAt:     utils.Int64ToPgTimestamptz(startTime),
		CreatedAt_2:   utils.Int64ToPgTimestamptz(endTime),
	})
}

func (r *bankRecordRepositoryImpl) GetUserBankRecordsFromDate(ctx context.Context, userId int64, accountNumber string, startTime int64) ([]sqlc.FMRecordBank, error) {
	return r.queries.GetUserBankRecordsFromDate(ctx, sqlc.GetUserBankRecordsFromDateParams{
		UserID:        userId,
		AccountNumber: accountNumber,
		CreatedAt:     utils.Int64ToPgTimestamptz(startTime),
	})
}

func (r *bankRecordRepositoryImpl) GetUserBankRecordsToDate(ctx context.Context, userId int64, accountNumber string, endTime int64) ([]sqlc.FMRecordBank, error) {
	return r.queries.GetUserBankRecordsToDate(ctx, sqlc.GetUserBankRecordsToDateParams{
		UserID:        userId,
		AccountNumber: accountNumber,
		CreatedAt:     utils.Int64ToPgTimestamptz(endTime),
	})
}

func (r *bankRecordRepositoryImpl) GetUserBankRecordsByType(ctx context.Context, userId int64, accountNumber, recordType string) ([]sqlc.FMRecordBank, error) {
	return r.queries.GetUserBankRecordsByType(ctx, sqlc.GetUserBankRecordsByTypeParams{
		UserID:        userId,
		AccountNumber: accountNumber,
		RecordType:    recordType,
	})
}

func (r *bankRecordRepositoryImpl) GetUserBankRecordsByTypeWithPeriod(ctx context.Context, userId int64, accountNumber, recordType string, startTime, endTime int64) ([]sqlc.FMRecordBank, error) {
	return r.queries.GetUserBankRecordsByTypeWithPeriod(ctx, sqlc.GetUserBankRecordsByTypeWithPeriodParams{
		UserID:        userId,
		AccountNumber: accountNumber,
		RecordType:    recordType,
		CreatedAt:     utils.Int64ToPgTimestamptz(startTime),
		CreatedAt_2:   utils.Int64ToPgTimestamptz(endTime),
	})
}

func (r *bankRecordRepositoryImpl) GetUserBankRecordsByTypeFromDate(ctx context.Context, userId int64, accountNumber, recordType string, startTime int64) ([]sqlc.FMRecordBank, error) {
	return r.queries.GetUserBankRecordsByTypeFromDate(ctx, sqlc.GetUserBankRecordsByTypeFromDateParams{
		UserID:        userId,
		AccountNumber: accountNumber,
		RecordType:    recordType,
		CreatedAt:     utils.Int64ToPgTimestamptz(startTime),
	})
}

func (r *bankRecordRepositoryImpl) GetUserBankRecordsByTypeToDate(ctx context.Context, userId int64, accountNumber, recordType string, endTime int64) ([]sqlc.FMRecordBank, error) {
	return r.queries.GetUserBankRecordsByTypeToDate(ctx, sqlc.GetUserBankRecordsByTypeToDateParams{
		UserID:        userId,
		AccountNumber: accountNumber,
		RecordType:    recordType,
		CreatedAt:     utils.Int64ToPgTimestamptz(endTime),
	})
}

func (r *bankRecordRepositoryImpl) GetExistingBankRecords(ctx context.Context, ids []int64) ([]int64, error) {
	return r.queries.GetExistingBankRecords(ctx, ids)
}
