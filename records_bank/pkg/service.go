package pkg

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"records_bank/postgres/sqlc"
	"records_bank/utils"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type BankRecordService struct {
	repo        BankRecordRepository
	httpClient  *http.Client
	kafkaReader *kafka.Reader
}

func NewBankRecordService(
	repo BankRecordRepository, tlsConfig *tls.Config, kafkaReader *kafka.Reader,
) *BankRecordService {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   tlsConfig,
			ForceAttemptHTTP2: true,
		},
		Timeout: utils.TIMEOUT,
	}

	return &BankRecordService{
		repo:        repo,
		httpClient:  client,
		kafkaReader: kafkaReader,
	}
}

func (s *BankRecordService) CreateBankRecordsBulk(ctx context.Context, userId int64, accountNumber string) error {
	txs, err := fetchBankAccountTxs(ctx, s.httpClient, accountNumber)
	if err != nil {
		return err
	}

	asyncTxs, err := s.SyncBankAccountTxs(ctx, txs)
	if err != nil {
		return err
	}

	txIds := []int64{}
	amounts := []float64{}
	txTypes := []string{}
	txDates := []pgtype.Timestamptz{}
	details := []string{}

	for _, record := range asyncTxs {
		txIds = append(txIds, record.ID)
		amounts = append(amounts, record.Amount)
		txTypes = append(txTypes, record.TxType)
		txDates = append(txDates, record.CreatedAt)
		details = append(details, record.Detail)
	}

	err = s.repo.CreateBankRecordsBulk(ctx, userId, accountNumber, amounts, txTypes, txIds, txDates, details)
	if err != nil {
		return err
	}

	return nil
}

func fetchBankAccountTxs(ctx context.Context, client *http.Client, idNumber string) ([]TxResponse, error) {
	baseURL := viper.GetString("fintech.url")
	url := fmt.Sprintf("%s/accounts/%s/transactions", baseURL, idNumber)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []TxResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "users-service/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return []TxResponse{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return []TxResponse{}, utils.NewBankRecordError(
			utils.ErrStatusCode, fmt.Sprint(resp.StatusCode),
		)
	}

	var txResp []TxResponse
	err = json.NewDecoder(resp.Body).Decode(&txResp)
	if err != nil {
		return []TxResponse{}, err
	}
	return txResp, nil
}

func (s *BankRecordService) SyncBankAccountTxs(ctx context.Context, txs []TxResponse) ([]TxResponse, error) {
	txIds := make([]int64, 0, len(txs))
	for _, tx := range txs {
		txIds = append(txIds, tx.ID)
	}

	existIds, err := s.repo.GetExistingBankRecords(ctx, txIds)
	if err != nil {
		return nil, err
	}

	var newTxs []TxResponse
	for _, tx := range txs {
		if !utils.Int64InSlice(tx.ID, existIds) {
			newTxs = append(newTxs, tx)
		}
	}

	return newTxs, nil
}

func (s *BankRecordService) GetBankRecord(ctx context.Context, id int64) (*sqlc.FMRecordBank, error) {
	record, err := s.repo.GetBankRecord(ctx, id)
	if err != nil {
		return nil, err
	}
	return &record, nil

}

func (s *BankRecordService) GetBankRecordsWithFilters(
	ctx context.Context, userId int64, accountNumber string, recordType string, startTime int64, endTime int64,
) ([]sqlc.FMRecordBank, error) {
	var records []sqlc.FMRecordBank
	var err error

	switch {
	case recordType != "" && startTime != 0 && endTime != 0:
		records, err = s.repo.GetUserBankRecordsByTypeWithPeriod(ctx, userId, accountNumber, recordType, startTime, endTime)
	case recordType != "" && startTime != 0:
		records, err = s.repo.GetUserBankRecordsByTypeFromDate(ctx, userId, accountNumber, recordType, startTime)
	case recordType != "" && endTime != 0:
		records, err = s.repo.GetUserBankRecordsByTypeToDate(ctx, userId, accountNumber, recordType, endTime)
	case recordType != "":
		records, err = s.repo.GetUserBankRecordsByType(ctx, userId, accountNumber, recordType)
	case startTime != 0 && endTime != 0:
		records, err = s.repo.GetUserBankRecordsWithPeriod(ctx, userId, accountNumber, startTime, endTime)
	case startTime != 0:
		records, err = s.repo.GetUserBankRecordsFromDate(ctx, userId, accountNumber, startTime)
	case endTime != 0:
		records, err = s.repo.GetUserBankRecordsToDate(ctx, userId, accountNumber, endTime)
	default:
		records, err = s.repo.GetUserBankRecords(ctx, userId, accountNumber)
	}

	if err != nil {
		return nil, err
	}
	return records, nil

}
