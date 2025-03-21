package pkg

import "github.com/jackc/pgx/v5/pgtype"

type TxResponse struct {
	ID        int64              `json:"id"`
	Amount    float64            `json:"amount"`
	TxType    string             `json:"tx_type"`
	Detail    string             `json:"detail"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
