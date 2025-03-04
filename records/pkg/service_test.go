package pkg

import (
	"errors"
	"records/pkg/mocks"
	"records/postgres/sqlc"
	"records/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateRecord(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("CreateRecord", int64(1), float32(100.0), time.Now().Unix(), "expense", "test record").
			Return(sqlc.FinanceManagerFMRecord{
				ID:              1,
				UserID:          1,
				Amount:          100.0,
				TransactionDate: utils.Int64ToPgDate(time.Now().Unix()),
				RecordType:      "expense",
				Detail:          "test record",
			}, nil)

		service := NewRecordService(mockRepo)

		_, err := service.CreateRecord(1, 100.0, time.Now().Unix(), "expense", "test record")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("CreateRecord", int64(1), float32(100.0), time.Now().Unix(), "expense", "test record").
			Return(sqlc.FinanceManagerFMRecord{}, errors.New("repository error"))

		service := NewRecordService(mockRepo)

		_, err := service.CreateRecord(1, 100.0, time.Now().Unix(), "expense", "test record")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetRecord(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("GetRecord", int64(1)).
			Return(&sqlc.FinanceManagerFMRecord{
				ID:              1,
				UserID:          1,
				Amount:          100.0,
				TransactionDate: utils.Int64ToPgDate(time.Now().Unix()),
				RecordType:      "expense",
				Detail:          "test record",
			}, nil)

		service := NewRecordService(mockRepo)

		_, err := service.GetRecord(1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("GetRecord", int64(1)).
			Return(&sqlc.FinanceManagerFMRecord{}, errors.New("repository error"))

		service := NewRecordService(mockRepo)

		_, err := service.GetRecord(1)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserRecordsWithFilters(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("GetUserRecords", int64(1)).
			Return([]sqlc.FinanceManagerFMRecord{
				{
					ID:              1,
					UserID:          1,
					Amount:          100.0,
					TransactionDate: utils.Int64ToPgDate(time.Now().Unix()),
					RecordType:      "expense",
					Detail:          "test record",
				},
			}, nil)

		service := NewRecordService(mockRepo)

		_, err := service.GetUserRecordsWithFilters(1, "", 0, 0)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("GetUserRecords", int64(1)).
			Return([]sqlc.FinanceManagerFMRecord{}, errors.New("repository error"))

		service := NewRecordService(mockRepo)

		_, err := service.GetUserRecordsWithFilters(1, "", 0, 0)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateRecord(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("UpdateRecord", int64(1), float32(100.0), time.Now().Unix(), "expense", "test record").
			Return(nil)

		service := NewRecordService(mockRepo)

		err := service.UpdateRecord(1, 100.0, time.Now().Unix(), "expense", "test record")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("UpdateRecord", int64(1), float32(100.0), time.Now().Unix(), "expense", "test record").
			Return(errors.New("repository error"))

		service := NewRecordService(mockRepo)

		err := service.UpdateRecord(1, 100.0, time.Now().Unix(), "expense", "test record")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteRecord(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("DeleteRecord", int64(1)).
			Return(nil)

		service := NewRecordService(mockRepo)

		err := service.DeleteRecord(1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error case", func(t *testing.T) {
		mockRepo := mocks.NewRecordRepository(t)

		mockRepo.
			On("DeleteRecord", int64(1)).
			Return(errors.New("repository error"))

		service := NewRecordService(mockRepo)

		err := service.DeleteRecord(1)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
