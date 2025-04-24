package utils

import (
	"records_bank/postgres/sqlc"
	"records_bank/proto"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func SqlcToProto(bankRecord *sqlc.FMRecordBank) (*proto.BankRecord, error) {
	return &proto.BankRecord{
		Id:              bankRecord.ID,
		UserId:          bankRecord.UserID,
		Amount:          bankRecord.Amount,
		TransactionDate: bankRecord.TransactionDate.Time.Unix(),
		Detail:          bankRecord.Detail,
		CreatedAt:       bankRecord.CreatedAt.Time.Unix(),
	}, nil
}

func Int64ToPgTimestamptz(date int64) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:             time.Unix(date, 0),
		InfinityModifier: 0,
		Valid:            true,
	}
}
