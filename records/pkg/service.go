package pkg

import (
	"context"
	"records/postgres/sqlc"
)

type RecordService struct {
	repo RecordRepository
}

func NewRecordService(repo RecordRepository) *RecordService {
	return &RecordService{repo: repo}
}

func (r *RecordService) CreateRecord(ctx context.Context, userId int64, amount float32, transactionDate int64, recordType string, detail string) (*sqlc.FMRecord, error) {
	record, err := r.repo.CreateRecord(ctx, userId, amount, transactionDate, recordType, detail)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RecordService) GetRecord(ctx context.Context, id int64) (*sqlc.FMRecord, error) {
	record, err := r.repo.GetRecord(ctx, id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *RecordService) GetUserRecordsWithFilters(ctx context.Context, userId int64, recordType string, startTime int64, endTime int64) ([]sqlc.FMRecord, error) {

	var records []sqlc.FMRecord
	var err error

	switch {
	case recordType != "" && startTime != 0 && endTime != 0:
		records, err = r.repo.GetUserRecordsByTypeWithPeriod(ctx, userId, recordType, startTime, endTime)
	case recordType != "" && startTime != 0:
		records, err = r.repo.GetUserRecordsByTypeFromDate(ctx, userId, recordType, startTime)
	case recordType != "" && endTime != 0:
		records, err = r.repo.GetUserRecordsByTypeToDate(ctx, userId, recordType, endTime)
	case recordType != "":
		records, err = r.repo.GetUserRecordsByType(ctx, userId, recordType)
	case startTime != 0 && endTime != 0:
		records, err = r.repo.GetUserRecordsWithPeriod(ctx, userId, startTime, endTime)
	case startTime != 0:
		records, err = r.repo.GetUserRecordsFromDate(ctx, userId, startTime)
	case endTime != 0:
		records, err = r.repo.GetUserRecordsToDate(ctx, userId, endTime)
	default:
		records, err = r.repo.GetUserRecords(ctx, userId)
	}

	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r *RecordService) UpdateRecord(ctx context.Context, id int64, amount float32, transactionDate int64, recordType string, detail string) error {
	err := r.repo.UpdateRecord(ctx, id, amount, transactionDate, recordType, detail)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecordService) DeleteRecord(ctx context.Context, id int64) error {
	err := r.repo.DeleteRecord(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
