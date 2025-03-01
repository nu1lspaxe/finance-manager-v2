package utils

import (
	"records/postgres/sqlc"
	"records/proto"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func SqlcToProto(record *sqlc.Record) (*proto.Record, error) {

	return &proto.Record{
		Id:        record.ID,
		UserId:    record.UserID,
		Amount:    record.Amount,
		CreatedAt: record.CreatedAt.Time.Unix(),
		UpdatedAt: record.UpdatedAt.Time.Unix(),
	}, nil
}

func Int64ToPgDate(date int64) pgtype.Date {
	return pgtype.Date{
		Time:             time.Unix(date, 0),
		InfinityModifier: 0,
		Valid:            true,
	}
}

func StringToPgText(text string) (pgtype.Text, error) {
	var pgText pgtype.Text
	if err := pgText.Scan(text); err != nil {
		return pgtype.Text{}, err
	}
	return pgText, nil
}
