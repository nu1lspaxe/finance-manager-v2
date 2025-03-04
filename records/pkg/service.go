package pkg

import "records/postgres/sqlc"

type RecordService struct {
	repo RecordRepository
}

func NewRecordService(repo RecordRepository) *RecordService {
	return &RecordService{repo: repo}
}

func (r *RecordService) CreateRecord(userId int64, amount float32, transactionDate int64, recordType string, detail string) (*sqlc.FinanceManagerFMRecord, error) {
	record, err := r.repo.CreateRecord(userId, amount, transactionDate, recordType, detail)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *RecordService) GetRecord(id int64) (*sqlc.FinanceManagerFMRecord, error) {
	record, err := r.repo.GetRecord(id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *RecordService) GetUserRecordsWithFilters(userId int64, recordType string, startTime int64, endTime int64) (*[]sqlc.FinanceManagerFMRecord, error) {

	var records []sqlc.FinanceManagerFMRecord
	var err error

	switch {
	case recordType != "" && startTime != 0 && endTime != 0:
		records, err = r.repo.GetUserRecordsByTypeWithPeriod(userId, recordType, startTime, endTime)
	case recordType != "" && startTime != 0:
		records, err = r.repo.GetUserRecordsByTypeFromDate(userId, recordType, startTime)
	case recordType != "" && endTime != 0:
		records, err = r.repo.GetUserRecordsByTypeToDate(userId, recordType, endTime)
	case recordType != "":
		records, err = r.repo.GetUserRecordsByType(userId, recordType)
	case startTime != 0 && endTime != 0:
		records, err = r.repo.GetUserRecordsWithPeriod(userId, startTime, endTime)
	case startTime != 0:
		records, err = r.repo.GetUserRecordsFromDate(userId, startTime)
	case endTime != 0:
		records, err = r.repo.GetUserRecordsToDate(userId, endTime)
	default:
		records, err = r.repo.GetUserRecords(userId)
	}

	if err != nil {
		return nil, err
	}
	return &records, nil
}

func (r *RecordService) UpdateRecord(id int64, amount float32, transactionDate int64, recordType string, detail string) error {
	err := r.repo.UpdateRecord(id, amount, transactionDate, recordType, detail)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecordService) DeleteRecord(id int64) error {
	err := r.repo.DeleteRecord(id)
	if err != nil {
		return err
	}
	return nil
}
