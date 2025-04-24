package pkg

import (
	"context"
	"records/postgres/sqlc"
	"records/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordRepository interface {
	CreateRecord(ctx context.Context, userId int64, amount float32, transactionDate int64, recordType string, detail string) (sqlc.FMRecord, error)
	GetRecord(ctx context.Context, id int64) (*sqlc.FMRecord, error)
	GetUserRecords(ctx context.Context, userId int64) ([]sqlc.FMRecord, error)
	GetUserRecordsWithPeriod(ctx context.Context, userId int64, startTime int64, endTime int64) ([]sqlc.FMRecord, error)
	GetUserRecordsFromDate(ctx context.Context, userId int64, date int64) ([]sqlc.FMRecord, error)
	GetUserRecordsToDate(ctx context.Context, userId int64, date int64) ([]sqlc.FMRecord, error)
	GetUserRecordsByType(ctx context.Context, userId int64, recordType string) ([]sqlc.FMRecord, error)
	GetUserRecordsByTypeWithPeriod(ctx context.Context, userId int64, recordType string, startTime int64, endTime int64) ([]sqlc.FMRecord, error)
	GetUserRecordsByTypeFromDate(ctx context.Context, userId int64, recordType string, date int64) ([]sqlc.FMRecord, error)
	GetUserRecordsByTypeToDate(ctx context.Context, userId int64, recordType string, date int64) ([]sqlc.FMRecord, error)
	UpdateRecord(ctx context.Context, id int64, amount float32, transactionDate int64, recordType string, detail string) error
	DeleteRecord(ctx context.Context, id int64) error
}

type recordRepositoryImpl struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRecordRepository(pool *pgxpool.Pool) RecordRepository {
	return &recordRepositoryImpl{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r recordRepositoryImpl) CreateRecord(ctx context.Context, userId int64, amount float32, transactionDate int64, recordType string, detail string) (sqlc.FMRecord, error) {
	txDate := utils.Int64ToPgDate(transactionDate)

	record, err := r.queries.CreateRecord(ctx, sqlc.CreateRecordParams{
		UserID:          userId,
		Amount:          amount,
		TransactionDate: txDate,
		RecordType:      recordType,
		Detail:          detail,
	})
	if err != nil {
		return sqlc.FMRecord{}, err
	}
	return record, nil
}

func (r recordRepositoryImpl) GetRecord(ctx context.Context, id int64) (*sqlc.FMRecord, error) {
	record, err := r.queries.GetRecord(ctx, id)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r recordRepositoryImpl) GetUserRecords(ctx context.Context, userId int64) ([]sqlc.FMRecord, error) {
	records, err := r.queries.GetUserRecords(ctx, userId)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsWithPeriod(ctx context.Context, userId int64, startTime int64, endTime int64) ([]sqlc.FMRecord, error) {
	sDate := utils.Int64ToPgDate(startTime)
	eDate := utils.Int64ToPgDate(endTime)

	records, err := r.queries.GetUserRecordsWithPeriod(ctx, sqlc.GetUserRecordsWithPeriodParams{
		UserID:            userId,
		TransactionDate:   sDate,
		TransactionDate_2: eDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsFromDate(ctx context.Context, userId int64, date int64) ([]sqlc.FMRecord, error) {
	txDate := utils.Int64ToPgDate(date)

	records, err := r.queries.GetUserRecordsFromDate(ctx, sqlc.GetUserRecordsFromDateParams{
		UserID:          userId,
		TransactionDate: txDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsToDate(ctx context.Context, userId int64, date int64) ([]sqlc.FMRecord, error) {
	txDate := utils.Int64ToPgDate(date)

	records, err := r.queries.GetUserRecordsToDate(ctx, sqlc.GetUserRecordsToDateParams{
		UserID:          userId,
		TransactionDate: txDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsByType(ctx context.Context, userId int64, recordType string) ([]sqlc.FMRecord, error) {
	records, err := r.queries.GetUserRecordsByType(ctx, sqlc.GetUserRecordsByTypeParams{
		UserID:     userId,
		RecordType: recordType,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsByTypeWithPeriod(ctx context.Context, userId int64, recordType string, startTime int64, endTime int64) ([]sqlc.FMRecord, error) {
	sDate := utils.Int64ToPgDate(startTime)
	eDate := utils.Int64ToPgDate(endTime)

	records, err := r.queries.GetUserRecordsByTypeWithPeriod(ctx, sqlc.GetUserRecordsByTypeWithPeriodParams{
		UserID:            userId,
		RecordType:        recordType,
		TransactionDate:   sDate,
		TransactionDate_2: eDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsByTypeFromDate(ctx context.Context, userId int64, recordType string, date int64) ([]sqlc.FMRecord, error) {
	txDate := utils.Int64ToPgDate(date)

	records, err := r.queries.GetUserRecordsByTypeFromDate(ctx, sqlc.GetUserRecordsByTypeFromDateParams{
		UserID:          userId,
		RecordType:      recordType,
		TransactionDate: txDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsByTypeToDate(ctx context.Context, userId int64, recordType string, date int64) ([]sqlc.FMRecord, error) {
	txDate := utils.Int64ToPgDate(date)

	records, err := r.queries.GetUserRecordsByTypeToDate(ctx, sqlc.GetUserRecordsByTypeToDateParams{
		UserID:          userId,
		RecordType:      recordType,
		TransactionDate: txDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) UpdateRecord(ctx context.Context, id int64, amount float32, transactionDate int64, recordType string, detail string) error {
	txDate := utils.Int64ToPgDate(transactionDate)

	err := r.queries.UpdateRecord(ctx, sqlc.UpdateRecordParams{
		ID:              id,
		Amount:          amount,
		TransactionDate: txDate,
		RecordType:      recordType,
		Detail:          detail,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r recordRepositoryImpl) DeleteRecord(ctx context.Context, id int64) error {
	err := r.queries.DeleteRecord(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
