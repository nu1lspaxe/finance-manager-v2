package pkg

import (
	"context"
	"records/postgres/sqlc"
	"records/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RecordRepository interface {
	CreateRecord(userId int64, amount float32, transactionDate int64, recordType string, detail string) (sqlc.FinanceManagerFMRecord, error)
	GetRecord(id int64) (*sqlc.FinanceManagerFMRecord, error)
	GetUserRecords(userId int64) ([]sqlc.FinanceManagerFMRecord, error)
	GetUserRecordsWithPeriod(userId int64, startTime int64, endTime int64) ([]sqlc.FinanceManagerFMRecord, error)
	GetUserRecordsFromDate(userId int64, date int64) ([]sqlc.FinanceManagerFMRecord, error)
	GetUserRecordsToDate(userId int64, date int64) ([]sqlc.FinanceManagerFMRecord, error)
	GetUserRecordsByType(userId int64, recordType string) ([]sqlc.FinanceManagerFMRecord, error)
	GetUserRecordsByTypeWithPeriod(userId int64, recordType string, startTime int64, endTime int64) ([]sqlc.FinanceManagerFMRecord, error)
	GetUserRecordsByTypeFromDate(userId int64, recordType string, date int64) ([]sqlc.FinanceManagerFMRecord, error)
	GetUserRecordsByTypeToDate(userId int64, recordType string, date int64) ([]sqlc.FinanceManagerFMRecord, error)
	UpdateRecord(id int64, amount float32, transactionDate int64, recordType string, detail string) error
	DeleteRecord(id int64) error
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

func (r recordRepositoryImpl) CreateRecord(userId int64, amount float32, transactionDate int64, recordType string, detail string) (sqlc.FinanceManagerFMRecord, error) {

	txDate := utils.Int64ToPgDate(transactionDate)

	record, err := r.queries.CreateRecord(context.Background(), sqlc.CreateRecordParams{
		UserID:          userId,
		Amount:          amount,
		TransactionDate: txDate,
		RecordType:      recordType,
		Detail:          detail,
	})
	if err != nil {
		return sqlc.FinanceManagerFMRecord{}, err
	}
	return record, nil
}

func (r recordRepositoryImpl) GetRecord(id int64) (*sqlc.FinanceManagerFMRecord, error) {
	record, err := r.queries.GetRecord(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r recordRepositoryImpl) GetUserRecords(userId int64) ([]sqlc.FinanceManagerFMRecord, error) {
	records, err := r.queries.GetUserRecords(context.Background(), userId)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsWithPeriod(userId int64, startTime int64, endTime int64) ([]sqlc.FinanceManagerFMRecord, error) {
	sDate := utils.Int64ToPgDate(startTime)
	eDate := utils.Int64ToPgDate(endTime)

	records, err := r.queries.GetUserRecordsWithPeriod(context.Background(), sqlc.GetUserRecordsWithPeriodParams{
		UserID:            userId,
		TransactionDate:   sDate,
		TransactionDate_2: eDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsFromDate(userId int64, date int64) ([]sqlc.FinanceManagerFMRecord, error) {
	txDate := utils.Int64ToPgDate(date)

	records, err := r.queries.GetUserRecordsFromDate(context.Background(), sqlc.GetUserRecordsFromDateParams{
		UserID:          userId,
		TransactionDate: txDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsToDate(userId int64, date int64) ([]sqlc.FinanceManagerFMRecord, error) {
	txDate := utils.Int64ToPgDate(date)

	records, err := r.queries.GetUserRecordsToDate(context.Background(), sqlc.GetUserRecordsToDateParams{
		UserID:          userId,
		TransactionDate: txDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsByType(userId int64, recordType string) ([]sqlc.FinanceManagerFMRecord, error) {
	records, err := r.queries.GetUserRecordsByType(context.Background(), sqlc.GetUserRecordsByTypeParams{
		UserID:     userId,
		RecordType: recordType,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsByTypeWithPeriod(userId int64, recordType string, startTime int64, endTime int64) ([]sqlc.FinanceManagerFMRecord, error) {
	sDate := utils.Int64ToPgDate(startTime)
	eDate := utils.Int64ToPgDate(endTime)

	records, err := r.queries.GetUserRecordsByTypeWithPeriod(context.Background(), sqlc.GetUserRecordsByTypeWithPeriodParams{
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

func (r recordRepositoryImpl) GetUserRecordsByTypeFromDate(userId int64, recordType string, date int64) ([]sqlc.FinanceManagerFMRecord, error) {
	txDate := utils.Int64ToPgDate(date)

	records, err := r.queries.GetUserRecordsByTypeFromDate(context.Background(), sqlc.GetUserRecordsByTypeFromDateParams{
		UserID:          userId,
		RecordType:      recordType,
		TransactionDate: txDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) GetUserRecordsByTypeToDate(userId int64, recordType string, date int64) ([]sqlc.FinanceManagerFMRecord, error) {
	txDate := utils.Int64ToPgDate(date)

	records, err := r.queries.GetUserRecordsByTypeToDate(context.Background(), sqlc.GetUserRecordsByTypeToDateParams{
		UserID:          userId,
		RecordType:      recordType,
		TransactionDate: txDate,
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryImpl) UpdateRecord(id int64, amount float32, transactionDate int64, recordType string, detail string) error {
	txDate := utils.Int64ToPgDate(transactionDate)

	err := r.queries.UpdateRecord(context.Background(), sqlc.UpdateRecordParams{
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

func (r recordRepositoryImpl) DeleteRecord(id int64) error {
	err := r.queries.DeleteRecord(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}
